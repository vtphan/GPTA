// Author: Vinhthuy Phan, 2018
package restHandlers

import (
	"fmt"
	"github.com/GPTA/src/frontEnd"
	"github.com/GPTA/src/repository"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
	"sort"
	"strconv"
)

// -----------------------------------------------------------------------------------
type SubmissionData struct {
	Flag      string
	Start     int64
	At        int64
	Completed int64
}

// -----------------------------------------------------------------------------------
func AnalyzeSubmissionsHandler(w http.ResponseWriter, r *http.Request) {
	pidStr := r.FormValue("pid")
	pid, err := strconv.Atoi(pidStr) // Convert string to int
	if err != nil {
		fmt.Println("Error converting pid to int:", err)
		http.Error(w, "Invalid problem ID", http.StatusBadRequest)
		return
	}

	// Use GetProblemUploadTime to retrieve the problem upload time
	start, err := repository.GetProblemUploadTime(pid)
	if err != nil {
		fmt.Println("Error retrieving problem upload time:", err)
		http.Error(w, "Failed to retrieve problem upload time", http.StatusInternalServerError)
		return
	}

	records := make(map[int][]*SubmissionData)

	// Use GetSubmissionsByProblemID to retrieve the submissions for the problem
	submissions, err := repository.GetSubmissionsByProblemID(pid)
	if err != nil {
		fmt.Println("Error retrieving submissions:", err)
		http.Error(w, "Failed to retrieve submissions", http.StatusInternalServerError)
		return
	}

	for _, submission := range submissions {
		sid := submission.StudentID
		priority := submission.SubmissionCategory
		at := submission.CodeSubmittedAt
		completed := submission.Completed

		if _, ok := records[sid]; !ok {
			records[sid] = make([]*SubmissionData, 0)
		}

		flag := "unknown"
		if priority == 1 {
			flag = "Got it!"
		} else if priority == 2 {
			flag = "Help!"
		}

		records[sid] = append(
			records[sid],
			&SubmissionData{
				Flag:      flag,
				Start:     start.UnixNano(),
				At:        at.UnixNano(),
				Completed: completed.UnixNano(),
			})
	}

	// Sort submissions by submission time (at)
	for sid, _ := range records {
		sort.Slice(records[sid], func(i, j int) bool {
			return records[sid][i].At < records[sid][j].At
		})
	}

	w.Header().Set("Content-Type", "text/html")
	t, _ := template.New("").Parse(frontEnd.ANALYZE_SUBMISSIONS_TEMPLATE)
	err = t.Execute(w, records)
	if err != nil {
		fmt.Println(err)
	}
}

// -----------------------------------------------------------------------------------

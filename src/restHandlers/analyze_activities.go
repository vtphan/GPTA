// Author: Vinhthuy Phan, 2018
package restHandlers

import (
	"fmt"
	"github.com/GPTA/src/frontEnd"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
	"time"
)

// -----------------------------------------------------------------------------------
type DailyActivityData struct {
	Pids     map[int]bool
	Sids     map[int]bool
	Count    int
	PidCount int
	SidCount int
}

// -----------------------------------------------------------------------------------
func ViewActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("pc") != models.Passcode {
		fmt.Fprintf(w, "Unauthorized")
		return
	}
	submissions, err := repository.GetSubmissions()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving submissions: %v", err), http.StatusInternalServerError)
		return
	}

	data := make(map[int64]*DailyActivityData)

	for _, submission := range submissions {
		date := time.Date(submission.CodeSubmittedAt.Year(), submission.CodeSubmittedAt.Month(), submission.CodeSubmittedAt.Day(), 0, 0, 0, 0, submission.CodeSubmittedAt.Location()).UnixNano()
		if _, ok := data[date]; !ok {
			data[date] = &DailyActivityData{
				Pids:  make(map[int]bool),
				Sids:  make(map[int]bool),
				Count: 0,
			}
		}
		data[date].Count++
		data[date].Pids[submission.ProblemID] = true
		data[date].Sids[submission.StudentID] = true
	}

	for d, _ := range data {
		data[d].PidCount = len(data[d].Pids)
		data[d].SidCount = len(data[d].Sids)
	}
	w.Header().Set("Content-Type", "text/html")
	t, err := template.New("").Parse(frontEnd.ACTIVITY_VIEW_TEMPLATE)
	if err != nil {
		fmt.Println(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

// -----------------------------------------------------------------------------------

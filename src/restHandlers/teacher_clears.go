// Author: Vinhthuy Phan, 2018
package restHandlers

import (
	// "encoding/json"
	"fmt"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"log"
	"strconv"
	"time"

	// "log"
	"net/http"
)

// -----------------------------------------------------------------------------------
// When problems are deactivated, boards cleared, no new submissions are possibile.
// -----------------------------------------------------------------------------------
func TeacherDeactivatesProblemsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	models.CodeSnapshotSem.Lock()
	defer models.CodeSnapshotSem.Unlock()
	filename := r.FormValue("filename")
	problemIDStr := r.FormValue("problem_id")
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		http.Error(w, "Invalid problem_id", http.StatusBadRequest)
		return
	}
	if prob, ok := models.ActiveProblems[filename]; ok {
		prob.Active = false
		models.PeerTutorAllowed = false
		if len(prob.Answers) > 0 {
			fmt.Fprintf(w, "1")
		} else {
			fmt.Fprintf(w, "0")
		}
		tempSnapshots := models.Snapshots
		for studentID, _ := range models.StudentSnapshot {
			models.StudentSnapshot[studentID] = map[int]int{}
		}
		models.Snapshots = make([]*models.Snapshot, 0)
		idx := 0
		for _, s := range tempSnapshots {
			if s.ProblemID != prob.Info.Pid {
				models.Snapshots = append(models.Snapshots, s)
				for studentID, _ := range models.StudentSnapshot {
					models.StudentSnapshot[studentID][prob.Info.Pid] = idx
				}
				idx++
			}
		}
		err := repository.UpdateProblemEndTime(time.Now(), problemID)
		if err != nil {
			log.Println(err)
		}
		for studentID, _ := range models.Students {
			for i, b := range models.Students[studentID].Boards {
				if b.Pid == prob.Info.Pid {
					models.Students[studentID].Boards = append(models.Students[studentID].Boards[:i], models.Students[studentID].Boards[i+1:]...)
					break
				}
			}
		}

	} else {
		fmt.Fprintf(w, "-1")
	}

	// filenames := make([]string, 0)
	// for fname, prob := range ActiveProblems {
	// 	if prob.Active {
	// 		prob.Active = false
	// 		if len(prob.Answers) > 0 {
	// 			filenames = append(filenames, fname)
	// 		}
	// 	}
	// }
	// for stid, _ := range Students {
	// 	Students[stid].Boards = make([]*Board, 0)
	// 	Students[stid].SubmissionStatus = 0
	// }
	// js, err := json.Marshal(filenames)
	// if err == nil {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.Write(js)
	// } else {
	// 	log.Fatal(err)
	// }
}

// -----------------------------------------------------------------------------------
// Clear submissions, boards, statuses, and set all problems inactive.
// -----------------------------------------------------------------------------------
func TeacherClearsSubmissionsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	models.WorkingSubs = make([]*models.Submission, 0)
	fmt.Fprintf(w, "Done.")
}

// Author: Vinhthuy Phan, 2018
package main

import (
	"encoding/json"
	"fmt"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"log"
	"net/http"
)

//-----------------------------------------------------------------

func student_checks_inHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	// attendance is taken automatically by authorization when this handler is called.
	// Next: return student attendance report
	attendances, err := repository.GetAttendanceByStudentID(uid)
	if err != nil {
		http.Error(w, "Failed to retrieve attendance records", http.StatusInternalServerError)
		log.Printf("Error retrieving attendance for student ID %d: %v", uid, err)
		return
	}

	// Converting attendance timestamps to Unix format
	dates := make([]int64, 0)
	for _, att := range attendances {
		dates = append(dates, att.AttendanceAt.Unix())
	}

	// Marshalling dates into JSON
	js, err := json.Marshal(dates)
	if err != nil {
		http.Error(w, "Failed to encode attendance records", http.StatusInternalServerError)
		log.Printf("Error marshalling attendance records: %v", err)
		return
	}

	// Sending the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// -----------------------------------------------------------------
func student_periodic_updateHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	submissionStat := &models.StudentSubmissionStatus{}
	if len(models.Students[uid].SubmissionStatus) > 0 {
		submissionStat = models.Students[uid].SubmissionStatus[0]
		models.Students[uid].SubmissionStatus = models.Students[uid].SubmissionStatus[1:]
	} else {
		submissionStat = &models.StudentSubmissionStatus{
			Filename:      "",
			AttemptNumber: 0,
			Status:        0,
		}
	}
	// submission_stat := Students[uid].SubmissionStatus
	thank_stat := models.Students[uid].ThankStatus
	board_stat := 0
	if len(models.Students[uid].Boards) > 0 {
		board_stat = 1
		for _, b := range models.Students[uid].Boards {
			if b.Type == "peer_feedback" {
				board_stat = 2
				break
			}
		}
	}
	// Students[uid].SubmissionStatus = 0 // reset status after notifying student
	models.Students[uid].ThankStatus = 0
	snapShotFeedbackStatus := 0
	if len(models.Students[uid].SnapShotFeedbackQueue) > 0 {
		snapShotFeedbackStatus = 1
	}
	fmt.Fprintf(w, "%d;%d;%d;%d;%d;%s", submissionStat.Status, board_stat, thank_stat, snapShotFeedbackStatus, submissionStat.AttemptNumber, submissionStat.Filename)
}

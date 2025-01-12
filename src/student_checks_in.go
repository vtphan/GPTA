// Author: Vinhthuy Phan, 2018
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//-----------------------------------------------------------------

func studentChecksInHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	// Retrieve attendance records for the given student ID.
	var attendances []Attendance
	err := Database.Where("student_id = ?", uid).Find(&attendances).Error
	if err != nil {
		http.Error(w, "Failed to fetch attendance records", http.StatusInternalServerError)
		log.Printf("Error fetching attendance records: %v", err)
		return
	}

	// Convert attendance timestamps to UNIX format.
	dates := make([]int64, len(attendances))
	for i, attendance := range attendances {
		dates[i] = attendance.AttendanceAt.Unix()
	}

	// Marshal the dates into JSON format.
	js, err := json.Marshal(dates)
	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		log.Printf("Error encoding JSON: %v", err)
		return
	}

	// Set the response header and write the JSON.
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// -----------------------------------------------------------------
func student_periodic_updateHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	submissionStat := &StudentSubmissionStatus{}
	if len(Students[uid].SubmissionStatus) > 0 {
		submissionStat = Students[uid].SubmissionStatus[0]
		Students[uid].SubmissionStatus = Students[uid].SubmissionStatus[1:]
	} else {
		submissionStat = &StudentSubmissionStatus{
			Filename:      "",
			AttemptNumber: 0,
			Status:        0,
		}
	}
	// submission_stat := Students[uid].SubmissionStatus
	thank_stat := Students[uid].ThankStatus
	board_stat := 0
	if len(Students[uid].Boards) > 0 {
		board_stat = 1
		for _, b := range Students[uid].Boards {
			if b.Type == "peer_feedback" {
				board_stat = 2
				break
			}
		}
	}
	// Students[uid].SubmissionStatus = 0 // reset status after notifying student
	Students[uid].ThankStatus = 0
	snapShotFeedbackStatus := 0
	if len(Students[uid].SnapShotFeedbackQueue) > 0 {
		snapShotFeedbackStatus = 1
	}
	fmt.Fprintf(w, "%d;%d;%d;%d;%d;%s", submissionStat.Status, board_stat, thank_stat, snapShotFeedbackStatus, submissionStat.AttemptNumber, submissionStat.Filename)
}

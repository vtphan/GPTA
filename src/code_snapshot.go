package main

import (
	"encoding/json"
	"fmt"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"log"
	"net/http"
	"strconv"
	"time"
)

func addCodeSnapshot(studentID int, problemID int, code string, status int, lastUpdate time.Time, event string) int {
	result, err := repository.AddCodeSnapshot(studentID, problemID, code, status, lastUpdate, event)
	if err != nil {
		log.Fatal("Could not save the snapshot for error: ", err)
		return -1
	}
	snapshotID := result.ID
	idx, ok := models.StudentSnapshot[studentID][problemID]
	if !ok {
		idx = len(models.Snapshots)
		models.StudentSnapshot[studentID][problemID] = idx
		name := repository.GetStudentName(studentID)
		if err != nil {
			log.Fatal("Could not retrieve student name: ", err)
			return -1
		}
		problemName := ""
		for _, problem := range models.ActiveProblems {
			if problem.Active == true && problem.Info.Pid == problemID {
				problemName = problem.Info.Filename
				break
			}
		}
		models.Snapshots = append(models.Snapshots, &models.Snapshot{
			ID:          int(snapshotID),
			StudentName: name,
			StudentID:   studentID,
			ProblemName: problemName,
			ProblemID:   problemID,
			Status:      status,
			FirstUpdate: lastUpdate,
			LastUpdated: lastUpdate,
			LinesOfCode: models.GetLinesOfCode(code),
			Code:        code,
		})
	} else {
		currentStatus := models.Snapshots[idx].Status
		if currentStatus > status {
			status = currentStatus
		}
		models.Snapshots[idx] = &models.Snapshot{
			ID:          int(snapshotID),
			StudentName: models.Snapshots[idx].StudentName,
			StudentID:   studentID,
			ProblemName: models.Snapshots[idx].ProblemName,
			ProblemID:   problemID,
			Status:      status,
			FirstUpdate: models.Snapshots[idx].FirstUpdate,
			LastUpdated: lastUpdate,
			LinesOfCode: models.GetLinesOfCode(code),
			Code:        code,
			NumFeedback: models.Snapshots[idx].NumFeedback,
		}
	}
	return snapshotID
}

func codeSnapshotHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	code := r.FormValue("code")
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	studentID, _ := strconv.Atoi(r.FormValue("uid"))
	snapshotEvent := r.FormValue("event")
	if snapshotEvent == "" {
		snapshotEvent = "at_regular_interval"
	}
	addCodeSnapshot(studentID, problemID, code, 0, time.Now(), snapshotEvent)
}

func codeSnapshotFeedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	snapshotID, _ := strconv.Atoi(r.FormValue("snapshot_id"))
	feedback := r.FormValue("feedback")
	authorID, _ := strconv.Atoi(r.FormValue("uid"))
	authorRole := r.FormValue("role")
	now := time.Now()

	mid, err := repository.AddMessage(snapshotID, "", authorID, authorRole, now, 1)
	if err != nil {
		log.Fatal("Could not save feedback for error: ", err)
		return
	}

	// Use GetCodeSnapshot instead of raw SQL query
	codeSnapshot, err := repository.GetCodeSnapshot(snapshotID)
	if err != nil {
		log.Fatal("Failed to retrieve code snapshot: ", err)
		return
	}

	// Check if the author is trying to give feedback on their own code
	if authorRole == "student" && codeSnapshot.StudentID == authorID {
		fmt.Fprintf(w, "You can not give feedback to your own code.")
		return
	}

	// Update feedback count
	idx := models.StudentSnapshot[codeSnapshot.StudentID][codeSnapshot.ProblemID]
	models.Snapshots[idx].NumFeedback++

	// Add feedback message
	messageID := mid
	id, err := repository.AddMessageFeedback(messageID, feedback, authorID, authorRole, now)
	if err != nil {
		log.Fatal("Could not save feedback for error: ", err)
		return
	}
	feedbackID := id

	// Append feedback to the student's queue
	models.Students[codeSnapshot.StudentID].SnapShotFeedbackQueue = append(models.Students[codeSnapshot.StudentID].SnapShotFeedbackQueue, &models.SnapShotFeedback{
		FeedbackID:  int(feedbackID),
		Snapshot:    codeSnapshot.Code,
		Feedback:    feedback,
		ProblemName: codeSnapshot.Event, // Assuming Event holds the problem name
		Provider:    getName(uid, authorRole),
	})

	// Update student status
	repository.AddOrUpdateStudentStatus(codeSnapshot.StudentID, codeSnapshot.ProblemID, "", "Been helped", "", "")
	if authorRole == "student" {
		repository.AddOrUpdateStudentStatus(authorID, codeSnapshot.ProblemID, "", "", "", "Tutoring")
	}
	fmt.Println("Feedback on code snapshot saved!")
}

func messageFeedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	messageID, _ := strconv.Atoi(r.FormValue("message_id"))
	feedback := r.FormValue("feedback")
	authorID, _ := strconv.Atoi(r.FormValue("uid"))
	authorRole := r.FormValue("role")
	now := time.Now()

	id, err := repository.AddMessageFeedback(messageID, feedback, authorID, authorRole, now)
	if err != nil {
		log.Fatal("Could not save the feedback for error: ", err)
		return
	}

	details, err := repository.GetCodeSnapshotMessageDetails(messageID)
	if err != nil {
		log.Fatal(err)
	}

	// Process the details (assuming the query returned a single result)
	if len(details) > 0 {
		detail := details[0]
		studentID := detail.StudentID
		code := detail.Code
		filename := detail.Filename
		messageType := detail.MessageType

		if authorRole == "student" && studentID == authorID {
			fmt.Fprintf(w, "You cannot give feedback to your own code.")
			return
		}

		if messageType == "0" {
			repository.AddOrUpdateStudentStatus(studentID, detail.ProblemID, "", "Been helped", "", "")
			if authorRole == "student" {
				repository.AddOrUpdateStudentStatus(authorID, detail.ProblemID, "", "", "", "Tutoring")
			}
		}

		idx := models.StudentSnapshot[studentID][detail.ProblemID]
		models.Snapshots[idx].NumFeedback++
		models.Students[studentID].SnapShotFeedbackQueue = append(models.Students[studentID].SnapShotFeedbackQueue, &models.SnapShotFeedback{
			FeedbackID:  id,
			Snapshot:    code,
			Feedback:    feedback,
			ProblemName: filename,
			Provider:    getName(uid, authorRole),
		})
		fmt.Println("Feedback on message saved!")
	}
}

func getSnapshotFeedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	feedfback := models.Students[uid].SnapShotFeedbackQueue[0]
	models.Students[uid].SnapShotFeedbackQueue = models.Students[uid].SnapShotFeedbackQueue[1:]
	js, err := json.Marshal(feedfback)
	if err != nil {
		log.Fatal(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

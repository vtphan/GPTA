package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func addCodeSnapshot(studentID int, problemID int, code string, status int, lastUpdate time.Time, event string) int {
	result, err := AddCodeSnapshot(studentID, problemID, code, status, lastUpdate, event)
	if err != nil {
		log.Fatal("Could not save the snapshot for error: ", err)
		return -1
	}
	snapshotID := result.ID
	idx, ok := StudentSnapshot[studentID][problemID]
	if !ok {
		idx = len(Snapshots)
		StudentSnapshot[studentID][problemID] = idx
		name := GetStudentName(studentID)
		if err != nil {
			log.Fatal("Could not retrieve student name: ", err)
			return -1
		}
		problemName := ""
		for _, problem := range ActiveProblems {
			if problem.Active == true && problem.Info.Pid == problemID {
				problemName = problem.Info.Filename
				break
			}
		}
		Snapshots = append(Snapshots, &Snapshot{
			ID:          int(snapshotID),
			StudentName: name,
			StudentID:   studentID,
			ProblemName: problemName,
			ProblemID:   problemID,
			Status:      status,
			FirstUpdate: lastUpdate,
			LastUpdated: lastUpdate,
			LinesOfCode: getLinesOfCode(code),
			Code:        code,
		})
	} else {
		currentStatus := Snapshots[idx].Status
		if currentStatus > status {
			status = currentStatus
		}
		Snapshots[idx] = &Snapshot{
			ID:          int(snapshotID),
			StudentName: Snapshots[idx].StudentName,
			StudentID:   studentID,
			ProblemName: Snapshots[idx].ProblemName,
			ProblemID:   problemID,
			Status:      status,
			FirstUpdate: Snapshots[idx].FirstUpdate,
			LastUpdated: lastUpdate,
			LinesOfCode: getLinesOfCode(code),
			Code:        code,
			NumFeedback: Snapshots[idx].NumFeedback,
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

	mid, err := AddMessage(snapshotID, "", authorID, authorRole, now, 1)
	if err != nil {
		log.Fatal("Could not save feedback for error: ", err)
		return
	}

	// Use GetCodeSnapshot instead of raw SQL query
	codeSnapshot, err := GetCodeSnapshot(snapshotID)
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
	idx := StudentSnapshot[codeSnapshot.StudentID][codeSnapshot.ProblemID]
	Snapshots[idx].NumFeedback++

	// Add feedback message
	messageID := mid
	id, err := AddMessageFeedback(messageID, feedback, authorID, authorRole, now)
	if err != nil {
		log.Fatal("Could not save feedback for error: ", err)
		return
	}
	feedbackID := id

	// Append feedback to the student's queue
	Students[codeSnapshot.StudentID].SnapShotFeedbackQueue = append(Students[codeSnapshot.StudentID].SnapShotFeedbackQueue, &SnapShotFeedback{
		FeedbackID:  int(feedbackID),
		Snapshot:    codeSnapshot.Code,
		Feedback:    feedback,
		ProblemName: codeSnapshot.Event, // Assuming Event holds the problem name
		Provider:    getName(uid, authorRole),
	})

	// Update student status
	addOrUpdateStudentStatus(codeSnapshot.StudentID, codeSnapshot.ProblemID, "", "Been helped", "", "")
	if authorRole == "student" {
		addOrUpdateStudentStatus(authorID, codeSnapshot.ProblemID, "", "", "", "Tutoring")
	}
	fmt.Println("Feedback on code snapshot saved!")
}

func messageFeedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	messageID, _ := strconv.Atoi(r.FormValue("message_id"))
	feedback := r.FormValue("feedback")
	authorID, _ := strconv.Atoi(r.FormValue("uid"))
	authorRole := r.FormValue("role")
	now := time.Now()

	id, err := AddMessageFeedback(messageID, feedback, authorID, authorRole, now)
	if err != nil {
		log.Fatal("Could not save the feedback for error: ", err)
		return
	}

	details, err := GetCodeSnapshotMessageDetails(messageID)
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

		if studentID == authorID {
			fmt.Fprintf(w, "You cannot give feedback to your own code.")
			return
		}

		if messageType == "0" {
			addOrUpdateStudentStatus(studentID, detail.ProblemID, "", "Been helped", "", "")
			if authorRole == "student" {
				addOrUpdateStudentStatus(authorID, detail.ProblemID, "", "", "", "Tutoring")
			}
		}

		idx := StudentSnapshot[studentID][detail.ProblemID]
		Snapshots[idx].NumFeedback++
		Students[studentID].SnapShotFeedbackQueue = append(Students[studentID].SnapShotFeedbackQueue, &SnapShotFeedback{
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
	feedfback := Students[uid].SnapShotFeedbackQueue[0]
	Students[uid].SnapShotFeedbackQueue = Students[uid].SnapShotFeedbackQueue[1:]
	js, err := json.Marshal(feedfback)
	if err != nil {
		log.Fatal(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

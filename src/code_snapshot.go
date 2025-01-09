package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// test
func addCodeSnapshot(studentID int, problemID int, code string, status int, lastUpdate time.Time, event string) int {
	// Add the snapshot to the database
	snapshot := &CodeSnapshot{
		StudentID:     studentID,
		ProblemID:     problemID,
		Status:        status,
		LastUpdatedAt: lastUpdate,
		Code:          code,
		Event:         event,
	}

	if err := DB.Create(snapshot).Error; err != nil {
		log.Fatal("Could not save the snapshot for error: ", err)
		return -1
	}

	// Get the snapshot ID from the newly created record
	snapshotID := snapshot.ID

	// Check if the student-problem mapping exists in `StudentSnapshot`
	idx, ok := StudentSnapshot[studentID][problemID]
	if !ok {
		// Initialize the index and student-problem mapping
		idx = len(Snapshots)
		StudentSnapshot[studentID][problemID] = idx

		// Get the student's name using GORM
		var student Student
		if err := DB.First(&student, studentID).Error; err != nil {
			log.Fatal("Error fetching student data: ", err)
			return -1
		}

		// Get the problem name from the active problems
		problemName := ""
		for _, problem := range ActiveProblems {
			if problem.Active && problem.Info.Pid == problemID {
				problemName = problem.Info.Filename
				break
			}
		}

		// Append the new snapshot to the `Snapshots` array
		Snapshots = append(Snapshots, &Snapshot{
			ID:          snapshotID,
			StudentName: student.Name,
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
		// Update the existing snapshot
		currentStatus := Snapshots[idx].Status
		if currentStatus > status {
			status = currentStatus
		}
		Snapshots[idx] = &Snapshot{
			ID:          snapshotID,
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
	// Parse form values
	snapshotID, _ := strconv.Atoi(r.FormValue("snapshot_id"))
	feedback := r.FormValue("feedback")
	authorID, _ := strconv.Atoi(r.FormValue("uid"))
	authorRole := r.FormValue("role")
	now := time.Now()

	// Retrieve snapshot details
	var snapshot CodeSnapshot
	err := DB.Where("id = ?", snapshotID).First(&snapshot).Error
	if err != nil {
		log.Fatal("Error fetching code snapshot: ", err)
		http.Error(w, "Could not find snapshot", http.StatusNotFound)
		return
	}

	// Prevent students from giving feedback on their own code
	if authorRole == "student" && snapshot.StudentID == authorID {
		fmt.Fprintf(w, "You cannot give feedback to your own code.")
		return
	}

	// Insert the feedback
	newFeedback := SnapshotFeedback{
		SnapshotID: snapshotID,
		Feedback:   feedback,
		AuthorID:   authorID,
		AuthorRole: authorRole,
		GivenAt:    now,
	}
	if err := DB.Create(&newFeedback).Error; err != nil {
		log.Fatal("Could not save feedback: ", err)
		http.Error(w, "Could not save feedback", http.StatusInternalServerError)
		return
	}

	// Update the number of feedbacks for the snapshot
	idx := StudentSnapshot[snapshot.StudentID][snapshot.ProblemID]
	Snapshots[idx].NumFeedback++

	// Add feedback to student's snapshot feedback queue
	studentName := getName(uid, authorRole)
	Students[snapshot.StudentID].SnapShotFeedbackQueue = append(Students[snapshot.StudentID].SnapShotFeedbackQueue, &SnapShotFeedback{
		FeedbackID:  newFeedback.ID,
		Snapshot:    snapshot.Code,
		Feedback:    feedback,
		ProblemName: "Problem Name Placeholder", // Replace with actual problem name if available
		Provider:    studentName,
	})

	// Update student status
	addOrUpdateStudentStatus(snapshot.StudentID, snapshot.ProblemID, "", "Been helped", "", "")
	if authorRole == "student" {
		addOrUpdateStudentStatus(authorID, snapshot.ProblemID, "", "", "", "Tutoring")
	}

	// Respond with success
	fmt.Println("Feedback on code snapshot saved!")
	fmt.Fprintf(w, "Feedback saved successfully!")
}

type SnapshotMessage struct {
	StudentID   int
	ProblemID   int
	Code        string
	Filename    string
	MessageType int
}

func messageFeedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	// Parse form values
	messageID, _ := strconv.Atoi(r.FormValue("message_id"))
	feedback := r.FormValue("feedback")
	authorID, _ := strconv.Atoi(r.FormValue("uid"))
	authorRole := r.FormValue("role")
	now := time.Now()

	// Insert the feedback using GORM
	newFeedback := SnapshotFeedback{
		SnapshotID: messageID, // assuming SnapshotID corresponds to messageID in the context
		Feedback:   feedback,
		AuthorID:   authorID,
		AuthorRole: authorRole,
		GivenAt:    now,
	}
	if err := DB.Create(&newFeedback).Error; err != nil {
		log.Fatal("Could not save feedback: ", err)
		http.Error(w, "Could not save feedback", http.StatusInternalServerError)
		return
	}

	// Retrieve snapshot details based on messageID
	var snapshotDetails SnapshotMessage
	err := DB.Table("code_snapshot").
		Select("cs.student_id, cs.problem_id, cs.code, p.filename, m.type").
		Joins("join problem p on cs.problem_id = p.id").
		Joins("join message m on m.snapshot_id = cs.id").
		Where("m.id = ?", messageID).
		Scan(&snapshotDetails).Error
	if err != nil {
		log.Fatal("Error fetching message details: ", err)
		http.Error(w, "Could not find message details", http.StatusNotFound)
		return
	}

	// Prevent users from giving feedback to their own code
	if snapshotDetails.StudentID == authorID {
		fmt.Fprintf(w, "You cannot give feedback to your own code.")
		return
	}

	// Update student status if message type is 0
	if snapshotDetails.MessageType == 0 {
		addOrUpdateStudentStatus(snapshotDetails.StudentID, snapshotDetails.ProblemID, "", "Been helped", "", "")
		if authorRole == "student" {
			addOrUpdateStudentStatus(authorID, snapshotDetails.ProblemID, "", "", "", "Tutoring")
		}
	}

	// Update the number of feedbacks for the snapshot
	idx := StudentSnapshot[snapshotDetails.StudentID][snapshotDetails.ProblemID]
	Snapshots[idx].NumFeedback++

	// Append feedback to student's snapshot feedback queue
	feedbackID := newFeedback.ID
	studentName := getName(uid, authorRole)
	Students[snapshotDetails.StudentID].SnapShotFeedbackQueue = append(Students[snapshotDetails.StudentID].SnapShotFeedbackQueue, &SnapShotFeedback{
		FeedbackID:  int(feedbackID),
		Snapshot:    snapshotDetails.Code,
		Feedback:    feedback,
		ProblemName: snapshotDetails.Filename,
		Provider:    studentName,
	})

	// Respond with success
	fmt.Println("Feedback on message saved!")
	fmt.Fprintf(w, "Feedback saved successfully!")
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

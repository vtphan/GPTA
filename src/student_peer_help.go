// Author: Vinhthuy Phan, 2018
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

// -----------------------------------------------------------------------------------
func studentGetHelpCode(w http.ResponseWriter, r *http.Request, who string, uid int) {
	filename := r.FormValue("filename")
	pid := 0
	prob, ok := ActiveProblems[filename]
	HelpSubSem.Lock()
	defer HelpSubSem.Unlock()
	selected := &HelpSubmission{}
	selected.Status = 1
	if ok {
		if prob.Active {
			// fmt.Fprint(w, "This problem is not active.")
			pid = prob.Info.Pid
			if _, ok := HelpEligibleStudents[pid][uid]; ok {

				for idx, sub := range WorkingHelpSubs {
					if sub.Pid != pid || sub.Uid == uid {
						continue
					}
					if _, ok := SeenHelpSubmissions[uid][sub.Sid]; !ok {
						selected = sub
						WorkingHelpSubs = append(WorkingHelpSubs[:idx], WorkingHelpSubs[idx+1:]...)
						SeenHelpSubmissions[uid][sub.Sid] = true
						selected.Status = 0
						break
					}
				}

				// fmt.Fprintf(w, "You are elligible to help in this problem.")

			} else {
				selected.Status = 2
			}
		}
	}
	js, err := json.Marshal(selected)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}

}

//-----------------------------------------------------------------------------------

func student_return_without_feedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	HelpSubSem.Lock()
	defer HelpSubSem.Unlock()
	tmp := r.FormValue("submission_id")
	submissionID, _ := strconv.Atoi(tmp)
	submission := HelpSubmissions[submissionID]
	WorkingHelpSubs = append(WorkingHelpSubs, submission)
	fmt.Fprint(w, "No feedback is given. This request is returned to the help queue.")
}

func student_send_help_messageHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	submissionID, _ := strconv.Atoi(r.FormValue("submission_id"))
	message := r.FormValue("message")

	// Use the GORM function to insert the help message
	err := AddHelpMessage(submissionID, uid, message, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	// Fetch the help submission for the given submissionID
	helpSub := HelpSubmissions[submissionID]
	studentID := helpSub.Uid

	// Compose the final message
	message = helpSub.Content + "\n\nFeedback: " + message

	// Create a new board entry for the student
	b := &Board{
		Content:      message,
		Answer:       "",
		Attempts:     0,
		Filename:     "peer_feedback.txt",
		Pid:          submissionID, // Assuming the submission ID maps to the Pid
		StartingTime: time.Now(),
		Type:         "peer_feedback",
	}

	// Append the new board entry to the student's boards
	Students[studentID].Boards = append(Students[studentID].Boards, b)

	// Send a response back to the client
	fmt.Fprint(w, "Dear "+who+", Your feedback has been sent.")
}

func sendThankYouHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	messageID, _ := strconv.Atoi(r.FormValue("message_id"))
	useful := r.FormValue("useful") == "yes" // Convert to a boolean

	// Use the GORM function to update the HelpMessage
	err := UpdateHelpMessage(messageID, useful, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	// If the message is marked as useful, update the student's thank status
	if useful {
		var studentID int
		err := Database.Model(&HelpMessage{}). // Use the HelpMessage model
							Where("id = ?", messageID).
							Pluck("student_id", &studentID). // Retrieve the student_id
							Error
		if err != nil {
			log.Fatal(err)
		}
		Students[studentID].ThankStatus = 1
	}

	// Respond back to the client
	fmt.Fprintf(w, "Thank you, %s! Your feedback has been recorded.", who)
}

func studentSendBackFeedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	backFeedback := r.FormValue("feedback")
	feedbackID, _ := strconv.Atoi(r.FormValue("feedback_id"))
	authorRole := r.FormValue("role")
	useful := r.FormValue("useful") // Assuming this field contains "yes" or "no" to determine if feedback is useful

	// Convert the "useful" string to a boolean
	var isUseful bool
	if useful == "yes" {
		isUseful = true
	} else if useful == "no" {
		isUseful = false
	} else {
		// Handle the case where "useful" is neither "yes" nor "no"
		http.Error(w, "Invalid value for 'useful'", http.StatusBadRequest)
		return
	}

	// Check if the feedback already exists for the given feedbackID, uid, and authorRole
	var existingFeedback MessageBackFeedback
	err := Database.Where("message_feedback_id = ? AND author_id = ? AND author_role = ?", feedbackID, uid, authorRole).First(&existingFeedback).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Fatal(err)
	}

	if err == nil { // Feedback exists, update it
		// Use the GORM function to update the feedback
		err := UpdateMessageBackFeedback(feedbackID, isUseful, time.Now())
		if err != nil {
			log.Fatal(err)
		}
	} else { // Feedback does not exist, create new feedback
		// Use the GORM function to add new feedback
		err := AddMessageBackFeedback(feedbackID, uid, authorRole, backFeedback, time.Now())
		if err != nil {
			log.Fatal(err)
		}
	}

	// Send a response back to the client
	fmt.Fprintf(w, "Your feedback has been recorded, %s.", who)
}

// Author: Vinhthuy Phan, 2018
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

// -----------------------------------------------------------------------------------
func studentGetHelpCode(w http.ResponseWriter, r *http.Request, who string, uid int) {
	filename := r.FormValue("filename")
	pid := 0
	prob, ok := models.ActiveProblems[filename]
	models.HelpSubSem.Lock()
	defer models.HelpSubSem.Unlock()
	selected := &models.HelpSubmission{}
	selected.Status = 1
	if ok {
		if prob.Active {
			// fmt.Fprint(w, "This problem is not active.")
			pid = prob.Info.Pid
			if _, ok := models.HelpEligibleStudents[pid][uid]; ok {

				for idx, sub := range models.WorkingHelpSubs {
					if sub.Pid != pid || sub.Uid == uid {
						continue
					}
					if _, ok := models.SeenHelpSubmissions[uid][sub.Sid]; !ok {
						selected = sub
						models.WorkingHelpSubs = append(models.WorkingHelpSubs[:idx], models.WorkingHelpSubs[idx+1:]...)
						models.SeenHelpSubmissions[uid][sub.Sid] = true
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
	models.HelpSubSem.Lock()
	defer models.HelpSubSem.Unlock()
	tmp := r.FormValue("submission_id")
	submissionID, _ := strconv.Atoi(tmp)
	submission := models.HelpSubmissions[submissionID]
	models.WorkingHelpSubs = append(models.WorkingHelpSubs, submission)
	fmt.Fprint(w, "No feedback is given. This request is returned to the help queue.")
}

func student_send_help_messageHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	submissionID, _ := strconv.Atoi(r.FormValue("submission_id"))
	message := r.FormValue("message")
	res, err := repository.AddHelpMessage(submissionID, uid, message, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	messageID := res.ID // todo test this ID coming correctly
	helpSub := models.HelpSubmissions[submissionID]
	studentID := helpSub.Uid
	message = helpSub.Content + "\n\nFeedback: " + message
	b := &models.Board{
		Content:      message,
		Answer:       "",
		Attempts:     0,
		Filename:     "peer_feedback.txt",
		Pid:          int(messageID),
		StartingTime: time.Now(),
		Type:         "peer_feedback",
	}
	models.Students[studentID].Boards = append(models.Students[studentID].Boards, b)
	fmt.Fprint(w, "Dear "+who+", Your feedback has been sent.")

}
func sendThankYouHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	messageID, _ := strconv.Atoi(r.FormValue("message_id"))
	useful := r.FormValue("useful")
	err := repository.UpdateHelpMessage(useful, time.Now(), messageID)
	if err != nil {
		log.Fatal(err)
	}
	if useful == "yes" {
		studentID, err := repository.GetStudentIDByMessageID(messageID)
		if err != nil {
			log.Printf("Error retrieving student ID: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Update the in-memory Students map
		if student, exists := models.Students[studentID]; exists {
			student.ThankStatus = 1
		}
	}

}

func studentSendBackFeedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	backFeedback := r.FormValue("feedback")
	feedbackID, _ := strconv.Atoi(r.FormValue("feedback_id"))
	authorRole := r.FormValue("role")

	existingFeedback, err := repository.FetchExistingMessageBackFeedback(feedbackID, uid, authorRole)
	if err != nil {
		log.Fatal(err)
	}

	if existingFeedback != nil {
		err = repository.UpdateMessageBackFeedback(backFeedback, time.Now(), feedbackID, uid, authorRole)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = repository.AddMessageBackFeedback(feedbackID, uid, authorRole, backFeedback)
		if err != nil {
			log.Fatal(err)
		}
	}
}

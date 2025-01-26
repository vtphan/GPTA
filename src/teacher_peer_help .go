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
func teacherGetHelpCode(w http.ResponseWriter, r *http.Request, who string, uid int) {

	models.HelpSubSem.Lock()
	defer models.HelpSubSem.Unlock()
	selected := &models.HelpSubmission{}
	selected.Status = 1

	// fmt.Fprint(w, "This problem is not active.")

	for idx, sub := range models.WorkingHelpSubs {
		// if _, ok := SeenHelpSubmissions[uid][sub.Sid]; !ok {
		selected = sub
		models.WorkingHelpSubs = append(models.WorkingHelpSubs[:idx], models.WorkingHelpSubs[idx+1:]...)
		models.SeenHelpSubmissions[uid][sub.Sid] = true
		selected.Status = 0
		break
		// }
	}

	// fmt.Fprintf(w, "You are elligible to help in this problem.")

	js, err := json.Marshal(selected)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}

}

//-----------------------------------------------------------------------------------

func teacher_return_without_feedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	models.HelpSubSem.Lock()
	defer models.HelpSubSem.Unlock()
	tmp := r.FormValue("submission_id")
	submissionID, _ := strconv.Atoi(tmp)
	submission := models.HelpSubmissions[submissionID]
	models.WorkingHelpSubs = append(models.WorkingHelpSubs, submission)
	fmt.Fprint(w, "No feedback is given. This request is returned to the help queue.")
}

func teacher_send_help_messageHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	submission_id, _ := strconv.Atoi(r.FormValue("submission_id"))
	message := r.FormValue("message")
	_, err := repository.AddHelpMessage(submission_id, uid, message, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	helpSub := models.HelpSubmissions[submission_id]
	student_id := helpSub.Uid
	message = helpSub.Content + "\n\nFeedback: " + message
	b := &models.Board{
		Content:      message,
		Answer:       "",
		Attempts:     0,
		Filename:     "peer_feedback.txt",
		Pid:          0,
		StartingTime: time.Now(),
		Type:         "peer_feedback",
	}
	models.Students[student_id].Boards = append(models.Students[student_id].Boards, b)
	fmt.Fprint(w, "Your feedback has been sent.")

}

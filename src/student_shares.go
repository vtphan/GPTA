// Author: Vinhthuy Phan, 2018
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// -----------------------------------------------------------------------------------
func student_sharesHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	content, filename := r.FormValue("content"), r.FormValue("filename")
	answer := r.FormValue("answer")
	test_cases := r.FormValue("testcases")
	priority, _ := strconv.Atoi(r.FormValue("priority"))
	sid := 0
	correct_answer := ""
	complete := false
	var err error
	msg := "Your submission will be looked at soon."

	attempt_number := -1
	pid := 0
	prob, ok := ActiveProblems[filename]
	now := time.Now()
	snapshotID := -1
	if ok {
		if !prob.Active {
			msg = "Problem is no longer active. But the teacher will look at your submission."
		} else {
			pid = prob.Info.Pid
			if _, ok := prob.Attempts[uid]; !ok {
				ActiveProblems[filename].Attempts[uid] = prob.Info.Attempts
			}
			if ActiveProblems[filename].Attempts[uid] == 0 {
				fmt.Fprintf(w, "This is not submitted because either you have reached the submission limit or your solution was previously graded correctly.")
				return
			}

			// Decrement attempts **only if students are not asking for help**
			if priority < 2 {
				ActiveProblems[filename].Attempts[uid] -= 1
				if ActiveProblems[filename].Attempts[uid] <= 3 {
					msg += fmt.Sprintf(" You have %d attempt(s) left.", ActiveProblems[filename].Attempts[uid])
				}
			}
			attempt_number = prob.Info.Attempts - ActiveProblems[filename].Attempts[uid]
			// Autograding if possible
			correct_answer = ActiveProblems[filename].Info.Answer
			decision := ""
			addOrUpdateStudentStatus(uid, pid, "", "", "submitted", "")
			if answer != "" {
				scoring_mesg := ""
				if correct_answer == answer {
					decision = "correct"
					scoring_mesg = add_or_update_score("correct", pid, uid, 0, -1)
					ActiveProblems[filename].Attempts[uid] = 0 // This prevents further submission
					complete = true
					err = IncrementProblemStatGradedCorrect(pid)
					if err != nil {
						log.Fatal(err)
					}
					addOrUpdateStudentStatus(uid, pid, "", "", "Graded Correct", "")
				} else if ActiveProblems[filename].Info.ExactAnswer {
					decision = "incorrect"
					scoring_mesg = add_or_update_score("incorrect", pid, uid, 0, -1)
					complete = true
					err = IncrementProblemStatGradedIncorrect(pid)
					if err != nil {
						log.Fatal(err)
					}
					addOrUpdateStudentStatus(uid, pid, "", "", "Graded Incorrect", "")
				} else {
					scoring_mesg = "Answer appears to be incorrect. It will be looked at."
				}
				ActiveProblems[filename].Answers = append(ActiveProblems[filename].Answers, answer)

				fmt.Fprintf(w, scoring_mesg)
			}

			// Add submitted but not graded code to code snapshot.
			snapshotID = addCodeSnapshot(uid, pid, content, 1, now, "at_submission")

			var result SubmissionTable
			if complete {
				result, err = AddSubmissionComplete(pid, uid, content, priority, attempt_number, now, now, snapshotID, answer)
			} else {
				result, err = AddSubmission(pid, uid, content, priority, attempt_number, now, snapshotID, answer)
			}
			if err != nil {

				log.Fatal(err)
			}
			sid = result.ID

			err = IncrementProblemStatSubmission(pid)
			if err != nil {
				log.Fatal(err)
			}
			if complete {
				err := CompleteSubmission(time.Now(), decision, sid)
				if err != nil {
					log.Fatal(err)
				}
			}
			if test_cases != "" {
				existingTestCaseID, err := FetchExistingTestCase(uid, pid)
				if err != nil {
					log.Fatal(err)
				}

				if existingTestCaseID != 0 {
					err = UpdateTestCase(test_cases, now, existingTestCaseID)
				} else {
					_, err = AddTestCase(pid, uid, test_cases, now)
				}

				if err != nil {
					log.Fatal(err)
				}
			}
			if ActiveProblems[filename].Attempts[uid] == 0 {
				if PeerTutorAllowed {
					if _, ok := HelpEligibleStudents[pid][uid]; !ok {
						HelpEligibleStudents[pid][uid] = true
						SeenHelpSubmissions[uid] = map[int]bool{}
						// fmt.Fprintf(w, "You are now elligible to help you friends. To help please click on 'Help Friends' button.")
						msg = msg + "\nYou are now elligible to help you friends. To help please click on 'Help Friends' button."

						_, err = AddHelpEligible(pid, uid, now)
						if err != nil {
							log.Fatal(err)
						}
						addOrUpdateStudentStatus(uid, pid, "", "", "", "Qualified")
					}
				}
			}

		}
	}
	if !complete {
		SubSem.Lock()
		defer SubSem.Unlock()
		sub := &Submission{
			Sid:           int(sid),
			Uid:           uid,
			Pid:           pid,
			Content:       content,
			Filename:      filename,
			Priority:      priority,
			AttemptNumber: attempt_number,
			At:            now,
			Name:          r.FormValue("name"),
			SnapshotID:    snapshotID,
		}
		WorkingSubs = append(WorkingSubs, sub)
		Submissions[int(sid)] = sub
		fmt.Fprintf(w, msg)
	}
}

//-----------------------------------------------------------------------------------

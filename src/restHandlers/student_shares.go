// Author: Vinhthuy Phan, 2018
package restHandlers

import (
	"fmt"
	"github.com/GPTA/src/Grade"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"log"
	"net/http"
	"strconv"
	"time"
)

// -----------------------------------------------------------------------------------
func StudentSharesHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
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
	prob, ok := models.ActiveProblems[filename]
	now := time.Now()
	snapshotID := -1
	if ok {
		if !prob.Active {
			msg = "Problem is no longer active. But the teacher will look at your submission."
		} else {
			pid = prob.Info.Pid
			if _, ok := prob.Attempts[uid]; !ok {
				models.ActiveProblems[filename].Attempts[uid] = prob.Info.Attempts
			}
			if models.ActiveProblems[filename].Attempts[uid] == 0 {
				fmt.Fprintf(w, "This is not submitted because either you have reached the submission limit or your solution was previously graded correctly.")
				return
			}

			// Decrement attempts **only if students are not asking for help**
			if priority < 2 {
				models.ActiveProblems[filename].Attempts[uid] -= 1
				if models.ActiveProblems[filename].Attempts[uid] <= 3 {
					msg += fmt.Sprintf(" You have %d attempt(s) left.", models.ActiveProblems[filename].Attempts[uid])
				}
			}
			attempt_number = prob.Info.Attempts - models.ActiveProblems[filename].Attempts[uid]
			// Autograding if possible
			correct_answer = models.ActiveProblems[filename].Info.Answer
			decision := ""
			repository.AddOrUpdateStudentStatus(uid, pid, "", "", "submitted")
			if answer != "" && correct_answer != "" {
				scoring_mesg := ""
				if correct_answer == answer {
					decision = "correct"
					scoring_mesg = repository.AddOrUpdateScore("correct", pid, uid, 0, -1)
					models.ActiveProblems[filename].Attempts[uid] = 0 // This prevents further submission
					complete = true
					err = repository.IncrementProblemStatGradedCorrect(pid)
					if err != nil {
						log.Fatal(err)
					}
					repository.AddOrUpdateStudentStatus(uid, pid, "", "", "Graded Correct")
				} else if models.ActiveProblems[filename].Info.ExactAnswer {
					decision = "incorrect"
					scoring_mesg = repository.AddOrUpdateScore("incorrect", pid, uid, 0, -1)
					complete = true
					err = repository.IncrementProblemStatGradedIncorrect(pid)
					if err != nil {
						log.Fatal(err)
					}
					repository.AddOrUpdateStudentStatus(uid, pid, "", "", "Graded Incorrect")
				} else {
					scoring_mesg = "Answer appears to be incorrect. It will be looked at."
				}
				models.ActiveProblems[filename].Answers = append(models.ActiveProblems[filename].Answers, answer)
				var grade Grade.GradeRequest
				grade = Grade.GradeRequest{
					StudentID: uid,
					ProblemID: pid,
					Grade:     decision,
				}
				err = Grade.GradeStudent(grade)
				if err != nil {
					fmt.Println("Error in GradeStudent")
				}
				fmt.Fprintf(w, scoring_mesg)
			}

			// Add submitted but not graded code to code snapshot.
			snapshotID = AddCodeSnapshot(uid, pid, content, 1, now, "at_submission")

			var result models.SubmissionTable
			if complete {
				result, err = repository.AddSubmissionComplete(pid, uid, content, priority, attempt_number, now, now, snapshotID, answer)
			} else {
				result, err = repository.AddSubmission(pid, uid, content, priority, attempt_number, now, snapshotID, answer)
			}
			if err != nil {

				log.Fatal(err)
			}
			sid = result.ID

			err = repository.IncrementProblemStatSubmission(pid)
			if err != nil {
				log.Fatal(err)
			}
			if complete {
				err := repository.CompleteSubmission(time.Now(), decision, sid)
				if err != nil {
					log.Fatal(err)
				}
			}
			if test_cases != "" {
				existingTestCaseID, err := repository.FetchExistingTestCase(uid, pid)
				if err != nil {
					log.Fatal(err)
				}

				if existingTestCaseID != 0 {
					err = repository.UpdateTestCase(test_cases, now, existingTestCaseID)
				} else {
					_, err = repository.AddTestCase(pid, uid, test_cases, now)
				}

				if err != nil {
					log.Fatal(err)
				}
			}
			if models.ActiveProblems[filename].Attempts[uid] == 0 {
				if models.PeerTutorAllowed {
					if _, ok := models.HelpEligibleStudents[pid][uid]; !ok {
						models.HelpEligibleStudents[pid][uid] = true
						models.SeenHelpSubmissions[uid] = map[int]bool{}
						// fmt.Fprintf(w, "You are now elligible to help you friends. To help please click on 'Help Friends' button.")
						msg = msg + "\nYou are now elligible to help you friends. To help please click on 'Help Friends' button."

						_, err = repository.AddHelpEligible(pid, uid, now)
						if err != nil {
							log.Fatal(err)
						}
						repository.AddOrUpdateStudentStatus(uid, pid, "", "", "")
					}
				}
			}

		}
	}
	if !complete {
		models.SubSem.Lock()
		defer models.SubSem.Unlock()
		sub := &models.Submission{
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
		models.WorkingSubs = append(models.WorkingSubs, sub)
		models.Submissions[int(sid)] = sub
		fmt.Fprintf(w, msg)
	}
}

//-----------------------------------------------------------------------------------

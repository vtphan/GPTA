// Author: Vinhthuy Phan, 2018
package restHandlers

import (
	"fmt"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// -----------------------------------------------------------------------------------
func extract_partial_credits(content string) int {
	re := regexp.MustCompile(`(\d)+ for effort`)
	result := re.FindSubmatch([]byte(content))
	if len(result) >= 2 {
		points, _ := strconv.Atoi(string(result[1]))
		return points
	} else {
		return -1
	}
}

// -----------------------------------------------------------------------------------
func TeacherGradesHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	content, decision := r.FormValue("content"), r.FormValue("decision")
	sid, _ := strconv.Atoi(r.FormValue("sid"))
	changed := r.FormValue("changed")
	// student_id, _ := strconv.Atoi(r.FormValue("student_id"))
	// pid, _ := strconv.Atoi(r.FormValue("pid"))
	mesg := ""
	sub, ok := models.Submissions[sid]
	if !ok {
		fmt.Fprintf(w, "Unknown submission cannot be graded.")
		return
	}
	student_id := sub.Uid
	if changed == "True" {
		// If the original file is changed, there's feedback.  Copy it to whiteboard.
		if prob, ok := models.ActiveProblems[sub.Filename]; ok {
			_, err := repository.AddFeedback(uid, student_id, content, time.Now(), sub.Sid)
			if err != nil {
				log.Fatal(err)
			}
			mesg = "Feedback saved to student's board."
			models.BoardsSem.Lock()
			defer models.BoardsSem.Unlock()
			b := &models.Board{
				Content:      content,
				Answer:       prob.Info.Answer,
				Attempts:     0, // This tells the client this is an existing problem
				Filename:     sub.Filename,
				Pid:          sub.Pid,
				StartingTime: time.Now(),
				Type:         "feedback",
			}
			models.Students[student_id].Boards = append(models.Students[student_id].Boards, b)
		}
	}

	// If submission is dismissed, do not take that attempt away from the student.
	if decision == "dismissed" {
		// Students[student_id].SubmissionStatus = 2
		subStat := &models.StudentSubmissionStatus{
			Filename:      sub.Filename,
			AttemptNumber: sub.AttemptNumber,
			Status:        2,
		}
		models.Students[sub.Uid].SubmissionStatus = append(models.Students[sub.Uid].SubmissionStatus, subStat)

		models.ActiveProblems[sub.Filename].Attempts[student_id] += 1
		fmt.Fprintf(w, "Submission dismissed.")
	} else if decision == "ungraded" {
		// Students[student_id].SubmissionStatus = 5
		subStat := &models.StudentSubmissionStatus{
			Filename:      sub.Filename,
			AttemptNumber: sub.AttemptNumber,
			Status:        2,
		}
		models.Students[sub.Uid].SubmissionStatus = append(models.Students[sub.Uid].SubmissionStatus, subStat)

		fmt.Fprintf(w, mesg)
	} else {
		// Update score based on the grading decision
		partial_credits := -1
		if decision != "correct" {
			partial_credits = extract_partial_credits(content)
		}
		scoring_mesg := repository.AddOrUpdateScore(decision, sub.Pid, sub.Uid, uid, partial_credits)
		mesg = scoring_mesg + "\n" + mesg
		if decision == "correct" {
			// Students[sub.Uid].SubmissionStatus = 4
			subStat := &models.StudentSubmissionStatus{
				Filename:      sub.Filename,
				AttemptNumber: sub.AttemptNumber,
				Status:        4,
			}

			now := time.Now()
			models.ActiveProblems[sub.Filename].Attempts[student_id] = 0 // This prevents further submission.
			pid := sub.Pid
			if models.PeerTutorAllowed {
				if _, ok := models.HelpEligibleStudents[pid][sub.Uid]; !ok {
					models.HelpEligibleStudents[pid][sub.Uid] = true
					models.SeenHelpSubmissions[sub.Uid] = map[int]bool{}
					// Add eligible timestamp to datbase
					_, err := repository.AddHelpEligible(pid, sub.Uid, now)
					if err != nil {
						log.Fatal(err)
					}
					subStat.Status = 5
				}
			}
			models.Students[sub.Uid].SubmissionStatus = append(models.Students[sub.Uid].SubmissionStatus, subStat)

			// Add the correct submission to codesnapshot.
			// AddCodeSnapshot(sub.Uid, pid, content, 3, now)
			err := repository.IncrementProblemStatGradedCorrect(pid)
			if err != nil {
				log.Fatal(err)
			}
			repository.AddOrUpdateStudentStatus(sub.Uid, pid, "", "", "Graded Correct")

		} else {
			// Students[student_id].SubmissionStatus = 3
			subStat := &models.StudentSubmissionStatus{
				Filename:      sub.Filename,
				AttemptNumber: sub.AttemptNumber,
				Status:        3,
			}
			models.Students[sub.Uid].SubmissionStatus = append(models.Students[sub.Uid].SubmissionStatus, subStat)

			// Add the incorrect submission to codesnapshot.
			// AddCodeSnapshot(sub.Uid, sub.Pid, content, 2, time.Now())
			err := repository.IncrementProblemStatGradedIncorrect(sub.Pid)
			if err != nil {
				log.Fatal(err)
			}
			repository.AddOrUpdateStudentStatus(sub.Uid, sub.Pid, "", "", "Graded Incorrect")
		}

		// Update submission complete time
		err := repository.CompleteSubmission(time.Now(), decision, sid)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, mesg)
	}
}

//-----------------------------------------------------------------------------------

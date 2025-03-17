// Author: Vinhthuy Phan, 2018
package restHandlers

import (
	"fmt"
	"github.com/GPTA/src/frontEnd"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

// -----------------------------------------------------------------------------------
// func insert_problems(uid int, problems []*ProblemInfo) {
// -----------------------------------------------------------------------------------
func insert_problem(uid int, problem *models.ProblemInfo) {
	// Create new problem
	pid := int64(0)
	if problem.Merit > 0 {
		// Find Tag id
		tagID, err := repository.FetchTagIDByDescription(problem.Tag)
		if err != nil {
			log.Fatal(err)
		}
		if tagID == 0 {
			tag, err := repository.AddTag(problem.Tag)
			if err != nil {
				log.Fatal(err)
			}
			tagID = int64(tag.ID)
		}

		// Insert only real problems into database
		result, err := repository.AddProblem(
			uid,
			problem.Description,
			problem.Answer,
			problem.Filename,
			problem.Merit,
			problem.Effort,
			problem.Attempts,
			problem.Topic_id,
			int(tagID))
		if err != nil {
			log.Fatal(err)
		}
		problem.Pid = result.ID
		models.ActiveProblems[problem.Filename] = &models.ActiveProblem{
			Info:     problem,
			Answers:  make([]string, 0),
			Active:   true,
			Attempts: make(map[int]int),
		}
		models.HelpEligibleStudents[int(pid)] = map[int]bool{}
		err = repository.AddProblemStatistics(problem.Pid)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// -----------------------------------------------------------------------------------
// TeacherMap starts one or more problems.
// -----------------------------------------------------------------------------------
func TeacherBroadcastsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	content := r.FormValue("content")
	answer := r.FormValue("answer")
	merit, _ := strconv.Atoi(r.FormValue("merit"))
	effort, _ := strconv.Atoi(r.FormValue("effort"))
	attempts, _ := strconv.Atoi(r.FormValue("attempts"))
	topic_id, _ := strconv.Atoi(r.FormValue("topic_id"))
	tag := r.FormValue("tag")
	filename := r.FormValue("filename")
	exact_answer := r.FormValue("exact_answer")

	// fmt.Printf("%d,Answer:%s, Merit:%d, Effort:%d, Attempts:%d, Tag:%s, Filename:%s\n", len(content), answer, merit, effort, attempts, tag, filename)

	problem := &models.ProblemInfo{
		Description: content,
		Filename:    filename,
		Answer:      answer,
		Merit:       merit,
		Effort:      effort,
		Attempts:    attempts,
		Topic_id:    topic_id,
		Tag:         tag,
		ExactAnswer: exact_answer == "True",
	}
	// fmt.Println("answer:", problem.Answer, problem.ExactAnswer)

	insert_problem(uid, problem)
	models.BoardsSem.Lock()
	defer models.BoardsSem.Unlock()
	for student_id, _ := range models.Students {
		studentCourses, exists := models.StudentClassesMap[student_id]
		if !exists {
			continue // Skip students who are not mapped to any course
		}

		// Check if the student is mapped to models.CourseId
		isEnrolled := false
		for _, course := range studentCourses {
			if course == models.CourseId {
				isEnrolled = true
				break
			}
		}

		if !isEnrolled {
			continue // Skip students not enrolled in the current course
		}
		b := &models.Board{
			Content:      problem.Description,
			Answer:       problem.Answer,
			Attempts:     problem.Attempts,
			Filename:     problem.Filename,
			Pid:          problem.Pid,
			StartingTime: time.Now(),
			Type:         "new",
		}
		models.Students[student_id].Boards = append(models.Students[student_id].Boards, b)
		if b.Pid != 0 && student_id != 0 {
			// Add student coding status as idle
			repository.AddOrUpdateStudentStatus(student_id, b.Pid, "Idle", "", "")
		}
	}
	fmt.Fprintf(w, "Content copied to white boards.")
}

//-----------------------------------------------------------------------------------

func TeacherWebBroadcastHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	temp := template.New("")
	t, err := temp.Parse(frontEnd.PROBLEM_FILE_UPLOAD_VIEW)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

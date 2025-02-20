package restHandlers

import (
	"github.com/GPTA/src/frontEnd"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"html/template"
	"log"
	"net/http"
	"time"
)

type ProblemData struct {
	ID                 int
	Filename           string
	UploadedAt         time.Time
	IsActive           bool
	Attendance         int
	NumActive          int
	NumHelpRequest     int
	NumGradedCorrect   int
	NumGradedIncorrect int
	NumNotGraded       int
}

type ProblemListData struct {
	Problems         []*ProblemData
	PeerTutorAllowed bool
	UserID           int
	UserRole         string
	Password         string
	Username         string
}

func ProblemListHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	password := r.FormValue("password")

	var problemsFromDB []models.Problem
	var err error

	if role == "student" {
		courseId := r.FormValue("course_id")
		if courseId == "" {
			http.Error(w, "Missing course_id parameter", http.StatusBadRequest)
			return
		}

		problemsFromDB, err = repository.GetProblems(courseId)
		if err != nil {
			http.Error(w, "Failed to fetch problems", http.StatusInternalServerError)
			log.Printf("Error fetching problems for student: %v", err)
			return
		}
	}

	if role == "teacher" {
		problemsFromDB, err = repository.GetProblems(models.CourseId) // Adjust based on actual implementation
		if err != nil {
			http.Error(w, "Failed to fetch problems", http.StatusInternalServerError)
			log.Printf("Error fetching problems for teacher: %v", err)
			return
		}
	}
	// Convert database problems to ProblemData format
	var problems = make([]*ProblemData, 0)
	for _, problem := range problemsFromDB {
		// Get stats for each problem
		nActive, nHelp, nNotGraded, nCorrect, nIncorrect := repository.GetProblemStats(problem.ID)

		problems = append(problems, &ProblemData{
			ID:                 problem.ID,
			Filename:           problem.Filename,
			UploadedAt:         problem.ProblemUploadedAt,
			IsActive:           problem.ProblemEndedAt == nil,
			Attendance:         len(repository.GetCurrentStudents()),
			NumActive:          nActive,
			NumHelpRequest:     nHelp,
			NumGradedCorrect:   nCorrect,
			NumGradedIncorrect: nIncorrect,
			NumNotGraded:       nNotGraded,
		})
	}

	// Prepare the ProblemListData for the template
	problemListData := &ProblemListData{
		Problems:         problems,
		PeerTutorAllowed: models.PeerTutorAllowed,
		UserID:           uid,
		UserRole:         role,
		Password:         password,
		Username:         GetName(uid, role),
	}

	// Render the template
	temp := template.New("")
	t, err := temp.Parse(frontEnd.PROBLEM_LIST_TEMPLATE)
	if err != nil {
		log.Printf("Error parsing template: %v", err)
	}
	w.Header().Set("Content-Type", "text/html")
	if err = t.Execute(w, problemListData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
		return
	}
}

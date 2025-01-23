package main

import (
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

func problemListHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	password := r.FormValue("password")

	// Fetch problems using the GetProblems function
	problemsFromDB, err := GetProblems()
	if err != nil {
		log.Fatalf("Error fetching problems: %v", err)
	}

	// Convert database problems to ProblemData format
	var problems = make([]*ProblemData, 0)
	for _, problem := range problemsFromDB {
		// Get stats for each problem
		nActive, nHelp, nNotGraded, nCorrect, nIncorrect := GetProblemStats(problem.ID)

		problems = append(problems, &ProblemData{
			ID:                 problem.ID,
			Filename:           problem.Filename,
			UploadedAt:         problem.ProblemUploadedAt,
			IsActive:           problem.ProblemEndedAt == nil,
			Attendance:         len(GetCurrentStudents()),
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
		PeerTutorAllowed: PeerTutorAllowed,
		UserID:           uid,
		UserRole:         role,
		Password:         password,
		Username:         getName(uid, role),
	}

	// Render the template
	temp := template.New("")
	t, err := temp.Parse(PROBLEM_LIST_TEMPLATE)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}
	w.Header().Set("Content-Type", "text/html")
	if err = t.Execute(w, problemListData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatalf("Error executing template: %v", err)
	}
}

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

	// Fetch problems from the database using GORM
	var problemsData []Problem
	err := Database.Model(&Problem{}).
		Select("id, filename, problem_uploaded_at, problem_ended_at").
		Scan(&problemsData).Error
	if err != nil {
		log.Fatal(err)
	}

	problems := make([]*ProblemData, 0)

	for _, problem := range problemsData {
		// Gather stats for each problem
		nActive, nHelp, nNotGraded, nCorrect, nIncorrect := getProblemStats(problem.ID)

		problems = append(problems, &ProblemData{
			ID:                 problem.ID,
			Filename:           problem.Filename,
			UploadedAt:         problem.ProblemUploadedAt,
			IsActive:           problem.ProblemEndedAt == nil,
			Attendance:         len(getCurrentStudents()),
			NumActive:          nActive,
			NumHelpRequest:     nHelp,
			NumGradedCorrect:   nCorrect,
			NumGradedIncorrect: nIncorrect,
			NumNotGraded:       nNotGraded,
		})
	}

	// Prepare data for rendering
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
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, problemListData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

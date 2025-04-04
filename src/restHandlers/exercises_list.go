package restHandlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/GPTA/src/frontEnd"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
)

type problemData struct {
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
	LatestFeedbackTime *time.Time
}

type problemListData struct {
	Problems         []*problemData
	PeerTutorAllowed bool
	UserID           int
	UserRole         string
	Password         string
	Username         string
}

func ExerciseListHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	password := r.FormValue("password")
	courseId := r.FormValue("course_id")
	if courseId != "" {
		models.CourseId = courseId
	}

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
	var problems = make([]*problemData, 0)
	for _, problem := range problemsFromDB {
		// Get stats for each problem
		nActive, nHelp, nNotGraded, nCorrect, nIncorrect := repository.GetProblemStats(problem.ID)

		// Create the problem data structure
		problemData := &problemData{
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
			LatestFeedbackTime: nil, // Initialize to nil
		}

		students := repository.GetAllStudents()

		// For teachers, find the latest feedback time across all students for this problem
		// For students, just check their own feedback
		messages, err := FetchMessages(students, problem.ID, uid, role)
		if err != nil {
			log.Printf("Error fetching messages for problem %d, student %d: %v", problem.ID, uid, err)
		} else {
			var latestTime *time.Time
			for _, msg := range messages {
				for _, feedback := range msg.Feedbacks {
					if latestTime == nil || feedback.GivenAt.After(*latestTime) {
						latestTime = &feedback.GivenAt
					}
				}
			}
			problemData.LatestFeedbackTime = latestTime
		}

		problems = append(problems, problemData)
	}

	// Prepare the ProblemListData for the template
	problemListData := &problemListData{
		Problems:         problems,
		PeerTutorAllowed: models.PeerTutorAllowed,
		UserID:           uid,
		UserRole:         role,
		Password:         password,
		Username:         GetName(uid, role),
	}
	for i, problem := range problemListData.Problems {
		fmt.Printf("Problem %d: %+v\n", i+1, *problem)
	}
	// Render the template
	temp := template.New("")
	t, err := temp.Parse(frontEnd.EXERCISE_LIST_TEMPLATE)
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

// student's feedback page
func StudentFeedbackProvisionHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	studentID, _ := strconv.Atoi(r.FormValue("student_id"))
	temp := template.New("")
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode}
	t, err := temp.Funcs(ownFuncs).Parse(frontEnd.EXERCISE_FEEDBACK_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	students := repository.GetAllStudents()
	var messages = make([]*models.MessageDashBoard, 0)
	_, ok := models.HelpEligibleStudents[problemID][uid]
	if role == "teacher" || uid == studentID || (models.PeerTutorAllowed && ok) {
		// Fetch messages using the new function
		messages, err = FetchMessages(students, problemID, studentID, role)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.Error(w, "You are not authorized to access!", http.StatusUnauthorized)
	}

	// Sort the messages by descding order of GivenAt
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].GivenAt.After(messages[j].GivenAt)
	})

	// TODO(shiplu): sort the messages array
	// sort.Slice(helpRequests, func(i, j int) bool { return helpRequests[i].GivenAt.Before(helpRequests[j].GivenAt) })
	latestSnapshot := &models.Snapshot{}
	if _, ok := models.StudentSnapshot[studentID][problemID]; ok {
		latestSnapshot = models.Snapshots[models.StudentSnapshot[studentID][problemID]]
	} else {
		latestSnapshot, _ = repository.GetLatestSnapshot(studentID, problemID)
	}

	// Get student status
	studentStats := &models.DashBoardStudentInfo{}
	if role == "teacher" || uid == studentID || (models.PeerTutorAllowed && ok) {
		// Fetch student status using the new function
		studentStats, err = repository.FetchStudentStatuses(problemID, studentID)
		if err != nil {
			log.Fatal(err) // or use http.Error depending on your context
		}

		// Do something with studentStats, for example, returning them as JSON
		// Example: json.NewEncoder(w).Encode(studentStats)
	} else {
		http.Error(w, "You are not authorized to access!", http.StatusUnauthorized)
	}

	data := &models.FeedbackProvisionDashBoard{
		StudentName:  students[studentID],
		ProblemName:  latestSnapshot.ProblemName,
		LastSnapshot: latestSnapshot,
		Messages:     messages,
		StudentID:    studentID,
		ProblemID:    problemID,
		UserID:       uid,
		UserRole:     role,
		Password:     r.FormValue("password"),
		Status:       *studentStats,
		Username:     GetName(uid, role),
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

// func ListFeedbackForExercise()

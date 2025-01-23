package main

import (
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type DashBoardStudentInfo struct {
	StudentID      int
	StudentName    string
	LastUpdatedAt  time.Time
	CodingStat     string
	HelpStat       string
	SubmissionStat string
	TutoringStat   string
}

type AnswerStatInfo struct {
	Answer  string
	Count   int
	Percent float64
}

type DashBoardInfo struct {
	StudentInfo        []*DashBoardStudentInfo
	ProblemName        string
	Code               string
	IsActive           bool
	ProblemID          int
	NumActive          int
	NumHelpRequest     int
	NumGradedCorrect   int
	NumGradedIncorrect int
	NumNotGraded       int
	AnswerStats        []*AnswerStatInfo
	UserID             int
	UserRole           string
	Password           string
	Username           string
}

func getName(uid int, role string) string {
	name := ""
	if role == "teacher" {
		name = TeacherIdToName[uid]
	} else {
		name = Students[uid].Name
	}
	return name
}

func problemDashboardHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	role := r.FormValue("role")
	password := r.FormValue("password")

	// Retrieve all students
	students := GetAllStudents()

	// Retrieve code snapshots for the problem
	codeSnapshots, err := GetCodeSnapshotsByProblemID(problemID)
	if err != nil {
		http.Error(w, "Error fetching code snapshots", http.StatusInternalServerError)
		return
	}

	// Initialize map for storing the last update for each student
	lastUpdateMap := make(map[int]time.Time)

	_, ok := HelpEligibleStudents[problemID][uid]

	// Populate the map based on the role or eligibility
	for _, snapshot := range codeSnapshots {
		if role == "teacher" || uid == snapshot.StudentID || (PeerTutorAllowed && ok) {
			lastUpdateMap[snapshot.StudentID] = snapshot.LastUpdatedAt
		}
	}

	// Retrieve problem description and end time
	code, problemEndedAt, err := GetProblemDetail(problemID)
	if err != nil {
		http.Error(w, "Error fetching problem details", http.StatusInternalServerError)
		return
	}

	// Retrieve latest submission times for the problem
	latestSubmissionTime := GetLatestSubmissionTime(problemID)

	// Retrieve student statuses for the problem
	studentStatuses, err := GetStudentStatusesByProblemID(problemID)
	if err != nil {
		http.Error(w, "Error fetching student statuses", http.StatusInternalServerError)
		return
	}

	var studentInfo []*DashBoardStudentInfo
	for _, status := range studentStatuses {
		if role == "teacher" || uid == status.StudentID || (PeerTutorAllowed && ok) {
			studentInfo = append(studentInfo, &DashBoardStudentInfo{
				StudentID:      status.StudentID,
				StudentName:    students[status.StudentID],
				LastUpdatedAt:  lastUpdateMap[status.StudentID],
				CodingStat:     status.CodingStat,
				HelpStat:       status.HelpStat,
				SubmissionStat: status.SubmissionStat,
				TutoringStat:   status.TutoringStat,
			})
		}
	}

	sort.SliceStable(studentInfo, func(i, j int) bool {
		if studentInfo[i].SubmissionStat == "submitted" && studentInfo[j].SubmissionStat == "submitted" {
			return latestSubmissionTime[studentInfo[i].StudentID].Before(latestSubmissionTime[studentInfo[j].StudentID])
		}
		if studentInfo[i].SubmissionStat == "submitted" {
			return true
		}
		if studentInfo[j].SubmissionStat == "submitted" {
			return false
		}
		if studentInfo[i].HelpStat == "Asked for help" {
			return true
		}
		if studentInfo[j].HelpStat == "Asked for help" {
			return false
		}
		return true
	})

	nActive, nHelp, nNotGraded, nCorrect, nIncorrect := GetProblemStats(problemID)

	var answerStats []*AnswerStatInfo
	if role != "student" {
		answerStats, err = GetAnswerStats(problemID)
		if err != nil {
			http.Error(w, "Error fetching answer stats", http.StatusInternalServerError)
			return
		}
	}

	dashBoardData := &DashBoardInfo{
		StudentInfo:        studentInfo,
		ProblemID:          problemID,
		ProblemName:        GetProblemNameFromID(problemID),
		Code:               code,
		IsActive:           problemEndedAt.IsZero(),
		NumActive:          nActive,
		NumHelpRequest:     nHelp,
		NumGradedCorrect:   nCorrect,
		NumGradedIncorrect: nIncorrect,
		NumNotGraded:       nNotGraded,
		AnswerStats:        answerStats,
		UserID:             uid,
		UserRole:           role,
		Password:           password,
		Username:           getName(uid, role),
	}

	temp := template.New("")
	ownFuncs := template.FuncMap{"formatTimeSince": formatTimeSince}
	t, err := temp.Funcs(ownFuncs).Parse(PROBLEM_DASHBOARD_TEMPLATE)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, dashBoardData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

package restHandlers

import (
	"github.com/GPTA/src/frontEnd"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func GetName(uid int, role string) string {
	name := ""
	if role == "teacher" {
		name = models.TeacherIdToName[uid]
	} else {
		name = models.Students[uid].Name
	}
	return name
}

func ScaffoldingDashboardHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	// Parse query params
	problemID, err := strconv.Atoi(r.URL.Query().Get("problem_id"))
	if err != nil {
		http.Error(w, "Invalid problem_id", http.StatusBadRequest)
		return
	}

	// Fetch scaffoldings from DB
	var scaffoldings []models.Scaffolding
	result := models.DB.Where("problem_id = ?", problemID).Find(&scaffoldings)
	if result.Error != nil {
		http.Error(w, "Error fetching scaffoldings", http.StatusInternalServerError)
		return
	}

	// Organize scaffoldings by strategy and level
	scaffoldMap := make(map[int]map[int][]models.Scaffolding)
	for _, s := range scaffoldings {
		if _, exists := scaffoldMap[s.ScaffoldingStrategy]; !exists {
			scaffoldMap[s.ScaffoldingStrategy] = make(map[int][]models.Scaffolding)
		}
		scaffoldMap[s.ScaffoldingStrategy][s.ScaffoldingLevel] = append(scaffoldMap[s.ScaffoldingStrategy][s.ScaffoldingLevel], s)
	}

	// Template processing
	temp := template.New("")
	ownFuncs := template.FuncMap{}
	t, err := temp.Funcs(ownFuncs).Parse(frontEnd.SCAFFOLDING_TEMPLATE)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, scaffoldMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ProblemDashboardHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	role := r.FormValue("role")
	password := r.FormValue("password")

	// Retrieve all students
	students := repository.GetAllStudents()

	// Retrieve code snapshots for the problem
	codeSnapshots, err := repository.GetCodeSnapshotsByProblemID(problemID)
	if err != nil {
		http.Error(w, "Error fetching code snapshots", http.StatusInternalServerError)
		return
	}

	// Initialize map for storing the last update for each student
	lastUpdateMap := make(map[int]time.Time)

	_, ok := models.HelpEligibleStudents[problemID][uid]

	// Populate the map based on the role or eligibility
	for _, snapshot := range codeSnapshots {
		if role == "teacher" || uid == snapshot.StudentID || (models.PeerTutorAllowed && ok) {
			lastUpdateMap[snapshot.StudentID] = snapshot.LastUpdatedAt
		}
	}

	// Retrieve problem description and end time
	code, problemEndedAt, err := repository.GetProblemDetail(problemID)
	if err != nil {
		http.Error(w, "Error fetching problem details", http.StatusInternalServerError)
		return
	}

	// Retrieve latest submission times for the problem
	latestSubmissionTime := repository.GetLatestSubmissionTime(problemID)

	// Retrieve student statuses for the problem
	studentStatuses, err := repository.GetStudentStatusesByProblemID(problemID)
	if err != nil {
		http.Error(w, "Error fetching student statuses", http.StatusInternalServerError)
		return
	}

	var studentInfo []*models.DashBoardStudentInfo
	for _, status := range studentStatuses {
		if role == "teacher" || uid == status.StudentID || (models.PeerTutorAllowed && ok) {
			studentInfo = append(studentInfo, &models.DashBoardStudentInfo{
				StudentID:      status.StudentID,
				StudentName:    students[status.StudentID],
				LastUpdatedAt:  lastUpdateMap[status.StudentID],
				CodingStat:     status.CodingStat,
				HelpStat:       status.HelpStat,
				SubmissionStat: status.SubmissionStat,
				Percentage:     status.Percentage,
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

	nActive, nHelp, nNotGraded, nCorrect, nIncorrect := repository.GetProblemStats(problemID)

	var answerStats []*models.AnswerStatInfo
	if role != "student" {
		answerStats, err = repository.GetAnswerStats(problemID)
		if err != nil {
			http.Error(w, "Error fetching answer stats", http.StatusInternalServerError)
			return
		}
	}

	dashBoardData := &models.DashBoardInfo{
		StudentInfo:        studentInfo,
		ProblemID:          problemID,
		ProblemName:        repository.GetProblemNameFromID(problemID),
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
		Username:           GetName(uid, role),
	}

	temp := template.New("")
	ownFuncs := template.FuncMap{"formatTimeSince": FormatTimeSince}
	t, err := temp.Funcs(ownFuncs).Parse(frontEnd.PROBLEM_DASHBOARD_TEMPLATE)
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

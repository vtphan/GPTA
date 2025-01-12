package main

import (
	"html/template"
	"log"
	"math"
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

func getCurrentStudents() []int {
	var currentStudents []int
	err := Database.Model(&Attendance{}). // Use the Attendance model
						Select("student_id").                                              // Select only the student_id field
						Where("DATE(attendance_at) = ?", time.Now().Format("2006-01-02")). // Filter by today's date
						Scan(&currentStudents).Error                                       // Scan the results into currentStudents slice

	if err != nil {
		log.Fatal(err)
	}

	return currentStudents
}

func getAllStudents() map[int]string {
	var students []struct {
		ID   int
		Name string
	}
	err := Database.Model(&Student{}).Select("id, name").Scan(&students).Error
	if err != nil {
		log.Fatal(err)
	}

	studentsMap := make(map[int]string)
	for _, student := range students {
		studentsMap[student.ID] = student.Name
	}
	return studentsMap
}

func getProblemStats(problemID int) (int, int, int, int, int) {
	var stats struct {
		Active          int
		Submission      int
		HelpRequest     int
		GradedCorrect   int
		GradedIncorrect int
	}

	err := Database.Model(&ProblemStatistics{}). // Use the ProblemStatistics model
							Select("active, submission, help_request, graded_correct, graded_incorrect"). // Select fields
							Where("problem_id = ?", problemID).                                           // Filter by problem_id
							Scan(&stats).Error                                                            // Scan the result into the stats struct

	if err != nil {
		log.Fatal(err)
	}

	// Return the values, calculating the ungraded submissions
	return stats.Active, stats.HelpRequest, stats.Submission - stats.GradedCorrect - stats.GradedIncorrect, stats.GradedCorrect, stats.GradedIncorrect
}

func getProblemNameFromID(problemID int) string {
	var problem Problem
	err := Database.Where("id = ?", problemID).First(&problem).Error
	if err != nil {
		log.Fatal(err)
	}
	return problem.Filename
}

func getLatestSubmissionTime(problemID int) map[int]time.Time {
	latestSubmissions := make(map[int]time.Time)

	// Assuming Submission is a model representing the submissions table
	var submissions []struct {
		StudentID      int
		SubmissionTime time.Time
	}

	// Perform the query using GORM's ORM methods
	err := Database.Model(&Submission{}). // Use the Submission model
						Select("student_id, MAX(code_submitted_at) AS submission_time"). // Select student_id and the max submission time
						Where("problem_id = ?", problemID).                              // Filter by problem_id
						Group("student_id").                                             // Group by student_id to get the latest submission per student
						Scan(&submissions).Error                                         // Scan the result into the submissions slice

	if err != nil {
		log.Fatal(err)
	}

	// Map the results to student ID and submission time
	for _, submission := range submissions {
		latestSubmissions[submission.StudentID] = submission.SubmissionTime
	}

	return latestSubmissions
}

func problemDashboardHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	role := r.FormValue("role")
	password := r.FormValue("password")

	students := getAllStudents()
	var lastUpdateMap = make(map[int]time.Time)

	// Fetch last update times for students using GORM ORM
	var snapshots []struct {
		StudentID     int
		LastUpdatedAt time.Time
	}
	err := Database.Model(&CodeSnapshot{}). // Use the CodeSnapshot model
						Select("student_id, max(last_updated_at) as last_updated_at").
						Where("problem_id = ?", problemID).
						Group("student_id").
						Scan(&snapshots).Error
	if err != nil {
		log.Fatal(err)
	}

	_, ok := HelpEligibleStudents[problemID][uid]

	for _, snapshot := range snapshots {
		if role == "teacher" || uid == snapshot.StudentID || (PeerTutorAllowed && ok) {
			lastUpdateMap[snapshot.StudentID] = snapshot.LastUpdatedAt
		}
	}

	// Fetch problem details using GORM ORM
	var problem Problem
	err = Database.Where("id = ?", problemID).First(&problem).Error
	if err != nil {
		log.Fatal(err)
	}

	latestSubmissionTime := getLatestSubmissionTime(problemID)

	// Fetch student statuses using GORM ORM
	var studentStatuses []struct {
		StudentID      int
		CodingStat     string
		HelpStat       string
		SubmissionStat string
		TutoringStat   string
	}
	err = Database.Model(&StudentStatus{}). // Use the StudentStatus model
						Select("student_id, coding_stat, help_stat, submission_stat, tutoring_stat").
						Where("problem_id = ?", problemID).
						Scan(&studentStatuses).Error
	if err != nil {
		log.Fatal(err)
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

	// Sort student info
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

	// Fetch problem stats using GORM ORM
	nActive, nHelp, nNotGraded, nCorrect, nIncorrect := getProblemStats(problemID)

	// Fetch answer stats using GORM ORM (if role is not student)
	var answerStats []*AnswerStatInfo
	if role != "student" {
		var answers []struct {
			Answer string
			Count  int
		}
		err := Database.Model(&Submission{}). // Use the Submission model
							Select("answer, count(*) as cnt").
							Where("problem_id = ? AND answer IS NOT NULL AND LENGTH(answer) > 0", problemID).
							Group("answer").
							Scan(&answers).Error
		if err != nil {
			log.Fatal(err)
		}

		var total int
		for _, answer := range answers {
			answerStats = append(answerStats, &AnswerStatInfo{
				Answer: answer.Answer,
				Count:  answer.Count,
			})
			total += answer.Count
		}

		// Calculate percentage for each answer
		for i, answer := range answerStats {
			answerStats[i].Percent = math.Round(float64(answer.Count)*10000/float64(total)) / 100
		}
	}

	// Prepare dashboard data
	dashBoardData := &DashBoardInfo{
		StudentInfo:        studentInfo,
		ProblemID:          problemID,
		ProblemName:        getProblemNameFromID(problemID),
		Code:               problem.Filename,
		IsActive:           problem.ProblemEndedAt == nil || problem.ProblemEndedAt.IsZero(),
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

	// Render the template
	temp := template.New("")
	ownFuncs := template.FuncMap{"formatTimeSince": formatTimeSince}
	t, err := temp.Funcs(ownFuncs).Parse(PROBLEM_DASHBOARD_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, dashBoardData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

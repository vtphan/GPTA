package main

import (
	"fmt"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type FeedbackDashBaord struct {
	Name            string
	Role            string
	Feedback        string
	FeedbackID      int
	CurrentUserVote string
	Downvote        int
	Upvote          int
	GivenAt         time.Time
}

type MessageDashBoard struct {
	ID         int
	Name       string
	Role       string
	Message    string
	Type       int // 0 = help request, 1 = unsolicited
	Event      string
	GivenAt    time.Time
	Code       string
	SnapshotID int
	Feedbacks  []*FeedbackDashBaord
}

type FeedbackProvisionDashBoard struct {
	StudentName  string
	ProblemName  string
	Status       DashBoardStudentInfo
	LastSnapshot *Snapshot
	Messages     []*MessageDashBoard
	StudentID    int
	ProblemID    int
	UserID       int
	UserRole     string
	Password     string
	Username     string
}

type SubmissionInfo struct {
	ID          int
	Code        string
	Grade       string
	SubmittedAt time.Time
	SnapshotID  int
}

type SubmissionDashboard struct {
	Submissions []*SubmissionInfo
	StudentName string
	ProblemName string
	StudentID   int
	ProblemID   int
	UserID      int
	UserRole    string
	Password    string
	Username    string
}

type TemplateDate struct {
	Feedback       FeedbackProvisionDashBoard
	Submission     SubmissionDashboard
	Status         DashBoardStudentInfo
	ChatgptaServer string
	UserID         int
	UserRole       string
	Password       string
	Username       string
	CourseName     string
}

func getCurrentUserVote(feedbackID int, userID int, userRole string) string {
	var feedback MessageBackFeedback

	// Query the database using GORM
	err := DB.Where("message_feedback_id = ? AND author_id = ? AND author_role = ?", feedbackID, userID, userRole).
		Select("useful").First(&feedback).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "" // Return an empty string if no record is found
		}
		log.Printf("Error fetching vote: %v", err)
		return ""
	}

	return feedback.Useful
}

func getMessageFeedbacks(messageID int, userID int, userRole string) []*FeedbackDashBaord {
	var feedbacks []MessageFeedback

	// Query the database using GORM
	err := DB.Where("message_id = ?", messageID).Find(&feedbacks).Error
	if err != nil {
		log.Printf("Error fetching message feedbacks: %v", err)
		return nil
	}

	feedbackDashboards := make([]*FeedbackDashBaord, 0)

	for _, feedback := range feedbacks {
		// Fetch author name based on role
		var name string
		if feedback.AuthorRole == "teacher" {
			name = getTeacherName(feedback.AuthorID)
		} else {
			name = getStudentName(feedback.AuthorID)
		}

		// Populate the dashboard entry
		feedbackDashboards = append(feedbackDashboards, &FeedbackDashBaord{
			Name:            name,
			Role:            feedback.AuthorRole,
			Feedback:        feedback.Feedback,
			FeedbackID:      feedback.ID,
			CurrentUserVote: getCurrentUserVote(feedback.ID, userID, userRole),
			Downvote:        getBackFeedbackCount(feedback.ID, "no"),
			Upvote:          getBackFeedbackCount(feedback.ID, "yes"),
			GivenAt:         feedback.GivenAt,
		})
	}

	return feedbackDashboards
}

func getLatestSnapshot(studentID int, problemID int) *Snapshot {
	var snapshot Snapshot

	// Query the latest snapshot using GORM
	err := DB.Model(&Snapshot{}).
		Select("id, code, MAX(last_updated) as last_updated").
		Where("problem_id = ? AND student_id = ?", problemID, studentID).
		Group("problem_id, student_id").
		Scan(&snapshot).Error

	if err != nil {
		log.Printf("Error fetching latest snapshot: %v", err)
		return nil
	}

	// Add the problem name
	snapshot.ProblemName = getProblemNameFromID(problemID)

	return &snapshot
}

func getTeacherName(authorID int) string {
	var name string

	// Use GORM to fetch the teacher's name
	err := DB.Model(&Teacher{}).
		Select("name").
		Where("id = ?", authorID).
		Scan(&name).Error

	if err != nil {
		log.Printf("Error fetching teacher name for ID %d: %v", authorID, err)
		return ""
	}

	return name
}

func getStudentName(studentID int) string {
	var name string

	// Use GORM to fetch the student's name
	err := DB.Model(&Student{}).
		Select("name").
		Where("id = ?", studentID).
		Scan(&name).Error

	if err != nil {
		log.Printf("Error fetching student name for ID %d: %v", studentID, err)
		return ""
	}

	return name
}

func getBackFeedbackCount(feedbackID int, backFeedbackType string) int {
	var voteCount int64

	err := DB.Model(&MessageBackFeedback{}).
		Where("useful = ? AND message_feedback_id = ?", backFeedbackType, feedbackID).
		Count(&voteCount).Error

	if err != nil {
		log.Printf("Error counting back feedbacks for feedback ID %d: %v", feedbackID, err)
		return 0
	}

	return int(voteCount)
}

func studentDashboardFeedbackProvisionHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	studentID, _ := strconv.Atoi(r.FormValue("student_id"))
	temp := template.New("")
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode}
	t, err := temp.Funcs(ownFuncs).Parse(FEEDBACK_PROVISION_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}

	students := getAllStudents()
	var messages []*MessageDashBoard
	_, ok := HelpEligibleStudents[problemID][uid]

	if role == "teacher" || uid == studentID || (PeerTutorAllowed && ok) {
		var rawMessages []struct {
			ID         int
			SnapshotID int
			Message    string
			AuthorID   int
			AuthorRole string
			GivenAt    time.Time
			Type       int
			Code       string
			Event      string
		}

		err = DB.Table("messages").
			Select("M.id, M.snapshot_id, M.message, M.author_id, M.author_role, M.given_at, M.type, C.code, C.event").
			Joins("JOIN code_snapshots C ON M.snapshot_id = C.id").
			Where("C.problem_id = ? AND C.student_id = ?", problemID, studentID).
			Scan(&rawMessages).Error

		if err != nil {
			log.Printf("Error fetching messages: %v", err)
			http.Error(w, "Error fetching messages", http.StatusInternalServerError)
			return
		}

		for _, raw := range rawMessages {
			name := ""
			if raw.AuthorRole == "teacher" {
				name = getTeacherName(raw.AuthorID)
			} else {
				name = students[raw.AuthorID]
			}
			messages = append(messages, &MessageDashBoard{
				ID:         raw.ID,
				Name:       name,
				Role:       raw.AuthorRole,
				Message:    raw.Message,
				Type:       raw.Type,
				Event:      raw.Event,
				GivenAt:    raw.GivenAt,
				SnapshotID: raw.SnapshotID,
				Code:       raw.Code,
				Feedbacks:  getMessageFeedbacks(raw.ID, uid, role),
			})
		}
	} else {
		http.Error(w, "You are not authorized to access!", http.StatusUnauthorized)
		return
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].GivenAt.After(messages[j].GivenAt)
	})

	var latestSnapshot *Snapshot
	if _, ok := StudentSnapshot[studentID][problemID]; ok {
		latestSnapshot = Snapshots[StudentSnapshot[studentID][problemID]]
	} else {
		latestSnapshot = getLatestSnapshot(studentID, problemID)
	}

	var studentStatus DashBoardStudentInfo
	err = DB.Model(&StudentStatus{}).
		Select("coding_stat, help_stat, submission_stat, tutoring_stat").
		Where("problem_id = ? AND student_id = ?", problemID, studentID).
		Scan(&studentStatus).Error

	if err != nil {
		log.Printf("Error fetching student status: %v", err)
	}

	data := &FeedbackProvisionDashBoard{
		StudentName:  students[studentID],
		ProblemName:  latestSnapshot.ProblemName,
		LastSnapshot: latestSnapshot,
		Messages:     messages,
		StudentID:    studentID,
		ProblemID:    problemID,
		UserID:       uid,
		UserRole:     role,
		Password:     r.FormValue("password"),
		Status:       studentStatus,
		Username:     getName(uid, role),
	}

	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func studentDashboardSubmissionHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	studentID, _ := strconv.Atoi(r.FormValue("student_id"))
	temp := template.New("")
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode}
	t, err := temp.Funcs(ownFuncs).Parse(SUBMISSION_VIEW_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}

	var submissions []*SubmissionInfo
	_, ok := HelpEligibleStudents[problemID][uid]
	if role == "teacher" || uid == studentID || (PeerTutorAllowed && ok) {
		err := DB.Where("student_id = ? AND problem_id = ?", studentID, problemID).
			Find(&submissions).Error
		if err != nil {
			log.Fatal(err)
		}

		sort.SliceStable(submissions, func(i, j int) bool {
			if submissions[i].Grade == "" && submissions[j].Grade == "" {
				return submissions[i].SubmittedAt.Before(submissions[j].SubmittedAt)
			}
			if submissions[i].Grade == "" {
				return true
			}
			if submissions[j].Grade == "" {
				return false
			}
			return submissions[i].SubmittedAt.Before(submissions[j].SubmittedAt)
		})
	} else {
		http.Error(w, "You are not authorized to access!", http.StatusUnauthorized)
		return
	}

	data := &SubmissionDashboard{
		StudentName: getStudentName(studentID),
		ProblemName: getProblemNameFromID(problemID),
		Submissions: submissions,
		StudentID:   studentID,
		ProblemID:   problemID,
		UserID:      uid,
		UserRole:    role,
		Password:    r.FormValue("password"),
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func hasMessageBackFeedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	feedbackID, _ := strconv.Atoi(r.FormValue("feedback_id"))
	userRole := r.FormValue("role")
	fmt.Print(feedbackID, userRole, uid)

	var feedback MessageBackFeedback
	err := DB.Where("message_feedback_id = ? AND author_id = ? AND author_role = ?", feedbackID, uid, userRole).
		First(&feedback).Error
	if err != nil {
		if gorm.ErrRecordNotFound == err {
			fmt.Fprint(w, "no")
		} else {
			log.Fatal(err)
		}
		return
	}
	fmt.Fprint(w, "yes")
}

func studentDashboardCodeSpaceHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	studentID, _ := strconv.Atoi(r.FormValue("student_id"))
	temp := template.New("")
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode}
	t, err := temp.Funcs(ownFuncs).Parse(CODE_SNAPSHOT_TAB_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	students := getAllStudents()

	// Get latest snapshot from DB
	latestSnapshot := &Snapshot{}
	if _, ok := StudentSnapshot[studentID][problemID]; ok {
		latestSnapshot = Snapshots[StudentSnapshot[studentID][problemID]]
	} else {
		latestSnapshot = getLatestSnapshot(studentID, problemID)
	}

	// Get all student messages from DB
	var messages []MessageDashBoard
	_, ok := HelpEligibleStudents[problemID][uid]
	if role == "teacher" || uid == studentID || (PeerTutorAllowed && ok) {
		var msgList []MessageDashBoard
		err := DB.Joins("JOIN code_snapshots C ON M.snapshot_id = C.id").
			Where("C.problem_id = ? AND C.student_id = ?", problemID, studentID).
			Find(&msgList).Error
		if err != nil {
			log.Fatal(err)
		}

		// Fetching related message feedbacks for each message
		for i := range msgList {
			msgList[i].Feedbacks = getMessageFeedbacks(msgList[i].ID, uid, role)
		}

		// Sorting by 'GivenAt' in descending order
		sort.Slice(msgList, func(i, j int) bool {
			return msgList[i].GivenAt.After(msgList[j].GivenAt)
		})

		messages = msgList
	} else {
		http.Error(w, "You are not authorized to access!", http.StatusUnauthorized)
	}

	// Get all submissions from DB
	var submissions []SubmissionInfo
	if role == "teacher" || uid == studentID || (PeerTutorAllowed && ok) {
		err := DB.Where("student_id = ? AND problem_id = ?", studentID, problemID).
			Find(&submissions).Error
		if err != nil {
			log.Fatal(err)
		}

		// Sorting submissions by 'SubmittedAt' and 'Grade'
		sort.SliceStable(submissions, func(i, j int) bool {
			if submissions[i].Grade == "" && submissions[j].Grade == "" {
				return submissions[i].SubmittedAt.After(submissions[j].SubmittedAt)
			}
			if submissions[i].Grade == "" {
				return true
			}
			if submissions[j].Grade == "" {
				return false
			}
			return submissions[i].SubmittedAt.After(submissions[j].SubmittedAt)
		})
	} else {
		http.Error(w, "You are not authorized to access!", http.StatusUnauthorized)
	}

	// Get student status
	var studentStats DashBoardStudentInfo
	err = DB.Table("student_statuses").
		Where("problem_id = ? AND student_id = ?", problemID, studentID).
		First(&studentStats).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Fatal(err)
	}

	feedback := &FeedbackProvisionDashBoard{
		StudentName:  students[studentID],
		ProblemName:  latestSnapshot.ProblemName,
		LastSnapshot: latestSnapshot,
		Messages:     convertMessagesToPointers(messages), // Convert to []*MessageDashBoard
		StudentID:    studentID,
		ProblemID:    problemID,
		UserID:       uid,
		UserRole:     role,
		Password:     r.FormValue("password"),
	}

	submission := &SubmissionDashboard{
		StudentName: getStudentName(studentID),
		ProblemName: getProblemNameFromID(problemID),
		Submissions: convertSubmissionsToPointers(submissions), // Convert to []*SubmissionInfo
		StudentID:   studentID,
		ProblemID:   problemID,
		UserID:      uid,
		UserRole:    role,
		Password:    r.FormValue("password"),
		Username:    getName(uid, role),
	}

	data := TemplateDate{
		Submission:     *submission,
		Feedback:       *feedback,
		Status:         studentStats,
		ChatgptaServer: Config.ChatgptaServer,
		UserID:         uid,
		UserRole:       role,
		Password:       r.FormValue("password"),
		Username:       getName(uid, role),
		CourseName:     Config.CourseName,
	}

	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func convertMessagesToPointers(messages []MessageDashBoard) []*MessageDashBoard {
	var result []*MessageDashBoard
	for i := range messages {
		result = append(result, &messages[i])
	}
	return result
}

// Converts []SubmissionInfo to []*SubmissionInfo
func convertSubmissionsToPointers(submissions []SubmissionInfo) []*SubmissionInfo {
	var result []*SubmissionInfo
	for i := range submissions {
		result = append(result, &submissions[i])
	}
	return result
}

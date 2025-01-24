package main

import (
	"fmt"
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

func GetMessageFeedbacks(messageID int, userID int, userRole string) []*FeedbackDashBaord {
	messageFeedbacks, err := GetMessageFeedbacksByMessageID(messageID)
	if err != nil {
		return nil
	}

	var feedbacks []*FeedbackDashBaord
	for _, feedback := range messageFeedbacks {
		name := ""
		if feedback.AuthorRole == "teacher" {
			name = GetTeacherName(feedback.AuthorID)
		} else {
			name = GetStudentName(feedback.AuthorID)
		}

		feedbacks = append(feedbacks, &FeedbackDashBaord{
			Name:            name,
			Role:            feedback.AuthorRole,
			Feedback:        feedback.Feedback,
			FeedbackID:      feedback.ID,
			CurrentUserVote: GetCurrentUserVote(feedback.ID, userID, userRole),
			Downvote:        GetBackFeedbackCount(feedback.ID, "no"),
			Upvote:          GetBackFeedbackCount(feedback.ID, "yes"),
			GivenAt:         feedback.GivenAt,
		})
	}

	return feedbacks
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
	students := GetAllStudents()
	var messages = make([]*MessageDashBoard, 0)
	_, ok := HelpEligibleStudents[problemID][uid]
	if role == "teacher" || uid == studentID || (PeerTutorAllowed && ok) {
		// Fetch messages using the new function
		messages, err = FetchMessagesForStudent(students, problemID, studentID, role)
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
	latestSnapshot := &Snapshot{}
	if _, ok := StudentSnapshot[studentID][problemID]; ok {
		latestSnapshot = Snapshots[StudentSnapshot[studentID][problemID]]
	} else {
		latestSnapshot, _ = GetLatestSnapshot(studentID, problemID)
	}

	// Get student status
	studentStats := &DashBoardStudentInfo{}
	if role == "teacher" || uid == studentID || (PeerTutorAllowed && ok) {
		// Fetch student status using the new function
		studentStats, err = FetchStudentStatuses(problemID, studentID)
		if err != nil {
			log.Fatal(err) // or use http.Error depending on your context
		}

		// Do something with studentStats, for example, returning them as JSON
		// Example: json.NewEncoder(w).Encode(studentStats)
	} else {
		http.Error(w, "You are not authorized to access!", http.StatusUnauthorized)
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
		Status:       *studentStats,
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

	var submissions = make([]*SubmissionInfo, 0)
	_, ok := HelpEligibleStudents[problemID][uid]
	if role == "teacher" || uid == studentID || (PeerTutorAllowed && ok) {
		// Fetch submissions using the new function
		submissionRecords, err := FetchSubmission(studentID, problemID)
		if err != nil {
			log.Fatal(err)
		}

		// Convert SubmissionTable to SubmissionInfo and append to submissionInfos
		var submissionInfos []*SubmissionInfo
		for _, submission := range submissionRecords {
			submissionInfos = append(submissionInfos, &SubmissionInfo{
				ID:          submission.ID,
				SnapshotID:  submission.SnapshotID,
				Code:        submission.StudentCode,
				Grade:       submission.Verdict,
				SubmittedAt: submission.CodeSubmittedAt,
			})
		}

		// Sort submissions if needed (optional, as GORM already sorts them)
		sort.SliceStable(submissionInfos, func(i, j int) bool {
			if submissionInfos[i].Grade == "" && submissionInfos[j].Grade == "" {
				return submissionInfos[i].SubmittedAt.Before(submissionInfos[j].SubmittedAt)
			}
			if submissionInfos[i].Grade == "" {
				return true
			}
			if submissionInfos[j].Grade == "" {
				return false
			}
			return submissionInfos[i].SubmittedAt.Before(submissionInfos[j].SubmittedAt)
		})
	} else {
		http.Error(w, "You are not authorized to access!", http.StatusUnauthorized)
	}

	// TODO(shiplu): sort the messages array
	// sort.Slice(helpRequests, func(i, j int) bool { return helpRequests[i].GivenAt.Before(helpRequests[j].GivenAt) })
	data := &SubmissionDashboard{
		StudentName: GetStudentName(studentID),
		ProblemName: GetProblemNameFromID(problemID),
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

	// Check if feedback exists
	exists, err := CheckMessageBackFeedback(feedbackID, uid, userRole)
	if err != nil {
		log.Printf("Error checking feedback: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Respond to the client
	if exists {
		fmt.Fprint(w, "yes")
	} else {
		fmt.Fprint(w, "no")
	}
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
	students := GetAllStudents()

	// Get latest snapshot from DB
	latestSnapshot := &Snapshot{}
	if _, ok := StudentSnapshot[studentID][problemID]; ok {
		latestSnapshot = Snapshots[StudentSnapshot[studentID][problemID]]
	} else {
		latestSnapshot, _ = GetLatestSnapshot(studentID, problemID)
	}

	// Get all student messages from DB
	var messages = make([]*MessageDashBoard, 0)
	_, ok := HelpEligibleStudents[problemID][uid]
	if role == "teacher" || uid == studentID || (PeerTutorAllowed && ok) {
		messages, err = FetchMessages(problemID, studentID, uid, role, students)
		if err != nil {
			log.Fatal(err) // Consider graceful error handling
		}
		// Process `messages` as needed
	} else {
		http.Error(w, "You are not authorized to access!", http.StatusUnauthorized)
	}

	// Sort the messages by descding order of GivenAt
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].GivenAt.After(messages[j].GivenAt)
	})

	feedback := &FeedbackProvisionDashBoard{
		StudentName:  students[studentID],
		ProblemName:  latestSnapshot.ProblemName,
		LastSnapshot: latestSnapshot,
		Messages:     messages,
		StudentID:    studentID,
		ProblemID:    problemID,
		UserID:       uid,
		UserRole:     role,
		Password:     r.FormValue("password"),
	}

	// Get all submissions from DB.
	var submissions = make([]*SubmissionInfo, 0)
	_, ok = HelpEligibleStudents[problemID][uid]
	if role == "teacher" || uid == studentID || (PeerTutorAllowed && ok) {
		submissions, err = FetchSubmissions(problemID, studentID)
		if err != nil {
			log.Printf("Error fetching submissions: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "You are not authorized to access!", http.StatusUnauthorized)
	}

	submission := &SubmissionDashboard{
		StudentName: GetStudentName(studentID),
		ProblemName: GetProblemNameFromID(problemID),
		Submissions: submissions,
		StudentID:   studentID,
		ProblemID:   problemID,
		UserID:      uid,
		UserRole:    role,
		Password:    r.FormValue("password"),
		Username:    getName(uid, role),
	}

	// Get student status
	studentStats, err := FetchStudentStatus(problemID, studentID)
	if err != nil {
		log.Fatal(err) // Consider handling the error gracefully instead of terminating the program
	}

	// Apply role-based authorization
	if !(role == "teacher" || uid == studentID || (PeerTutorAllowed && ok)) {
		http.Error(w, "You are not authorized to access!", http.StatusUnauthorized)
		return
	}

	data := TemplateDate{
		Submission:     *submission,
		Feedback:       *feedback,
		Status:         *studentStats,
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

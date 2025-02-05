package restHandlers

import (
	"fmt"
	"github.com/GPTA/src/frontEnd"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
)

func GetMessageFeedbacks(messageID int, userID int, userRole string) []*models.FeedbackDashBoard {
	messageFeedbacks, err := repository.GetMessageFeedbacksByMessageID(messageID)
	if err != nil {
		return nil
	}

	var feedbacks []*models.FeedbackDashBoard
	for _, feedback := range messageFeedbacks {
		name := ""
		if feedback.AuthorRole == "teacher" {
			name = repository.GetTeacherName(feedback.AuthorID)
		} else {
			name = repository.GetStudentName(feedback.AuthorID)
		}

		feedbacks = append(feedbacks, &models.FeedbackDashBoard{
			Name:            name,
			Role:            feedback.AuthorRole,
			Feedback:        feedback.Feedback,
			FeedbackID:      feedback.ID,
			CurrentUserVote: repository.GetCurrentUserVote(feedback.ID, userID, userRole),
			Downvote:        repository.GetBackFeedbackCount(feedback.ID, "no"),
			Upvote:          repository.GetBackFeedbackCount(feedback.ID, "yes"),
			GivenAt:         feedback.GivenAt,
		})
	}

	return feedbacks
}

func StudentDashboardFeedbackProvisionHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	studentID, _ := strconv.Atoi(r.FormValue("student_id"))
	temp := template.New("")
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode}
	t, err := temp.Funcs(ownFuncs).Parse(frontEnd.FEEDBACK_PROVISION_TEMPLATE)
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

func StudentDashboardSubmissionHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	studentID, _ := strconv.Atoi(r.FormValue("student_id"))
	temp := template.New("")
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode}
	t, err := temp.Funcs(ownFuncs).Parse(frontEnd.SUBMISSION_VIEW_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}

	var submissions = make([]*models.SubmissionInfo, 0)
	_, ok := models.HelpEligibleStudents[problemID][uid]
	if role == "teacher" || uid == studentID || (models.PeerTutorAllowed && ok) {
		// Fetch submissions using the new function
		submissionRecords, err := repository.FetchSubmission(studentID, problemID)
		if err != nil {
			log.Fatal(err)
		}

		// Convert SubmissionTable to SubmissionInfo and append to submissionInfos
		var submissionInfos []*models.SubmissionInfo
		for _, submission := range submissionRecords {
			submissionInfos = append(submissionInfos, &models.SubmissionInfo{
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
	data := &models.SubmissionDashboard{
		StudentName: repository.GetStudentName(studentID),
		ProblemName: repository.GetProblemNameFromID(problemID),
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

func HasMessageBackFeedbackHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	feedbackID, _ := strconv.Atoi(r.FormValue("feedback_id"))
	userRole := r.FormValue("role")

	// Check if feedback exists
	exists, err := repository.CheckMessageBackFeedback(feedbackID, uid, userRole)
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

func StudentDashboardCodeSpaceHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	problemID, _ := strconv.Atoi(r.FormValue("problem_id"))
	studentID, _ := strconv.Atoi(r.FormValue("student_id"))
	temp := template.New("")
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode}
	t, err := temp.Funcs(ownFuncs).Parse(frontEnd.CODE_SNAPSHOT_TAB_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	students := repository.GetAllStudents()

	// Get latest snapshot from DB
	latestSnapshot := &models.Snapshot{}
	if _, ok := models.StudentSnapshot[studentID][problemID]; ok {
		latestSnapshot = models.Snapshots[models.StudentSnapshot[studentID][problemID]]
	} else {
		latestSnapshot, _ = repository.GetLatestSnapshot(studentID, problemID)
	}

	// Get all student messages from DB
	var messages = make([]*models.MessageDashBoard, 0)
	_, ok := models.HelpEligibleStudents[problemID][uid]
	if role == "teacher" || uid == studentID || (models.PeerTutorAllowed && ok) {
		messages, err = FetchMessages(students, problemID, studentID, role)
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

	feedback := &models.FeedbackProvisionDashBoard{
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
	var submissions = make([]*models.SubmissionInfo, 0)
	_, ok = models.HelpEligibleStudents[problemID][uid]
	if role == "teacher" || uid == studentID || (models.PeerTutorAllowed && ok) {
		submissions, err = repository.FetchSubmissions(problemID, studentID)
		if err != nil {
			log.Printf("Error fetching submissions: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "You are not authorized to access!", http.StatusUnauthorized)
	}

	submission := &models.SubmissionDashboard{
		StudentName: repository.GetStudentName(studentID),
		ProblemName: repository.GetProblemNameFromID(problemID),
		Submissions: submissions,
		StudentID:   studentID,
		ProblemID:   problemID,
		UserID:      uid,
		UserRole:    role,
		Password:    r.FormValue("password"),
		Username:    GetName(uid, role),
	}

	// Get student status
	studentStats, err := repository.FetchStudentStatus(problemID, studentID)
	if err != nil {
		log.Fatal(err) // Consider handling the error gracefully instead of terminating the program
	}

	// Apply role-based authorization
	if !(role == "teacher" || uid == studentID || (models.PeerTutorAllowed && ok)) {
		http.Error(w, "You are not authorized to access!", http.StatusUnauthorized)
		return
	}

	data := models.TemplateDate{
		Submission: *submission,
		Feedback:   *feedback,
		Status:     *studentStats,
		UserID:     uid,
		UserRole:   role,
		Password:   r.FormValue("password"),
		Username:   GetName(uid, role),
		CourseName: models.Config.CourseName,
	}

	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func FetchMessages(students map[int]string, problemID, studentID int, role string) ([]*models.MessageDashBoard, error) {
	result, err := repository.FetchMessagesAndSnapshots(problemID, studentID)
	if err != nil {
		return nil, err
	}

	var messages []*models.MessageDashBoard
	for _, r := range result {
		name := ""
		if r.AuthorRole == "teacher" {
			name = repository.GetTeacherName(r.AuthorID)
		} else {
			name = students[r.AuthorID]
		}

		messages = append(messages, &models.MessageDashBoard{
			ID:         r.MessageID,
			Name:       name,
			Role:       r.AuthorRole,
			Message:    r.Message,
			Type:       r.MessageType,
			Event:      r.Event,
			GivenAt:    r.GivenAt,
			SnapshotID: r.SnapshotID,
			Code:       r.Code,
			Feedbacks:  GetMessageFeedbacks(r.MessageID, studentID, role),
		})
	}

	return messages, nil
}

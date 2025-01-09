package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type CodeSpaceData struct {
	Snapshots []*Snapshot
	UserID    int
	UserRole  string
	Password  string
}

type FeedbackData struct {
	FeedbackID      int
	Feedback        string
	FeedbackTime    time.Time
	Upvote          int
	Downvote        int
	CurrentUserVote string
	GivenBy         string
	Code            string
}
type SnapshotData struct {
	Snapshot       *Snapshot
	HelpRequestIDs []int
	UserID         int
	UserRole       string
	Feedbacks      []*FeedbackData
	Password       string
}

type HelpRequest struct {
	ID          int
	StudentName string
	Explanation string
	NumReply    int
	GivenAt     time.Time
	SnapshotID  int
	Snapshot    string
	ProblemName string
	UserID      int
	UserRole    string
	Password    string
}

type HelpRequestListData struct {
	HelpRequests  []*HelpRequest
	NumHelpNeeded int
	UserID        int
	UserRole      string
	Password      string
}

func getEditorMode(filename string) string {
	filename = strings.ToLower(filename)
	if strings.HasSuffix(filename, ".py") {
		return "python"
	}
	if strings.HasSuffix(filename, ".java") {
		return "text/x-java"
	}
	if strings.HasSuffix(filename, ".cpp") || strings.HasSuffix(filename, ".c++") || strings.HasSuffix(filename, ".c") {
		return "text/x-c++src"
	}
	return "text"
}

func formatTimeDuration(d time.Duration) string {
	m := int(d.Minutes())
	d1 := d - time.Duration(m*60*1000000000)
	s := int(d1.Seconds())
	str := ""
	if m > 0 {
		str = strconv.Itoa(m) + " min "
	}
	str += strconv.Itoa(s) + " sec"
	return str
}

func formatTimeSince(t time.Time) string {
	d := time.Now().Sub(t)
	return formatTimeDuration(d)
}

// test
func getVoteCount(feedbackID int, voteType string) int {
	var vote int64

	// Use GORM to count the votes
	if err := DB.Model(&SnapshotBackFeedback{}).
		Where("is_helpful = ? AND snapshot_feedback_id = ?", voteType, feedbackID).
		Count(&vote).Error; err != nil {
		log.Fatal("Error counting votes: ", err)
	}

	return int(vote)
}

func add(x int, y int) int {
	return x + y
}

func codespaceHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	temp := template.New("")
	ownFuncs := template.FuncMap{"formatTimeSince": formatTimeSince}
	t, err := temp.Funcs(ownFuncs).Parse(CODESPACE_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	var snapshots []*Snapshot
	if role == "student" {
		for _, s := range Snapshots {
			if s.StudentID == uid {
				snapshots = append(snapshots, s)
			} else if _, ok := HelpEligibleStudents[s.ProblemID][uid]; ok {
				snapshots = append(snapshots, s)
			}
		}
	} else {
		for _, s := range Snapshots {
			snapshots = append(snapshots, s)
		}
	}
	sort.Slice(snapshots, func(i, j int) bool { return snapshots[i].LastUpdated.After(snapshots[j].LastUpdated) })
	data := &CodeSpaceData{
		Snapshots: snapshots,
		UserID:    uid,
		UserRole:  role,
		Password:  r.FormValue("password"),
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func getCodeSnapshotHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	var snapshot CodeSnapshot
	snapshotID := 0
	studentID := 0
	problemID := 0

	if r.FormValue("snapshot_id") != "" {
		// Parse snapshot_id and fetch details using the CodeSnapshot model
		snapshotID, _ = strconv.Atoi(r.FormValue("snapshot_id"))
		if err := DB.First(&snapshot, snapshotID).Error; err != nil {
			log.Fatal("Error fetching snapshot: ", err)
		}
		studentID = snapshot.StudentID
		problemID = snapshot.ProblemID
	} else {
		studentID, _ = strconv.Atoi(r.FormValue("student_id"))
		problemID, _ = strconv.Atoi(r.FormValue("problem_id"))
	}

	role := r.FormValue("role")
	temp := template.New("")
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode, "add": add}
	t, err := temp.Funcs(ownFuncs).Parse(CODE_SNAPSHOT_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch helpIDs using GORM
	var helpIDs []int
	if err := DB.Model(&CodeExplanation{}).
		Where("student_id = ? AND problem_id = ?", studentID, problemID).
		Pluck("id", &helpIDs).Error; err != nil {
		log.Fatal("Error fetching help request IDs: ", err)
	}

	// Fetch feedback and related details using GORM
	var feedbacks []*FeedbackData
	var snapshotFeedbacks []struct {
		ID         int
		Feedback   string
		AuthorID   int
		AuthorRole string
		GivenAt    time.Time
		Code       string
	}
	if err := DB.Table("code_snapshot AS C").
		Select("F.id, F.feedback, F.author_id, F.author_role, F.given_at, C.code").
		Joins("JOIN snapshot_feedback F ON C.id = F.snapshot_id").
		Where("C.student_id = ? AND C.problem_id = ?", studentID, problemID).
		Order("F.given_at DESC").
		Scan(&snapshotFeedbacks).Error; err != nil {
		log.Fatal("Error fetching feedback: ", err)
	}

	for _, sf := range snapshotFeedbacks {
		upvote := getVoteCount(sf.ID, "yes")
		downvote := getVoteCount(sf.ID, "no")
		var currentUserVote string

		// Fetch current user vote
		if err := DB.Model(&SnapshotBackFeedback{}).
			Where("snapshot_feedback_id = ? AND author_id = ? AND author_role = ?", sf.ID, uid, role).
			Pluck("is_helpful", &currentUserVote).Error; err != nil && err != gorm.ErrRecordNotFound {
			log.Fatal("Error fetching current user vote: ", err)
		}

		// Fetch author name
		var authorName string
		if sf.AuthorRole == "teacher" {
			DB.Model(&Teacher{}).Where("id = ?", sf.AuthorID).Pluck("name", &authorName)
		} else {
			DB.Model(&Student{}).Where("id = ?", sf.AuthorID).Pluck("name", &authorName)
		}

		feedbacks = append(feedbacks, &FeedbackData{
			FeedbackID:      sf.ID,
			Feedback:        sf.Feedback,
			FeedbackTime:    sf.GivenAt,
			Upvote:          upvote,
			Downvote:        downvote,
			CurrentUserVote: currentUserVote,
			GivenBy:         authorName,
			Code:            sf.Code,
		})
	}

	idx := StudentSnapshot[studentID][problemID]
	data := &SnapshotData{
		Snapshot:       Snapshots[idx],
		HelpRequestIDs: helpIDs,
		UserID:         uid,
		UserRole:       role,
		Feedbacks:      feedbacks,
		Password:       r.FormValue("password"),
	}
	w.Header().Set("Content-Type", "text/html")
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	// Log event
	otherInfo := struct {
		Referral  string
		ProblemID int
	}{
		Referral:  r.URL.String(),
		ProblemID: problemID,
	}
	elog, _ := json.Marshal(otherInfo)
	logEvent("willing to help", uid, role, "click", string(elog))
}

func getNumberOfReply(snapshotID int) int {
	var count int64
	err := DB.Model(&SnapshotFeedback{}).Where("snapshot_id = ?", snapshotID).Count(&count).Error
	if err != nil {
		log.Fatal("Error fetching number of replies: ", err)
	}
	return int(count)
}

func helpRequestListHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	temp := template.New("")
	ownFuncs := template.FuncMap{"formatTimeSince": formatTimeSince}
	t, err := temp.Funcs(ownFuncs).Parse(HELP_REQUEST_LIST_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	var helpRequests []*HelpRequest
	var helpNeededCount int
	var numReply int
	if role == "student" {
		for _, s := range HelpSubmissions {
			if _, ok := HelpEligibleStudents[s.Pid][uid]; ok || s.Uid == uid {
				numReply = getNumberOfReply(s.SnapshotID)
				helpRequests = append(helpRequests, &HelpRequest{
					ID:          s.Sid,
					NumReply:    numReply,
					StudentName: Students[s.Uid].Name,
					GivenAt:     s.At,
				})
				if numReply == 0 {
					helpNeededCount++
				}
			}
		}
	} else {
		for _, s := range HelpSubmissions {
			numReply = getNumberOfReply(s.SnapshotID)
			helpRequests = append(helpRequests, &HelpRequest{
				ID:          s.Sid,
				NumReply:    numReply,
				StudentName: Students[s.Uid].Name,
				GivenAt:     s.At,
			})
			if numReply == 0 {
				helpNeededCount++
			}
		}
	}
	sort.Slice(helpRequests, func(i, j int) bool { return helpRequests[i].GivenAt.Before(helpRequests[j].GivenAt) })
	data := &HelpRequestListData{
		HelpRequests:  helpRequests,
		NumHelpNeeded: helpNeededCount,
		UserID:        uid,
		UserRole:      role,
		Password:      r.FormValue("password"),
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func viewHelpRequestHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	requestID, _ := strconv.Atoi(r.FormValue("request_id"))
	role := r.FormValue("role")
	pw := r.FormValue("password")
	temp := template.New("")
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode, "formatTimeSince": formatTimeSince}
	t, err := temp.Funcs(ownFuncs).Parse(HELP_REQUEST_VIEW_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	data := &HelpRequest{}
	problemID := -1
	for _, s := range HelpSubmissions {
		if s.Sid == requestID {
			data = &HelpRequest{
				StudentName: Students[s.Uid].Name,
				Explanation: s.Content,
				GivenAt:     s.At,
				SnapshotID:  s.SnapshotID,
				Snapshot:    s.Snapshot,
				ProblemName: s.Filename,
				UserID:      uid,
				UserRole:    role,
				Password:    pw,
			}
			problemID = s.Pid
		}
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
	otherInfo := struct {
		Referral  string
		ProblemID int
	}{
		Referral:  r.URL.String(),
		ProblemID: problemID,
	}
	elog, _ := json.Marshal(otherInfo)
	logEvent("willing to help", uid, role, "click", string(elog))

}

func setPeerTutorHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	turnOn, _ := strconv.Atoi(r.FormValue("turn_on"))
	if turnOn == 1 {
		PeerTutorAllowed = true
	} else if turnOn == 0 {
		PeerTutorAllowed = false
	}
}

func peerTutorHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	filename := r.FormValue("filename")
	if prob, ok := ActiveProblems[filename]; ok {
		if eligible, ok := HelpEligibleStudents[prob.Info.Pid][uid]; ok {
			if PeerTutorAllowed && eligible {
				fmt.Fprint(w, "redirect")
				return
			}
		}
	}
	fmt.Fprintf(w, "You are not eligible for peer tutoring for this problem!")
}

package restHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/GPTA/src/frontEnd"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type CodeSpaceData struct {
	Snapshots []*models.Snapshot
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
	Snapshot       *models.Snapshot
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
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	var parts []string

	if days > 0 {
		parts = append(parts, strconv.Itoa(days)+" day"+plural(days))
	}
	if hours > 0 {
		parts = append(parts, strconv.Itoa(hours)+" hr"+plural(hours))
	}
	if minutes > 0 {
		parts = append(parts, strconv.Itoa(minutes)+" min"+plural(minutes))
	}
	if seconds > 0 && len(parts) == 0 { // Only show seconds if no larger units exist
		parts = append(parts, strconv.Itoa(seconds)+" sec"+plural(seconds))
	}

	return strings.Join(parts, " ")
}
func plural(value int) string {
	if value != 1 {
		return "s"
	}
	return ""
}

func FormatTimeSince(t time.Time) string {
	d := time.Since(t)
	return formatTimeDuration(d)
}

func add(x int, y int) int {
	return x + y
}

func codespaceHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	temp := template.New("")
	ownFuncs := template.FuncMap{"FormatTimeSince": FormatTimeSince}
	t, err := temp.Funcs(ownFuncs).Parse(frontEnd.CODESPACE_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	var snapshots []*models.Snapshot
	if role == "student" {
		for _, s := range models.Snapshots {
			if s.StudentID == uid {
				snapshots = append(snapshots, s)
			} else if _, ok := models.HelpEligibleStudents[s.ProblemID][uid]; ok {
				snapshots = append(snapshots, s)
			}
		}
	} else {
		for _, s := range models.Snapshots {
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

func helpRequestListHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	role := r.FormValue("role")
	temp := template.New("")
	ownFuncs := template.FuncMap{"FormatTimeSince": FormatTimeSince}
	t, err := temp.Funcs(ownFuncs).Parse(frontEnd.HELP_REQUEST_LIST_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	var helpRequests []*HelpRequest
	var helpNeededCount int
	var numReply int
	if role == "student" {
		for _, s := range models.HelpSubmissions {
			if _, ok := models.HelpEligibleStudents[s.Pid][uid]; ok || s.Uid == uid {
				numReply, _ = repository.GetNumberOfReply(s.SnapshotID)
				helpRequests = append(helpRequests, &HelpRequest{
					ID:          s.Sid,
					NumReply:    numReply,
					StudentName: models.Students[s.Uid].Name,
					GivenAt:     s.At,
				})
				if numReply == 0 {
					helpNeededCount++
				}
			}
		}
	} else {
		for _, s := range models.HelpSubmissions {
			numReply, _ = repository.GetNumberOfReply(s.SnapshotID)
			helpRequests = append(helpRequests, &HelpRequest{
				ID:          s.Sid,
				NumReply:    numReply,
				StudentName: models.Students[s.Uid].Name,
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
	ownFuncs := template.FuncMap{"getEditorMode": getEditorMode, "FormatTimeSince": FormatTimeSince}
	t, err := temp.Funcs(ownFuncs).Parse(frontEnd.HELP_REQUEST_VIEW_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	data := &HelpRequest{}
	problemID := -1
	for _, s := range models.HelpSubmissions {
		if s.Sid == requestID {
			data = &HelpRequest{
				StudentName: models.Students[s.Uid].Name,
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
	LogEvent("willing to help", uid, role, "click", string(elog))

}

func SetPeerTutorHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	turnOn, _ := strconv.Atoi(r.FormValue("turn_on"))
	if turnOn == 1 {
		models.PeerTutorAllowed = true
	} else if turnOn == 0 {
		models.PeerTutorAllowed = false
	}
}

func PeerTutorHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	filename := r.FormValue("filename")
	if prob, ok := models.ActiveProblems[filename]; ok {
		if eligible, ok := models.HelpEligibleStudents[prob.Info.Pid][uid]; ok {
			if models.PeerTutorAllowed && eligible {
				fmt.Fprint(w, "redirect")
				return
			}
		}
	}
	fmt.Fprintf(w, "You are not eligible for peer tutoring for this problem!")
}

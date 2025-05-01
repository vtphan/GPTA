package restHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/GPTA/src/frontEnd"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/openAI"
	"github.com/GPTA/src/repository"
	"golang.org/x/net/html"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/teacher_signin", http.StatusSeeOther)
	} else {
		sessionToken := c.Value
		userSession, exists := models.Sessions[sessionToken]
		if !exists || userSession.IsExpired() {
			http.Redirect(w, r, "/teacher_signin", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/view_exercises?role=teacher&uid="+strconv.Itoa(models.TeacherNameToId[userSession.Username]), http.StatusFound)
		}
	}
}

func AdminSigninHandler(w http.ResponseWriter, r *http.Request) {
	// Get the form values (admin username and password)
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Check if the username exists in the admin map and match the password
	expectedPass, ok := models.AdminPass[username]
	if !ok || expectedPass != password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Admin login success, create a session
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(3 * time.Hour)

	models.Sessions[sessionToken] = models.Session{
		Username: username,
		Expiry:   expiresAt,
	}

	// Set the session token as a cookie in the response
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	// Redirect to the admin dashboard page after login
	http.Redirect(w, r, "/admin_dashboard", http.StatusSeeOther)
}

func AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the session token is valid
	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie == nil {
		http.Redirect(w, r, "/admin_signin", http.StatusSeeOther)
		return
	}

	// Validate session token (you may want to check if the session token exists and is not expired)
	session, ok := models.Sessions[sessionCookie.Value]
	if !ok || session.Expiry.Before(time.Now()) {
		// Session is invalid or expired, redirect to signin
		http.Redirect(w, r, "/admin_signin", http.StatusSeeOther)
		return
	}

	// If session is valid, show the admin dashboard
	fmt.Fprintf(w, frontEnd.ADMIN_DASHBOARD)
}
func TeacherSigninCompleteHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("username")
	password := r.FormValue("password")
	expectedPass, ok := models.TeacherPass[name]
	if !ok || expectedPass != password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create session token
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(3 * time.Hour)

	models.Sessions[sessionToken] = models.Session{
		Username: name,
		Expiry:   expiresAt,
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
	models.LoggedInTeacher = name
	fmt.Fprintf(w, "%d", models.TeacherNameToId[name])
}

func TeacherSigninHandler(w http.ResponseWriter, r *http.Request) {
	temp := template.New("")
	t, err := temp.Parse(frontEnd.TEACHER_LOGIN)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "Text/html")
	err = t.Execute(w, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

type SettingsData struct {
	Username string
}

func SettingsViewHandler(w http.ResponseWriter, r *http.Request) {
	username := models.LoggedInTeacher
	data := SettingsData{Username: username}
	temp := template.New("")
	t, err := temp.Parse(frontEnd.SETTINGS_VIEW)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "Text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

type VisualData struct {
	Text template.JS
}

func ViViewHandler(w http.ResponseWriter, r *http.Request) {
	text := openAI.ProcessStudentCodeSubmissions(w, r)
	data := VisualData{Text: template.JS(text)}
	temp := template.New("")
	t, err := temp.Parse(frontEnd.VIS_TEMPLATE)
	if err != nil {
		log.Printf(err.Error())
	}
	w.Header().Set("Content-Type", "Text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf(err.Error())
	}
}

func ScViewHandler(w http.ResponseWriter, r *http.Request) {
	text := openAI.ProcessIndividualSubmissions(w, r)
	data := VisualData{Text: template.JS(text)}
	temp := template.New("")
	t, err := temp.Parse(frontEnd.SC_TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "Text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func AssignScaffoldHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		StudentID   int    `json:"student_id"`
		ProblemID   int    `json:"problem_id"`
		Duration    int    `json:"duration"`
		Scaffolding string `json:"scaffolding"`
		SnapshotID  int    `json:"snapshot_id"`
		AuthorID    int    `json:"uid"`
		AuthorRole  string `json:"role"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	scaffold := models.AssignedScaffolding{
		StudentID:   data.StudentID,
		ProblemID:   data.ProblemID,
		Duration:    data.Duration,
		Scaffolding: ExtractTextFromHTML(data.Scaffolding),
	}

	result := models.DB.Create(&scaffold)
	if result.Error != nil {
		http.Error(w, "Failed to save scaffold", http.StatusInternalServerError)
		return
	}

	snapshotID := data.SnapshotID
	authorID := data.AuthorID
	authorRole := data.AuthorRole
	now := time.Now()
	feedback := scaffold.Scaffolding

	mid, err := repository.AddMessage(snapshotID, "", authorID, authorRole, now, 1)
	if err != nil {
		log.Fatal("Could not save feedback for error: ", err)
		return
	}

	// Use GetCodeSnapshot instead of raw SQL query
	codeSnapshot, err := repository.GetCodeSnapshot(snapshotID)
	if err != nil {
		log.Fatal("Failed to retrieve code snapshot: ", err)
		return
	}

	// Check if the author is trying to give feedback on their own code
	if authorRole == "student" && codeSnapshot.StudentID == authorID {
		fmt.Fprintf(w, "You can not give feedback to your own code.")
		return
	}

	// Update feedback count
	idx := models.StudentSnapshot[codeSnapshot.StudentID][codeSnapshot.ProblemID]
	models.Snapshots[idx].NumFeedback++

	// Add feedback message
	messageID := mid
	id, err := repository.AddMessageFeedback(messageID, feedback, authorID, authorRole, now)
	if err != nil {
		log.Fatal("Could not save feedback for error: ", err)
		return
	}
	feedbackID := id

	// Append feedback to the student's queue
	models.Students[codeSnapshot.StudentID].SnapShotFeedbackQueue = append(models.Students[codeSnapshot.StudentID].SnapShotFeedbackQueue, &models.SnapShotFeedback{
		FeedbackID:  int(feedbackID),
		Snapshot:    codeSnapshot.Code,
		Feedback:    feedback,
		ProblemName: codeSnapshot.Event, // Assuming Event holds the problem name
		Provider:    GetName(authorID, authorRole),
	})

	// Update student status
	repository.AddOrUpdateStudentStatus(codeSnapshot.StudentID, codeSnapshot.ProblemID, "", "Been helped", "")
	if authorRole == "student" {
		repository.AddOrUpdateStudentStatus(authorID, codeSnapshot.ProblemID, "", "", "")
	}
	fmt.Println("Feedback on code snapshot saved!")
	w.WriteHeader(http.StatusCreated)
}

func ExtractTextFromHTML(htmlStr string) string {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return htmlStr // fallback to original if parse fails
	}
	var f func(*html.Node) string
	f = func(n *html.Node) string {
		if n.Type == html.TextNode {
			return n.Data
		}
		if n.Type == html.ElementNode && (n.Data == "br" || n.Data == "p" || n.Data == "div") {
			return "\n" + concatText(n.FirstChild, f)
		}
		return concatText(n.FirstChild, f)
	}
	return strings.TrimSpace(f(doc))
}

func concatText(n *html.Node, f func(*html.Node) string) string {
	var result string
	for c := n; c != nil; c = c.NextSibling {
		result += f(c)
	}
	return result
}

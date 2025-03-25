package restHandlers

import (
	"fmt"
	"github.com/GPTA/src/frontEnd"
	"github.com/GPTA/src/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
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
	w.Header().Set("Content-Type", "text/html")
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
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (app App) indexHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	} else {
		sessionToken := c.Value
		userSession, exists := app.Sessions[sessionToken]
		if !exists || userSession.isExpired() {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/view_exercises?uid="+strconv.Itoa(userSession.user.ID), http.StatusFound)
		}
	}
}

func (app App) signinCompleteHandler(w http.ResponseWriter, r *http.Request) {
	user := &User{
		Name:     r.FormValue("username"),
		Password: app.getMD5Hash(r.FormValue("password")),
	}
	var count int64

	app.DB.Model(&user).Count(&count)
	if count == 1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(3 * time.Hour)

	app.Sessions[sessionToken] = session{
		user:   user,
		expiry: expiresAt,
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
	// fmt.Fprintf(w, "%d", TeacherNameToId[name])
	// http.Redirect(w, r, "view_exercises?role=teacher", http.StatusFound)
}

func (app App) signinHandler(w http.ResponseWriter, r *http.Request) {
	temp := template.New("")
	t, err := temp.Parse(TEACHER_LOGIN)
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

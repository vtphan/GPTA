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

func TeacherSigninCompleteHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("username")
	password := r.FormValue("password")
	courseID := r.FormValue("course_id") // Get course_id from request
	expectedPass, ok := models.TeacherPass[name]
	if !ok || expectedPass != password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	teacherId, ok := models.TeacherNameToId[name]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if teacherCourses, exists := models.TeacherClassesMap[teacherId]; !exists {
		// Teacher is not mapped to any courses
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	} else {
		// Check if the course_id exists in the teacherCourses slice
		isAuthorized := false
		for _, course := range teacherCourses {
			if course == courseID {
				isAuthorized = true
				break
			}
		}

		if !isAuthorized {
			// Teacher is not authorized for this course
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
	// Store the selected course ID
	models.CourseId = courseID

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

// Author: Vinhthuy Phan, 2018
package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
)

// -----------------------------------------------------------------
// Authorize localhost
// -----------------------------------------------------------------
func AuthorizeLocalhost(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Host == "localhost:8080" {
			fn(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("Unauthorized access: host is not local.")
			fmt.Fprint(w, "Unauthorized access: host is not local.")
		}
	}
}

// -----------------------------------------------------------------
func Authorize(fn func(http.ResponseWriter, *http.Request, string, int), userRole string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		givenRole := r.FormValue("role")
		if userRole != "" && userRole != givenRole {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("Unauthorized access:", r.FormValue("name"), " user role: ", givenRole, " expected role:", userRole)
			fmt.Fprint(w, `
				<html>
				<head>
					<title>Unauthorized</title>
					<script>
						setTimeout(function() {
							window.location.href = '/teacher_signin';
						}, 2000);
					</script>
				</head>
				<body>
					<p>Unauthorized access. Redirecting to login page...</p>
				</body>
				</html>
			`)
			return
		}
		uid, err := strconv.Atoi(r.FormValue("uid"))
		msg := ""
		if err == nil {
			ok := false
			var password string
			if givenRole == "teacher" {
				password, ok = models.TeacherMap[uid]
				if ok && password != r.FormValue("password") {
					ok = false
					msg += r.FormValue("uid") + " (TeacherMap): Password doesn't match. "
				}
				if !ok {
					c, err := r.Cookie("session_token")
					if err != nil {
						ok = false
						msg += "No Session token exists. "
					} else {
						sessionToken := c.Value
						userSession, exists := models.Sessions[sessionToken]
						if !exists || userSession.IsExpired() {
							ok = false
							msg += r.FormValue("uid") + " (TeacherMap): No Session token exists or Session expired. "
						} else {
							ok = true
						}
					}
				}
			} else {
				_, ok = models.Students[uid]
				if !ok {
					ok = repository.LoadAndAuthorizeStudent(uid, r.FormValue("password"))
				} else if models.Students[uid].Password != r.FormValue("password") {
					ok = false
					msg += models.Students[uid].Name + "(Student): Password doesn't match. "
				}
			}
			if ok {
				fn(w, r, r.FormValue("name"), uid)
				return
			}
		}
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("Unauthorized access:", r.FormValue("name"), msg)
		fmt.Fprint(w, `
			<html>
			<head>
				<title>Unauthorized</title>
				<script>
					setTimeout(function() {
						window.location.href = '/teacher_signin';
					}, 2000);
				</script>
			</head>
			<body>
				<p>Unauthorized access. Redirecting to login page...</p>
			</body>
			</html>
		`)
	}
}

// -----------------------------------------------------------------
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get session token from the cookie
	cookie, err := r.Cookie("session_token")
	if err == nil {
		// Delete the session if it exists
		delete(models.Sessions, cookie.Value)
	}

	// Expire the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour), // Expire the cookie
		Path:    "/",
	})

	// Redirect to login page
	http.Redirect(w, r, "/teacher_login", http.StatusSeeOther)
}

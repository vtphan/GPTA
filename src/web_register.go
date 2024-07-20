package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

func (app App) webRegisterHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if email == "" || !validateEmail(email) {
		http.Error(w, fmt.Sprintf("Invalid email address %s.", email), http.StatusBadRequest)
		log.Printf("Invalid email address %s", email)
		return
	}
	tag := randomString(20)
	tempTag := &TempUserTag{Email: email, Tag: tag}
	app.DB.Create(&tempTag)
	err := app.Mailer.SendEmail(email, "Complete your GPTA registration", fmt.Sprintf("Please complete registration at the following link. %s/complete_registration?verify=%s", Config.Address, tag))
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
	} else {
		http.Redirect(w, r, "register_email_sent", http.StatusFound)
	}
}

func (app App) webRegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	temp := template.New("")
	t, err := temp.Parse(REGISTER)
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

func (app App) webRegisterEmailSentCompletePageHandler(w http.ResponseWriter, r *http.Request) {
	temp := template.New("")
	t, err := temp.Parse(EMAIL_SENT)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}

func (app App) webRegiserCompleteRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	tag := r.FormValue("token")
	var tempUser TempUserTag
	result := app.DB.Where("tag = ?", tag).First(&tempUser)
	if result.Error != nil {
		log.Print(result.Error)
		http.Error(w, "User not found.", http.StatusNotFound)
		return
	}
	temp := template.New("")
	t, err := temp.Parse(GIVE_PASSWORD)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, tempUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

// validateEmail function checks if the email format is valid
func validateEmail(email string) bool {
	// Define a regular expression for validating email addresses
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

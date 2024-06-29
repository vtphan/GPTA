package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func (app App) webRegisterHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	tag := randomString(20)
	tempTag := &TempUserTag{Email: email, Tag: tag}
	app.DB.Create(&tempTag)
	err := app.Mailer.SendEmail(email, "Complete your GPTA registration", fmt.Sprintf("Please complete registration at the following link. %s/complete_register", Config.Address))
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
		log.Fatal(err)
	}
}

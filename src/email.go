package main

import (
	"gopkg.in/gomail.v2"
)

// EmailSender struct to hold email configuration
type EmailSender struct {
	SMTPServer string
	Port       int
	Email      string
	Password   string
}

// NewEmailSender is a constructor function for EmailSender
func NewEmailSender(smtpServer string, port int, email, password string) *EmailSender {
	return &EmailSender{
		SMTPServer: smtpServer,
		Port:       port,
		Email:      email,
		Password:   password,
	}
}

// SendEmail method to send an email
func (es *EmailSender) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", es.Email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(es.SMTPServer, es.Port, es.Email, es.Password)

	return d.DialAndSend(m)
}

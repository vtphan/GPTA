package main

import (
	"crypto/md5"
	"encoding/hex"

	"gorm.io/gorm"
)

// App represents the application with a database connection.
type App struct {
	DB       *gorm.DB
	Mailer   EmailSender
	Sessions map[string]session
}

func (app App) getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

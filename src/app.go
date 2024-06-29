package main

import (
	"gorm.io/gorm"
)

// App represents the application with a database connection.
type App struct {
	DB *gorm.DB
}

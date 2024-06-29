package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func createDatabaseIfNotExists(dbUser, dbPassword, dbHost, dbName string) error {
	// Create a connection string without specifying the database
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/", dbUser, dbPassword, dbHost)

	// Open a connection to MySQL server
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// Execute the CREATE DATABASE statement
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	return err
}

func (app *App) migrateSchema(c Configuration) {
	createDatabaseIfNotExists(c.DBUserName, c.DBPassWord, c.DBServerIP, c.Database)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.DBUserName, c.DBUserName, c.DBServerIP, c.DBServerPort, c.Database)
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	// Auto-migrate the schema
	DB.AutoMigrate(&User{}, &Attendance{}, &Tag{}, &Problem{}, &Submission{}, &Score{},
		&Feedback{}, &TestCase{}, &CodeExplanation{}, &HelpMessage{}, &CodeSnapshot{},
		&SnapshotFeedback{}, &SnapshotBackFeedback{}, &Message{}, &MessageFeedback{},
		&MessageBackFeedback{}, &HelpEligible{}, &UserEventLog{}, &StudentStatus{}, &ProblemStatistics{})
}

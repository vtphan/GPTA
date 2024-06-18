package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func migrateSchema(c Configuration) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.DBUserName, c.DBUserName, c.DBServerIP, c.DBServerPort, c.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Auto-migrate the schema
	db.AutoMigrate(&User{}, &Attendance{}, &Tag{}, &Problem{}, &Submission{}, &Score{}, &Feedback{}, &TestCase{}, &CodeExplanation{}, &HelpMessage{}, &CodeSnapshot{}, &SnapshotFeedback{}, &SnapshotBackFeedback{}, &Message{}, &MessageFeedback{}, &MessageBackFeedback{}, &HelpEligible{}, &UserEventLog{}, &StudentStatus{}, &ProblemStatistics{})
}

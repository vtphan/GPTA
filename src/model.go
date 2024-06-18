package main

import (
	"time"

	"gorm.io/gorm"
)

// User model
type User struct {
	gorm.Model
	ID       int    `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"size:100;unique"`
	Password string `gorm:"size:100"`
	Role     string `gorm:"size:50"` // Role can be "student" or "teacher"
	Email    string
	// Define relationships
	Attendances           []Attendance           `gorm:"foreignKey:StudentID"`
	Submissions           []Submission           `gorm:"foreignKey:StudentID"`
	Scores                []Score                `gorm:"foreignKey:StudentID"`
	Feedbacks             []Feedback             `gorm:"foreignKey:AuthorID"`
	TestCases             []TestCase             `gorm:"foreignKey:StudentID"`
	CodeSnapshots         []CodeSnapshot         `gorm:"foreignKey:StudentID"`
	CodeExplanations      []CodeExplanation      `gorm:"foreignKey:StudentID"`
	SnapshotFeedbacks     []SnapshotFeedback     `gorm:"foreignKey:AuthorID"`
	SnapshotBackFeedbacks []SnapshotBackFeedback `gorm:"foreignKey:AuthorID"`
	Messages              []Message              `gorm:"foreignKey:AuthorID"`
	MessageFeedbacks      []MessageFeedback      `gorm:"foreignKey:AuthorID"`
	MessageBackFeedbacks  []MessageBackFeedback  `gorm:"foreignKey:AuthorID"`
	Problems              []Problem              `gorm:"foreignKey:TeacherID"`
}

// Attendance model
type Attendance struct {
	gorm.Model
	ID           int `gorm:"primaryKey;autoIncrement"`
	StudentID    int `gorm:"not null"`
	AttendanceAt time.Time
	// Define foreign key
	Student User `gorm:"foreignKey:StudentID"`
}

// Tag model
type Tag struct {
	gorm.Model
	ID               int    `gorm:"primaryKey;autoIncrement"`
	TopicDescription string `gorm:"size:200;unique"`
}

// Problem model
type Problem struct {
	gorm.Model
	ID                 int `gorm:"primaryKey;autoIncrement"`
	TeacherID          int
	ProblemDescription string
	Answer             string
	Filename           string
	Merit              int
	Effort             int
	Attempts           int
	TopicID            int
	TagID              int
	ProblemUploadedAt  time.Time
	ProblemEndedAt     time.Time
	// Define foreign key
	Teacher User `gorm:"foreignKey:TeacherID"`
	Tag     Tag  `gorm:"foreignKey:TagID"`
}

// Submission model
type Submission struct {
	gorm.Model
	ID                 int `gorm:"primaryKey;autoIncrement"`
	ProblemID          int `gorm:"not null"`
	StudentID          int `gorm:"not null"`
	StudentCode        string
	SnapshotID         int `gorm:"default:0"`
	SubmissionCategory int
	CodeSubmittedAt    time.Time
	Completed          time.Time
	Verdict            string
	AttemptNumber      int
	Answer             string
	// Define foreign keys
	Problem  Problem      `gorm:"foreignKey:ProblemID"`
	Student  User         `gorm:"foreignKey:StudentID"`
	Snapshot CodeSnapshot `gorm:"foreignKey:SnapshotID"`
}

// Score model
type Score struct {
	gorm.Model
	ID                     int `gorm:"primaryKey;autoIncrement"`
	ProblemID              int `gorm:"not null"`
	StudentID              int
	TeacherID              int
	Score                  int
	GradedSubmissionNumber int
	ScoreGivenAt           time.Time
	// Define foreign keys
	Problem Problem `gorm:"foreignKey:ProblemID"`
	Student User    `gorm:"foreignKey:StudentID"`
	Teacher User    `gorm:"foreignKey:TeacherID"`
}

// Feedback model
type Feedback struct {
	gorm.Model
	ID              int `gorm:"primaryKey;autoIncrement"`
	AuthorID        int
	Feedback        string
	FeedbackGivenAt time.Time
	SubmissionID    int
	// Define foreign keys
	Author User `gorm:"foreignKey:AuthorID"`
}

// TestCase model
type TestCase struct {
	gorm.Model
	ID        int `gorm:"primaryKey;autoIncrement"`
	ProblemID int
	StudentID int
	TestCases string
	AddedAt   time.Time
	// Define foreign keys
	Problem Problem `gorm:"foreignKey:ProblemID"`
	Student User    `gorm:"foreignKey:StudentID"`
}

// CodeExplanation model
type CodeExplanation struct {
	gorm.Model
	ID              int `gorm:"primaryKey;autoIncrement"`
	ProblemID       int
	StudentID       int
	SnapshotID      int
	TryingWhat      string
	NeedHelpWith    string
	CodeSubmittedAt time.Time
	// Define foreign keys
	Problem  Problem      `gorm:"foreignKey:ProblemID"`
	Student  User         `gorm:"foreignKey:StudentID"`
	Snapshot CodeSnapshot `gorm:"foreignKey:SnapshotID"`
}

// HelpMessage model
type HelpMessage struct {
	gorm.Model
	ID                int `gorm:"primaryKey;autoIncrement"`
	CodeExplanationID int
	StudentID         int
	Message           string
	GivenAt           time.Time
	Useful            string
	UpdatedAt         time.Time
	// Define foreign keys
	CodeExplanation CodeExplanation `gorm:"foreignKey:CodeExplanationID"`
	Student         User            `gorm:"foreignKey:StudentID"`
}

// CodeSnapshot model
type CodeSnapshot struct {
	gorm.Model
	ID            int `gorm:"primaryKey;autoIncrement"`
	StudentID     int
	ProblemID     int
	Code          string
	LastUpdatedAt time.Time
	Status        int    `gorm:"default:0"`
	Event         string `gorm:"size:50"`
	// Define foreign keys
	Student User    `gorm:"foreignKey:StudentID"`
	Problem Problem `gorm:"foreignKey:ProblemID"`
}

// SnapshotFeedback model
type SnapshotFeedback struct {
	gorm.Model
	ID         int `gorm:"primaryKey;autoIncrement"`
	SnapshotID int
	Feedback   string
	AuthorID   int
	GivenAt    time.Time
	// Define foreign keys
	Snapshot CodeSnapshot `gorm:"foreignKey:SnapshotID"`
	Author   User         `gorm:"foreignKey:AuthorID"`
}

// SnapshotBackFeedback model
type SnapshotBackFeedback struct {
	gorm.Model
	ID                 int `gorm:"primaryKey;autoIncrement"`
	SnapshotFeedbackID int
	AuthorID           int
	IsHelpful          string
	GivenAt            time.Time
	// Define foreign keys
	SnapshotFeedback SnapshotFeedback `gorm:"foreignKey:SnapshotFeedbackID"`
	Author           User             `gorm:"foreignKey:AuthorID"`
}

// Message model
type Message struct {
	gorm.Model
	ID         int `gorm:"primaryKey;autoIncrement"`
	SnapshotID int
	Message    string
	AuthorID   int
	GivenAt    time.Time
	Type       int
	// Define foreign keys
	Snapshot CodeSnapshot `gorm:"foreignKey:SnapshotID"`
	Author   User         `gorm:"foreignKey:AuthorID"`
}

// MessageFeedback model
type MessageFeedback struct {
	gorm.Model
	ID        int `gorm:"primaryKey;autoIncrement"`
	MessageID int
	Feedback  string
	AuthorID  int
	GivenAt   time.Time
	// Define foreign keys
	Message Message `gorm:"foreignKey:MessageID"`
	Author  User    `gorm:"foreignKey:AuthorID"`
}

// MessageBackFeedback model
type MessageBackFeedback struct {
	gorm.Model
	ID                int `gorm:"primaryKey;autoIncrement"`
	MessageFeedbackID int
	AuthorID          int
	Useful            string
	GivenAt           time.Time
	// Define foreign keys
	MessageFeedback MessageFeedback `gorm:"foreignKey:MessageFeedbackID"`
	Author          User            `gorm:"foreignKey:AuthorID"`
}

// HelpEligible model
type HelpEligible struct {
	gorm.Model
	ID               int `gorm:"primaryKey;autoIncrement"`
	ProblemID        int
	StudentID        int
	BecameEligibleAt time.Time
	// Define foreign keys
	Problem Problem `gorm:"foreignKey:ProblemID"`
	Student User    `gorm:"foreignKey:StudentID"`
}

// UserEventLog model
type UserEventLog struct {
	gorm.Model
	ID           int    `gorm:"primaryKey;autoIncrement"`
	Name         string `gorm:"size:50"`
	UserID       int
	EventType    string `gorm:"size:50"`
	ReferralInfo string `gorm:"size:50"`
	EventTime    time.Time
	// Define foreign keys
	User User `gorm:"foreignKey:UserID"`
}

// StudentStatus model
type StudentStatus struct {
	gorm.Model
	ID             int `gorm:"primaryKey;autoIncrement"`
	StudentID      int
	ProblemID      int
	CodingStat     string `gorm:"size:50"`
	HelpStat       string `gorm:"size:50"`
	SubmissionStat string `gorm:"size:50"`
	TutoringStat   string `gorm:"size:50"`
	LastUpdatedAt  time.Time
	// Define foreign keys
	Student User    `gorm:"foreignKey:StudentID"`
	Problem Problem `gorm:"foreignKey:ProblemID"`
}

// ProblemStatistics model
type ProblemStatistics struct {
	gorm.Model
	ID              int `gorm:"primaryKey;autoIncrement"`
	ProblemID       int `gorm:"not null"`
	Active          int `gorm:"default:0"`
	Submission      int `gorm:"default:0"`
	HelpRequest     int `gorm:"default:0"`
	GradedCorrect   int `gorm:"default:0"`
	GradedIncorrect int `gorm:"default:0"`
	// Define foreign keys
	Problem Problem `gorm:"foreignKey:ProblemID"`
}

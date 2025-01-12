package main

import "time"

// Student represents the student table
type Student struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"size:100;unique"`
	Password string `gorm:"size:100"`
}

// Teacher represents the teacher table
type Teacher struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"size:100;unique"`
	Password string `gorm:"size:100"`
}

// Attendance represents the attendance table
type Attendance struct {
	ID           int `gorm:"primaryKey;autoIncrement"`
	StudentID    int `gorm:"not null"`
	AttendanceAt time.Time
}

// Tag represents the tag table
type Tag struct {
	ID               int    `gorm:"primaryKey;autoIncrement"`
	TopicDescription string `gorm:"size:200;unique"`
}

// Problem represents the problem table
type Problem struct {
	ID                 int `gorm:"primaryKey;autoIncrement"`
	TeacherID          int
	ProblemDescription string `gorm:"type:text"`
	Answer             string `gorm:"type:text"`
	Filename           string
	Merit              int
	Effort             int
	Attempts           int
	TopicID            int
	Tag                int
	ProblemUploadedAt  time.Time
	ProblemEndedAt     *time.Time `gorm:"default:null"`
}

// Submission represents the submission table
type Submission struct {
	ID                 int    `gorm:"primaryKey;autoIncrement"`
	ProblemID          int    `gorm:"not null"`
	StudentID          int    `gorm:"not null"`
	StudentCode        string `gorm:"type:text"`
	SnapshotID         int    `gorm:"default:0"`
	SubmissionCategory int
	CodeSubmittedAt    time.Time
	Completed          time.Time
	Verdict            string
	AttemptNumber      int
	Answer             string
}

// Score represents the score table
type Score struct {
	ID                     int `gorm:"primaryKey;autoIncrement"`
	ProblemID              int `gorm:"not null"`
	StudentID              int
	TeacherID              int
	Score                  int
	GradedSubmissionNumber int
	ScoreGivenAt           *time.Time
	Problem                Problem `gorm:"foreignKey:ProblemID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// Feedback represents the feedback table
type Feedback struct {
	ID              int `gorm:"primaryKey;autoIncrement"`
	TeacherID       int
	StudentID       int
	Feedback        string `gorm:"type:text"`
	FeedbackGivenAt time.Time
	SubmissionID    int
}

// TestCase represents the test_case table
type TestCase struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	ProblemID int
	StudentID int
	TestCases string `gorm:"type:text"`
	AddedAt   time.Time
}

// CodeExplanation represents the code_explanation table
type CodeExplanation struct {
	ID              int `gorm:"primaryKey;autoIncrement"`
	ProblemID       int
	StudentID       int
	SnapshotID      int
	TryingWhat      string `gorm:"type:text"`
	NeedHelpWith    string `gorm:"type:text"`
	CodeSubmittedAt time.Time
}

// HelpMessage represents the help_message table
type HelpMessage struct {
	ID                int `gorm:"primaryKey;autoIncrement"`
	CodeExplanationID int
	StudentID         int
	Message           string `gorm:"type:text"`
	GivenAt           time.Time
	Useful            string
	UpdatedAt         time.Time
}

// CodeSnapshot represents the code_snapshot table
type CodeSnapshot struct {
	ID            int `gorm:"primaryKey;autoIncrement"`
	StudentID     int
	ProblemID     int
	Code          string `gorm:"type:text"`
	LastUpdatedAt time.Time
	Status        int    `gorm:"default:0"`
	Event         string `gorm:"size:50"`
}

// SnapshotFeedback represents the snapshot_feedback table
type SnapshotFeedback struct {
	ID         int `gorm:"primaryKey;autoIncrement"`
	SnapshotID int
	Feedback   string `gorm:"type:text"`
	AuthorID   int
	AuthorRole string `gorm:"size:50"`
	GivenAt    time.Time
}

// SnapshotBackFeedback represents the snapshot_back_feedback table
type SnapshotBackFeedback struct {
	ID                 int `gorm:"primaryKey;autoIncrement"`
	SnapshotFeedbackID int
	AuthorID           int
	AuthorRole         string `gorm:"size:50"`
	IsHelpful          string `gorm:"size:50"`
	GivenAt            time.Time
}

// Message represents the message table
type Message struct {
	ID         int `gorm:"primaryKey;autoIncrement"`
	SnapshotID int
	Message    string `gorm:"type:text"`
	AuthorID   int
	AuthorRole string `gorm:"size:50"`
	GivenAt    time.Time
	Type       int
}

// MessageFeedback represents the message_feedback table
type MessageFeedback struct {
	ID         int `gorm:"primaryKey;autoIncrement"`
	MessageID  int
	Feedback   string `gorm:"type:text"`
	AuthorID   int
	AuthorRole string `gorm:"size:50"`
	GivenAt    time.Time
}

// MessageBackFeedback represents the message_back_feedback table
type MessageBackFeedback struct {
	ID                int `gorm:"primaryKey;autoIncrement"`
	MessageFeedbackID int
	AuthorID          int
	AuthorRole        string `gorm:"size:50"`
	Useful            string `gorm:"size:50"`
	GivenAt           time.Time
}

// HelpEligible represents the help_eligible table
type HelpEligible struct {
	ID               int `gorm:"primaryKey;autoIncrement"`
	ProblemID        int
	StudentID        int
	BecameEligibleAt time.Time
}

// UserEventLog represents the user_event_log table
type UserEventLog struct {
	ID           int    `gorm:"primaryKey;autoIncrement"`
	Name         string `gorm:"size:50"`
	UserID       int
	UserType     string `gorm:"size:50"`
	EventType    string `gorm:"size:50"`
	ReferralInfo string `gorm:"size:50"`
	EventTime    time.Time
}

// StudentStatus represents the student_status table
type StudentStatus struct {
	ID             int `gorm:"primaryKey;autoIncrement"`
	StudentID      int
	ProblemID      int
	CodingStat     string `gorm:"size:50"`
	HelpStat       string `gorm:"size:50"`
	SubmissionStat string `gorm:"size:50"`
	TutoringStat   string `gorm:"size:50"`
	LastUpdatedAt  time.Time
}

// ProblemStatistics represents the problem_statistics table
type ProblemStatistics struct {
	ID              int `gorm:"primaryKey;autoIncrement"`
	ProblemID       int `gorm:"not null"`
	Active          int `gorm:"default:0"`
	Submission      int `gorm:"default:0"`
	HelpRequest     int `gorm:"default:0"`
	GradedCorrect   int `gorm:"default:0"`
	GradedIncorrect int `gorm:"default:0"`
}

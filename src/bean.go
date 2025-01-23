package main

import "time"

type Student struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"size:100;unique"`
	Password string `gorm:"size:100"`
}

type Teacher struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"size:100;unique"`
	Password string `gorm:"size:100"`
}

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

type Feedback struct {
	ID              int `gorm:"primaryKey;autoIncrement"`
	TeacherID       int
	StudentID       int
	Feedback        string `gorm:"type:text"`
	FeedbackGivenAt time.Time
	SubmissionID    int
}

type Attendance struct {
	ID           int `gorm:"primaryKey;autoIncrement"`
	StudentID    int `gorm:"not null"`
	AttendanceAt time.Time
}

type Tag struct {
	ID               int    `gorm:"primaryKey;autoIncrement"`
	TopicDescription string `gorm:"size:200;unique"`
}
type TestCase struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	ProblemID int
	StudentID int
	TestCases string `gorm:"type:text"`
	AddedAt   time.Time
}

type HelpMessage struct {
	ID                int `gorm:"primaryKey;autoIncrement"`
	CodeExplanationID int
	StudentID         int
	Message           string `gorm:"type:text"`
	GivenAt           time.Time
	Useful            string
	UpdatedAt         time.Time
}

type CodeSnapshot struct {
	ID            int `gorm:"primaryKey;autoIncrement"`
	StudentID     int
	ProblemID     int
	Code          string `gorm:"type:text"`
	LastUpdatedAt time.Time
	Status        int    `gorm:"default:0"`
	Event         string `gorm:"size:50"`
}

type HelpEligible struct {
	ID               int `gorm:"primaryKey;autoIncrement"`
	ProblemID        int
	StudentID        int
	BecameEligibleAt time.Time
}

type UserEventLog struct {
	ID           int    `gorm:"primaryKey;autoIncrement"`
	Name         string `gorm:"size:50"`
	UserID       int
	UserType     string `gorm:"size:50"`
	EventType    string `gorm:"size:50"`
	ReferralInfo string `gorm:"size:50"`
	EventTime    time.Time
}
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
type Message struct {
	ID         int `gorm:"primaryKey;autoIncrement"`
	SnapshotID int
	Message    string `gorm:"type:text"`
	AuthorID   int
	AuthorRole string `gorm:"size:50"`
	GivenAt    time.Time
	Type       int
}
type MessageFeedback struct {
	ID         int `gorm:"primaryKey;autoIncrement"`
	MessageID  int
	Feedback   string `gorm:"type:text"`
	AuthorID   int
	AuthorRole string `gorm:"size:50"`
	GivenAt    time.Time
}

type ProblemStatistics struct {
	ID              int `gorm:"primaryKey;autoIncrement"`
	ProblemID       int `gorm:"not null"`
	Active          int `gorm:"default:0"`
	Submission      int `gorm:"default:0"`
	HelpRequest     int `gorm:"default:0"`
	GradedCorrect   int `gorm:"default:0"`
	GradedIncorrect int `gorm:"default:0"`
}

type MessageBackFeedback struct {
	ID                int `gorm:"primaryKey;autoIncrement"`
	MessageFeedbackID int
	AuthorID          int
	AuthorRole        string `gorm:"size:50"`
	Useful            string `gorm:"size:50"`
	GivenAt           time.Time
}

type SnapshotBackFeedback struct {
	ID                 int `gorm:"primaryKey;autoIncrement"`
	SnapshotFeedbackID int
	AuthorID           int
	AuthorRole         string `gorm:"size:50"`
	IsHelpful          string `gorm:"size:50"`
	GivenAt            time.Time
}

type SnapshotFeedback struct {
	ID         int `gorm:"primaryKey;autoIncrement"`
	SnapshotID int
	Feedback   string `gorm:"type:text"`
	AuthorID   int
	AuthorRole string `gorm:"size:50"`
	GivenAt    time.Time
}

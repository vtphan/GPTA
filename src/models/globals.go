// Author: Vinhthuy Phan, 2018
package models

import (
	"bufio"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

// ---------------------------------------------------------
type Configuration struct {
	CourseName string // todo - set this globally once the user selects the course
	IP         string
	Port       int
	Database   string
	DBServerIP string
	DBUserName string
	DBPassWord string
	Address    string
	LogFile    string
	PeerTutor  int
}

var Config *Configuration
var CourseId string

//---------------------------------------------------------
// Database
//---------------------------------------------------------

var DB *gorm.DB

const RecordNotFound = "record not found"

//---------------------------------------------------------
// Authentication
//---------------------------------------------------------

var TeacherMap = make(map[int]string)
var TeacherPass = make(map[string]string)
var TeacherClassesMap = make(map[int][]string)
var StudentClassesMap = make(map[int][]string)
var TeacherNameToId = make(map[string]int)
var TeacherIdToName = make(map[int]string)
var AdminPass = map[string]string{
	"admin": "123456", // Replace with real data
}
var Passcode string

//---------------------------------------------------------
// Semaphores
//---------------------------------------------------------

var BoardsSem sync.Mutex
var SubSem sync.Mutex
var BulletinSem sync.Mutex
var HelpSubSem sync.Mutex
var CodeSnapshotSem sync.Mutex

//---------------------------------------------------------
// Virtual boards for students and student submissions
//---------------------------------------------------------

type Board struct {
	Content      string
	Answer       string
	Attempts     int
	Filename     string
	Pid          int // problem id
	StartingTime time.Time
	Type         string
}
type StudentSubmissionStatus struct {
	Filename      string
	AttemptNumber int
	Status        int
	/*
		1 submission being looked at.
		2 teacher did not grade your submission (dismissed).
		3 your submission was not correct.
		4 your submission was correct.
	*/
}
type SnapShotFeedback struct {
	FeedbackID  int
	Snapshot    string
	Feedback    string
	ProblemName string
	Provider    string
}
type StudenInfo struct {
	Name                  string
	Password              string
	Boards                []*Board
	SubmissionStatus      []*StudentSubmissionStatus
	SnapShotFeedbackQueue []*SnapShotFeedback
	ThankStatus           int
	/*
		0 Nothing
		1 Got a new thanks for feedback
	*/
}

var Students = make(map[int]*StudenInfo)

//---------------------------------------------------------

var BulletinBoard = make([]string, 0)

// ---------------------------------------------------------
type Submission struct {
	Sid           int // submission id
	Uid           int // student id
	Pid           int // problem id
	Content       string
	Filename      string
	Priority      int
	AttemptNumber int
	At            time.Time
	Name          string
	SnapshotID    int
}

var WorkingSubs = make([]*Submission, 0)
var Submissions = make(map[int]*Submission)

//---------------------------------------------------------

type HelpSubmission struct {
	Sid        int // submission id
	Uid        int // student id
	Pid        int // problem id
	Status     int // 0=ok, 1=queue empty, 2=not elligible
	Content    string
	Filename   string
	At         time.Time
	SnapshotID int
	Snapshot   string
}

var WorkingHelpSubs = make([]*HelpSubmission, 0)
var HelpSubmissions = make(map[int]*HelpSubmission)

// ---------------------------------------------------------
type ProblemInfo struct {
	Description string
	Filename    string
	Answer      string
	Merit       int
	Effort      int
	Attempts    int
	Topic_id    int
	Tag         string
	Pid         int
	ExactAnswer bool
}

type ActiveProblem struct {
	Info     *ProblemInfo
	Answers  []string
	Active   bool
	Attempts map[int]int
}

type CodeSnapshotMessageDetails struct {
	StudentID   int    `gorm:"column:student_id"`
	ProblemID   int    `gorm:"column:problem_id"`
	Code        string `gorm:"column:code"`
	Filename    string `gorm:"column:filename"`
	MessageType string `gorm:"column:message_type"`
}

var ActiveProblems = make(map[string]*ActiveProblem)

//---------------------------------------------------------
// Utilities
//---------------------------------------------------------

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// -----------------------------------------------------------------------------
func WriteLog(filename, message string) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println(time.Now(), " ", message)
}

//---------------------------------------------------------

var HelpEligibleStudents = map[int]map[int]bool{}
var SeenHelpSubmissions = map[int]map[int]bool{}

//---------------------------------------------------------

// Snapshot contains information related to a code snapshot.
type Snapshot struct {
	ID          int
	StudentName string
	StudentID   int
	ProblemName string
	ProblemID   int
	Status      int
	FirstUpdate time.Time
	LastUpdated time.Time
	LinesOfCode int
	Code        string
	NumFeedback int
}

// Snapshots contains all the current snapshots from students.
var Snapshots = make([]*Snapshot, 0)

// StudentSnapshot is the mapping from (student id, problem id) -> snashpt index in `Snapshots` list.
var StudentSnapshot = map[int]map[int]int{}

// SnapshotStatus maps from integer status to string named status.
var SnapshotStatus = []string{"Not submitted", "Submitted: not graded", "Submitted: incorrect", "Submitted: correct"}

// SnapshotStatusMapping maps from string snapshot status to integer.
var SnapshotStatusMapping = map[string]int{
	"Not submitted":         0,
	"Submitted: not graded": 1,
	"Submitted: incorrect":  2,
	"Submitted: correct":    3,
}

func GetLinesOfCode(code string) int {
	scanner := bufio.NewScanner(strings.NewReader(code))
	scanner.Split(bufio.ScanLines)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count
}

// PeerTutorAllowed is the flag that decides whether peers are allowed to help other students or not.
var PeerTutorAllowed = false

type DashBoardStudentInfo struct {
	StudentID      int
	StudentName    string
	LastUpdatedAt  time.Time
	CodingStat     string
	HelpStat       string
	SubmissionStat string
	TutoringStat   string
}

type AnswerStatInfo struct {
	Answer  string
	Count   int
	Percent float64
}

type DashBoardInfo struct {
	StudentInfo        []*DashBoardStudentInfo
	ProblemName        string
	Code               string
	IsActive           bool
	ProblemID          int
	NumActive          int
	NumHelpRequest     int
	NumGradedCorrect   int
	NumGradedIncorrect int
	NumNotGraded       int
	AnswerStats        []*AnswerStatInfo
	UserID             int
	UserRole           string
	Password           string
	Username           string
}

type FeedbackDashBoard struct {
	Name            string
	Role            string
	Feedback        string
	FeedbackID      int
	CurrentUserVote string
	Downvote        int
	Upvote          int
	GivenAt         time.Time
}

type MessageDashBoard struct {
	ID         int
	Name       string
	Role       string
	Message    string
	Type       int // 0 = help request, 1 = unsolicited
	Event      string
	GivenAt    time.Time
	Code       string
	SnapshotID int
	Feedbacks  []*FeedbackDashBoard
}

type FeedbackProvisionDashBoard struct {
	StudentName  string
	ProblemName  string
	Status       DashBoardStudentInfo
	LastSnapshot *Snapshot
	Messages     []*MessageDashBoard
	StudentID    int
	ProblemID    int
	UserID       int
	UserRole     string
	Password     string
	Username     string
}

type SubmissionInfo struct {
	ID          int
	Code        string
	Grade       string
	SubmittedAt time.Time
	SnapshotID  int
}

type SubmissionDashboard struct {
	Submissions []*SubmissionInfo
	StudentName string
	ProblemName string
	StudentID   int
	ProblemID   int
	UserID      int
	UserRole    string
	Password    string
	Username    string
}

type TemplateDate struct {
	Feedback   FeedbackProvisionDashBoard
	Submission SubmissionDashboard
	Status     DashBoardStudentInfo
	UserID     int
	UserRole   string
	Password   string
	Username   string
	CourseName string
}

type StudentReport struct {
	Points   int
	Filename string
	Date     int64
}

type ScoreEntry struct {
	Name     string
	Points   int
	Attempts int
	Count    int
}

type TagsViewData struct {
	Tags            map[int]string
	SubmissionCount map[string]int
	Scores          map[int]*ScoreEntry
	PC              string
}

type ProblemPerformance struct {
	Pid       int
	Timestamp int64
	Correct   int
	Incorrect int
	Activity  float32
	Success   float32
	PC        string
}

type TagData struct {
	Description string
	Performance map[int]*ProblemPerformance
}

type MessageWithSnapshot struct {
	MessageID   int       `gorm:"column:id"`
	SnapshotID  int       `gorm:"column:snapshot_id"`
	Message     string    `gorm:"column:message"`
	AuthorID    int       `gorm:"column:author_id"`
	AuthorRole  string    `gorm:"column:author_role"`
	GivenAt     time.Time `gorm:"column:given_at"`
	MessageType int       `gorm:"column:type"`
	Code        string    `gorm:"column:code"`
	Event       string    `gorm:"column:event"`
}

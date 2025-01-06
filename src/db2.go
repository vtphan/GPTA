package main

import (
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"time"
)

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
	ProblemEndedAt     time.Time
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

//func InitDatabase(dbName string, username string, pass string, server string) {
//	var err error
//	DB, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:3306)/", username, pass, server)), &gorm.Config{})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	DB.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
//	DB, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", username, pass, server, dbName)), &gorm.Config{})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Initialize your tables and structs here
//}

func AddStudent(name string, password string) error {
	student := Student{Name: name, Password: password}
	if err := DB.Create(&student).Error; err != nil {
		return err
	}
	return nil
}

func AddTeacher(name string, password string) error {
	teacher := Teacher{Name: name, Password: password}
	if err := DB.Create(&teacher).Error; err != nil {
		return err
	}
	return nil
}

func AddProblem(teacherID int, problemDescription, answer, filename string, merit, effort, attempts int, topicID int, tag int, problemUploadedAt time.Time) error {
	problem := Problem{
		TeacherID:          teacherID,
		ProblemDescription: problemDescription,
		Answer:             answer,
		Filename:           filename,
		Merit:              merit,
		Effort:             effort,
		Attempts:           attempts,
		TopicID:            topicID,
		Tag:                tag,
		ProblemUploadedAt:  problemUploadedAt,
	}
	if err := DB.Create(&problem).Error; err != nil {
		return err
	}
	return nil
}

func AddSubmission(problemID, studentID int, studentCode string, submissionCategory int, attemptNumber int, codeSubmittedAt time.Time, snapshotID int, answer string) error {
	submission := Submission{
		ProblemID:          problemID,
		StudentID:          studentID,
		StudentCode:        studentCode,
		SubmissionCategory: submissionCategory,
		AttemptNumber:      attemptNumber,
		CodeSubmittedAt:    codeSubmittedAt,
		SnapshotID:         snapshotID,
		Answer:             answer,
	}
	if err := DB.Create(&submission).Error; err != nil {
		return err
	}
	return nil
}

func CompleteSubmission(id int, completed bool, verdict string) error {
	if err := DB.Model(&Submission{}).Where("id = ?", id).Updates(map[string]interface{}{"Completed": completed, "Verdict": verdict}).Error; err != nil {
		return err
	}
	return nil
}

func AddScore(problemID, studentID, teacherID int, score int, gradedSubmissionNumber int, scoreGivenAt *time.Time) error {
	scoreEntry := Score{
		ProblemID:              problemID,
		StudentID:              studentID,
		TeacherID:              teacherID,
		Score:                  score,
		GradedSubmissionNumber: gradedSubmissionNumber,
		ScoreGivenAt:           scoreGivenAt,
	}
	if err := DB.Create(&scoreEntry).Error; err != nil {
		return err
	}
	return nil
}

func AddFeedback(teacherID, studentID int, feedback string, feedbackGivenAt time.Time, submissionID int) error {
	feedbackEntry := Feedback{
		TeacherID:       teacherID,
		StudentID:       studentID,
		Feedback:        feedback,
		FeedbackGivenAt: feedbackGivenAt,
		SubmissionID:    submissionID,
	}
	if err := DB.Create(&feedbackEntry).Error; err != nil {
		return err
	}
	return nil
}

func UpdateScore(id int, teacherID int, score float64, gradedSubmissionNumber int) error {
	if err := DB.Model(&Score{}).Where("id = ?", id).Updates(map[string]interface{}{
		"TeacherID":              teacherID,
		"Score":                  score,
		"GradedSubmissionNumber": gradedSubmissionNumber,
	}).Error; err != nil {
		return err
	}
	return nil
}

func AddAttendance(studentID int, attendanceAt time.Time) error {
	attendance := Attendance{
		StudentID:    studentID,
		AttendanceAt: attendanceAt,
	}
	if err := DB.Create(&attendance).Error; err != nil {
		return err
	}
	return nil
}

func AddTag(topicDescription string) error {
	tag := Tag{TopicDescription: topicDescription}
	if err := DB.Create(&tag).Error; err != nil {
		return err
	}
	return nil
}

func AddTestCase(problemID, studentID int, testCases string, addedAt time.Time) error {
	testCase := TestCase{
		ProblemID: problemID,
		StudentID: studentID,
		TestCases: testCases,
		AddedAt:   addedAt,
	}
	if err := DB.Create(&testCase).Error; err != nil {
		return err
	}
	return nil
}

func UpdateTestCase(id int, testCases string, addedAt time.Time) error {
	if err := DB.Model(&TestCase{}).Where("id = ?", id).Updates(map[string]interface{}{
		"TestCases": testCases,
		"AddedAt":   addedAt,
	}).Error; err != nil {
		return err
	}
	return nil
}

func AddHelpSubmission(problemID, studentID, snapshotID int, tryingWhat, needHelpWith string, codeSubmittedAt time.Time) error {
	helpSubmission := CodeExplanation{
		ProblemID:       problemID,
		StudentID:       studentID,
		SnapshotID:      snapshotID,
		TryingWhat:      tryingWhat,
		NeedHelpWith:    needHelpWith,
		CodeSubmittedAt: codeSubmittedAt,
	}
	if err := DB.Create(&helpSubmission).Error; err != nil {
		return err
	}
	return nil
}

func AddHelpMessage(codeExplanationID, studentID int, message string, givenAt time.Time) error {
	helpMessage := HelpMessage{
		CodeExplanationID: codeExplanationID,
		StudentID:         studentID,
		Message:           message,
		GivenAt:           givenAt,
	}
	if err := DB.Create(&helpMessage).Error; err != nil {
		return err
	}
	return nil
}

func UpdateHelpMessage(id int, useful bool, updatedAt time.Time) error {
	if err := DB.Model(&HelpMessage{}).Where("id = ?", id).Updates(map[string]interface{}{
		"Useful":    useful,
		"UpdatedAt": updatedAt,
	}).Error; err != nil {
		return err
	}
	return nil
}

func AddCodeSnapshot(studentID, problemID int, code string, status int, lastUpdatedAt time.Time, event string) error {
	codeSnapshot := CodeSnapshot{
		StudentID:     studentID,
		ProblemID:     problemID,
		Code:          code,
		Status:        status,
		LastUpdatedAt: lastUpdatedAt,
		Event:         event,
	}
	if err := DB.Create(&codeSnapshot).Error; err != nil {
		return err
	}
	return nil
}

func AddSnapShotFeedback(snapshotID int, feedback string, authorID int, authorRole string, givenAt time.Time) error {
	snapshotFeedback := SnapshotFeedback{
		SnapshotID: snapshotID,
		Feedback:   feedback,
		AuthorID:   authorID,
		AuthorRole: authorRole,
		GivenAt:    givenAt,
	}
	if err := DB.Create(&snapshotFeedback).Error; err != nil {
		return err
	}
	return nil
}

func AddSnapshotBackFeedback(snapshotFeedbackID int, authorID int, authorRole string, isHelpful string, givenAt time.Time) error {
	snapshotBackFeedback := SnapshotBackFeedback{
		SnapshotFeedbackID: snapshotFeedbackID,
		AuthorID:           authorID,
		AuthorRole:         authorRole,
		IsHelpful:          isHelpful,
		GivenAt:            givenAt,
	}
	if err := DB.Create(&snapshotBackFeedback).Error; err != nil {
		return err
	}
	return nil
}

func UpdateSnapshotBackFeedback(snapshotFeedbackID int, isHelpful bool, givenAt time.Time) error {
	if err := DB.Model(&SnapshotBackFeedback{}).Where("snapshot_feedback_id = ?", snapshotFeedbackID).Updates(map[string]interface{}{
		"IsHelpful": isHelpful,
		"GivenAt":   givenAt,
	}).Error; err != nil {
		return err
	}
	return nil
}

func UpdateProblemEndTime(problemID int, problemEndedAt time.Time) error {
	if err := DB.Model(&Problem{}).Where("id = ?", problemID).Update("ProblemEndedAt", problemEndedAt).Error; err != nil {
		return err
	}
	return nil
}

func AddHelpEligible(problemID, studentID int, becameEligibleAt time.Time) error {
	helpEligible := HelpEligible{
		ProblemID:        problemID,
		StudentID:        studentID,
		BecameEligibleAt: becameEligibleAt,
	}
	if err := DB.Create(&helpEligible).Error; err != nil {
		return err
	}
	return nil
}

func AddUserEventLog(userID int, name, userType, eventType, referralInfo string, eventTime time.Time) error {
	userEventLog := UserEventLog{
		Name:         name,
		UserID:       userID,
		UserType:     userType,
		EventType:    eventType,
		ReferralInfo: referralInfo,
		EventTime:    eventTime,
	}
	if err := DB.Create(&userEventLog).Error; err != nil {
		return err
	}
	return nil
}

func AddStudentStatus(studentID, problemID int, codingStat, helpStat, submissionStat, tutoringStat string, lastUpdatedAt time.Time) error {
	studentStatus := StudentStatus{
		StudentID:      studentID,
		ProblemID:      problemID,
		CodingStat:     codingStat,
		HelpStat:       helpStat,
		SubmissionStat: submissionStat,
		TutoringStat:   tutoringStat,
		LastUpdatedAt:  lastUpdatedAt,
	}
	if err := DB.Create(&studentStatus).Error; err != nil {
		return err
	}
	return nil
}

func UpdateStudentCodingStat(studentID, problemID int, codingStat string, lastUpdatedAt time.Time) error {
	if err := DB.Model(&StudentStatus{}).Where("student_id = ? AND problem_id = ?", studentID, problemID).Update("CodingStat", codingStat).Update("LastUpdatedAt", lastUpdatedAt).Error; err != nil {
		return err
	}
	return nil
}

func UpdateStudentHelpStat(studentID, problemID int, helpStat string, lastUpdatedAt time.Time) error {
	if err := DB.Model(&StudentStatus{}).Where("student_id = ? AND problem_id = ?", studentID, problemID).Update("HelpStat", helpStat).Update("LastUpdatedAt", lastUpdatedAt).Error; err != nil {
		return err
	}
	return nil
}

func UpdateStudentSubmissionStat(studentID, problemID int, submissionStat string, lastUpdatedAt time.Time) error {
	if err := DB.Model(&StudentStatus{}).Where("student_id = ? AND problem_id = ?", studentID, problemID).Update("SubmissionStat", submissionStat).Update("LastUpdatedAt", lastUpdatedAt).Error; err != nil {
		return err
	}
	return nil
}

func UpdateStudentTutoringStat(studentID, problemID int, tutoringStat string, lastUpdatedAt time.Time) error {
	if err := DB.Model(&StudentStatus{}).Where("student_id = ? AND problem_id = ?", studentID, problemID).Update("TutoringStat", tutoringStat).Update("LastUpdatedAt", lastUpdatedAt).Error; err != nil {
		return err
	}
	return nil
}

func AddMessage(snapshotID int, message string, authorID int, authorRole string, givenAt time.Time, messageType int) error {
	messageEntry := Message{
		SnapshotID: snapshotID,
		Message:    message,
		AuthorID:   authorID,
		AuthorRole: authorRole,
		GivenAt:    givenAt,
		Type:       messageType,
	}
	if err := DB.Create(&messageEntry).Error; err != nil {
		return err
	}
	return nil
}

func AddMessageFeedback(messageID int, feedback string, authorID int, authorRole string, givenAt time.Time) error {
	messageFeedback := MessageFeedback{
		MessageID:  messageID,
		Feedback:   feedback,
		AuthorID:   authorID,
		AuthorRole: authorRole,
		GivenAt:    givenAt,
	}
	if err := DB.Create(&messageFeedback).Error; err != nil {
		return err
	}
	return nil
}

func AddMessageBackFeedback(messageFeedbackID int, authorID int, authorRole string, useful string, givenAt time.Time) error {
	messageBackFeedback := MessageBackFeedback{
		MessageFeedbackID: messageFeedbackID,
		AuthorID:          authorID,
		AuthorRole:        authorRole,
		Useful:            useful,
		GivenAt:           givenAt,
	}
	if err := DB.Create(&messageBackFeedback).Error; err != nil {
		return err
	}
	return nil
}

func UpdateMessageBackFeedback(messageFeedbackID int, useful bool, givenAt time.Time) error {
	if err := DB.Model(&MessageBackFeedback{}).Where("message_feedback_id = ?", messageFeedbackID).Updates(map[string]interface{}{
		"Useful":  useful,
		"GivenAt": givenAt,
	}).Error; err != nil {
		return err
	}
	return nil
}

func AddProblemStatistics(problemID int) error {
	problemStats := ProblemStatistics{
		ProblemID:       problemID,
		Active:          0,
		Submission:      0,
		HelpRequest:     0,
		GradedCorrect:   0,
		GradedIncorrect: 0,
	}
	if err := DB.Create(&problemStats).Error; err != nil {
		return err
	}
	return nil
}

func IncProblemStatActive(problemID int) error {
	if err := DB.Model(&ProblemStatistics{}).Where("problem_id = ?", problemID).UpdateColumn("Active", gorm.Expr("active + ?", 1)).Error; err != nil {
		return err
	}
	return nil
}

func IncProblemStatSubmission(problemID int) error {
	if err := DB.Model(&ProblemStatistics{}).Where("problem_id = ?", problemID).UpdateColumn("Submission", gorm.Expr("submission + ?", 1)).Error; err != nil {
		return err
	}
	return nil
}

func IncProblemStatHelp(problemID int) error {
	if err := DB.Model(&ProblemStatistics{}).Where("problem_id = ?", problemID).UpdateColumn("HelpRequest", gorm.Expr("help_request + ?", 1)).Error; err != nil {
		return err
	}
	return nil
}

func IncProblemStatGradedCorrect(problemID int) error {
	if err := DB.Model(&ProblemStatistics{}).Where("problem_id = ?", problemID).UpdateColumn("GradedCorrect", gorm.Expr("graded_correct + ?", 1)).Error; err != nil {
		return err
	}
	return nil
}

func IncProblemStatGradedIncorrect(problemID int) error {
	if err := DB.Model(&ProblemStatistics{}).Where("problem_id = ?", problemID).UpdateColumn("GradedIncorrect", gorm.Expr("graded_incorrect + ?", 1)).Error; err != nil {
		return err
	}
	return nil
}

func AddSubmissionComplete(problemID, studentID int, studentCode string, submissionCategory, attemptNumber int, codeSubmittedAt time.Time, completed time.Time, snapshotID int, answer string) error {
	submission := Submission{
		ProblemID:          problemID,
		StudentID:          studentID,
		StudentCode:        studentCode,
		SubmissionCategory: submissionCategory,
		AttemptNumber:      attemptNumber,
		CodeSubmittedAt:    codeSubmittedAt,
		Completed:          completed,
		SnapshotID:         snapshotID,
		Answer:             answer,
	}
	if err := DB.Create(&submission).Error; err != nil {
		return err
	}
	return nil
}

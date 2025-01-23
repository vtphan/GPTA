package main

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"math"
	"time"
)

func AddStudent(name string, password string) (Student, error) {
	student := Student{
		Name:     name,
		Password: password,
	}
	if err := DB.Create(&student).Error; err != nil {
		return student, fmt.Errorf("failed to insert student: %w", err)
	}
	return student, nil
}

func AddTeacher(name string, password string) (Teacher, error) {
	teacher := Teacher{
		Name:     name,
		Password: password,
	}
	if err := DB.Create(&teacher).Error; err != nil {
		return teacher, fmt.Errorf("failed to insert student: %w", err)
	}
	return teacher, nil
}

func FindStudentByName(name string) (Student, error) {
	var student Student
	err := DB.Where("name = ?", name).First(&student).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return student, fmt.Errorf("failed to query student: %w", err)
	}
	if err == gorm.ErrRecordNotFound {
		return student, fmt.Errorf(RecordNotFound)
	}
	return student, nil
}

func FindTeacherByName(name string) (Teacher, error) {
	var teacher Teacher
	err := DB.Where("name = ?", name).First(&teacher).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return teacher, fmt.Errorf("failed to query student: %w", err)
	}
	if err == gorm.ErrRecordNotFound {
		return teacher, fmt.Errorf(RecordNotFound)
	}
	return teacher, nil
}

func AddProblem(teacherID int, problemDescription, answer, filename string, merit, effort, attempts, topicID, tag int) (Problem, error) {
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
		ProblemUploadedAt:  time.Now(),
	}
	if err := DB.Create(&problem).Error; err != nil {
		return problem, fmt.Errorf("failed to insert problem: %w", err)
	}
	return problem, nil
}

func AddSubmission(problemID, studentID int, studentCode string, submissionCategory, attemptNumber int, codeSubmittedAt time.Time, snapshotID int, answer string) (Submission, error) {
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
		return submission, fmt.Errorf("failed to insert submission: %w", err)
	}
	return submission, nil
}

func AddSubmissionComplete(problemID, studentID int, studentCode string, submissionCategory, attemptNumber int, codeSubmittedAt, completedAt time.Time, snapshotID int, answer string) (Submission, error) {
	submission := Submission{
		ProblemID:          problemID,
		StudentID:          studentID,
		StudentCode:        studentCode,
		SubmissionCategory: submissionCategory,
		AttemptNumber:      attemptNumber,
		CodeSubmittedAt:    codeSubmittedAt,
		Completed:          completedAt,
		SnapshotID:         snapshotID,
		Answer:             answer,
	}
	if err := DB.Create(&submission).Error; err != nil {
		return submission, fmt.Errorf("failed to insert submission: %w", err)
	}
	return submission, nil
}

func CompleteSubmission(completed time.Time, verdict string, submissionID int) error {
	if err := DB.Model(&Submission{}).
		Where("id = ?", submissionID).
		Updates(map[string]interface{}{
			"completed": completed,
			"verdict":   verdict,
		}).Error; err != nil {
		return fmt.Errorf("failed to update submission: %w", err)
	}

	return nil
}

func AddScore(problemID, studentID, teacherID, score, gradedSubmissionNumber int, scoreGivenAt time.Time) (Score, error) {
	scoreRecord := Score{
		ProblemID:              problemID,
		StudentID:              studentID,
		TeacherID:              teacherID,
		Score:                  score,
		GradedSubmissionNumber: gradedSubmissionNumber,
		ScoreGivenAt:           &scoreGivenAt,
	}
	if err := DB.Create(&scoreRecord).Error; err != nil {
		return scoreRecord, fmt.Errorf("failed to insert score: %w", err)
	}

	return scoreRecord, nil
}

func UpdateScore(scoreID, teacherID, score, gradedSubmissionNumber int) error {
	if err := DB.Model(&Score{}).
		Where("id = ?", scoreID).
		Updates(map[string]interface{}{
			"teacher_id":               teacherID,
			"score":                    score,
			"graded_submission_number": gradedSubmissionNumber,
		}).Error; err != nil {
		return fmt.Errorf("failed to update score: %w", err)
	}
	return nil
}

func AddFeedback(teacherID, studentID int, feedback string, feedbackGivenAt time.Time, submissionID int) (Feedback, error) {
	feedbackRecord := Feedback{
		TeacherID:       teacherID,
		StudentID:       studentID,
		Feedback:        feedback,
		FeedbackGivenAt: feedbackGivenAt,
		SubmissionID:    submissionID,
	}

	// Insert the feedback record into the database using GORM
	if err := DB.Create(&feedbackRecord).Error; err != nil {
		return feedbackRecord, fmt.Errorf("failed to insert feedback: %w", err)
	}

	return feedbackRecord, nil
}

func AddAttendance(studentID int, attendanceAt time.Time) (Attendance, error) {
	attendanceRecord := Attendance{
		StudentID:    studentID,
		AttendanceAt: attendanceAt,
	}
	if err := DB.Create(&attendanceRecord).Error; err != nil {
		return attendanceRecord, fmt.Errorf("failed to insert attendance: %w", err)
	}

	return attendanceRecord, nil
}

func AddTag(topicDescription string) (Tag, error) {
	tagRecord := Tag{
		TopicDescription: topicDescription,
	}
	if err := DB.Create(&tagRecord).Error; err != nil {
		return tagRecord, fmt.Errorf("failed to insert tag: %w", err)
	}

	return tagRecord, nil
}

func AddTestCase(problemID, studentID int, testCases string, addedAt time.Time) (TestCase, error) {
	testCaseRecord := TestCase{
		ProblemID: problemID,
		StudentID: studentID,
		TestCases: testCases,
		AddedAt:   addedAt,
	}
	if err := DB.Create(&testCaseRecord).Error; err != nil {
		return testCaseRecord, fmt.Errorf("failed to insert test case: %w", err)
	}
	return testCaseRecord, nil
}

func UpdateTestCase(testCases string, addedAt time.Time, testCaseID int) error {
	if err := DB.Model(&TestCase{}).
		Where("id = ?", testCaseID).
		Updates(map[string]interface{}{
			"test_cases": testCases,
			"added_at":   addedAt,
		}).Error; err != nil {
		return fmt.Errorf("failed to update test case with ID %d: %w", testCaseID, err)
	}

	return nil
}

func AddHelpMessage(codeExplanationID int, studentID int, message string, givenAt time.Time) (HelpMessage, error) {
	helpMessage := HelpMessage{
		CodeExplanationID: codeExplanationID,
		StudentID:         studentID,
		Message:           message,
		GivenAt:           givenAt,
	}
	if err := DB.Create(&helpMessage).Error; err != nil {
		return helpMessage, fmt.Errorf("failed to add help message: %w", err)
	}
	return helpMessage, nil
}

func UpdateHelpMessage(useful string, updatedAt time.Time, helpMessageID int) error {
	if err := DB.Model(&HelpMessage{}).
		Where("id = ?", helpMessageID).
		Updates(map[string]interface{}{
			"useful":     useful,
			"updated_at": updatedAt,
		}).Error; err != nil {
		return fmt.Errorf("failed to update help message with ID %d: %w", helpMessageID, err)
	}
	return nil
}

func AddCodeSnapshot(studentID int, problemID int, code string, status int, lastUpdatedAt time.Time, event string) (CodeSnapshot, error) {
	codeSnapshot := CodeSnapshot{
		StudentID:     studentID,
		ProblemID:     problemID,
		Code:          code,
		LastUpdatedAt: lastUpdatedAt,
		Status:        status,
		Event:         event,
	}
	if err := DB.Create(&codeSnapshot).Error; err != nil {
		return codeSnapshot, fmt.Errorf("failed to insert code snapshot: %w", err)
	}
	return codeSnapshot, nil
}

func UpdateProblemEndTime(problemEndedAt time.Time, problemID int) error {
	if err := DB.Model(&Problem{}).
		Where("id = ?", problemID).
		Update("problem_ended_at", problemEndedAt).Error; err != nil {
		return fmt.Errorf("failed to update problem end time for problem ID %d: %w", problemID, err)
	}
	return nil
}

func AddHelpEligible(problemID int, studentID int, becameEligibleAt time.Time) (HelpEligible, error) {
	helpEligible := HelpEligible{
		ProblemID:        problemID,
		StudentID:        studentID,
		BecameEligibleAt: becameEligibleAt,
	}
	if err := DB.Create(&helpEligible).Error; err != nil {
		return helpEligible, fmt.Errorf("failed to insert help eligible: %w", err)
	}
	return helpEligible, nil
}

func AddUserEventLog(name string, userID int, userType string, eventType string, referralInfo string, eventTime time.Time) (UserEventLog, error) {
	userEventLog := UserEventLog{
		Name:         name,
		UserID:       userID,
		UserType:     userType,
		EventType:    eventType,
		ReferralInfo: referralInfo,
		EventTime:    eventTime,
	}
	if err := DB.Create(&userEventLog).Error; err != nil {
		return userEventLog, fmt.Errorf("failed to insert user event log: %w", err)
	}

	return userEventLog, nil
}

func AddStudentStatus(studentID int, problemID int, codingStat string, helpStat string, submissionStat string, tutoringStat string, lastUpdatedAt time.Time) (StudentStatus, error) {
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
		return studentStatus, fmt.Errorf("failed to insert student status: %w", err)
	}
	return studentStatus, nil
}

func UpdateStudentCodingStat(codingStat string, lastUpdatedAt time.Time, studentID int, problemID int) error {
	studentStatus := StudentStatus{
		CodingStat:    codingStat,
		LastUpdatedAt: lastUpdatedAt,
	}
	if err := DB.Model(&StudentStatus{}).
		Where("student_id = ? AND problem_id = ?", studentID, problemID).
		Updates(studentStatus).Error; err != nil {
		return fmt.Errorf("failed to update student coding status: %w", err)
	}
	return nil
}

func UpdateStudentSubmissionStat(submissionStat string, lastUpdatedAt time.Time, studentID int, problemID int) error {
	studentStatus := StudentStatus{
		SubmissionStat: submissionStat,
		LastUpdatedAt:  lastUpdatedAt,
	}
	if err := DB.Model(&StudentStatus{}).
		Where("student_id = ? AND problem_id = ?", studentID, problemID).
		Updates(studentStatus).Error; err != nil {
		return fmt.Errorf("failed to update student submission status: %w", err)
	}
	return nil
}

func UpdateStudentHelpStat(helpStat string, lastUpdatedAt time.Time, studentID int, problemID int) error {
	studentStatus := StudentStatus{
		HelpStat:      helpStat,
		LastUpdatedAt: lastUpdatedAt,
	}
	if err := DB.Model(&StudentStatus{}).
		Where("student_id = ? AND problem_id = ?", studentID, problemID).
		Updates(studentStatus).Error; err != nil {
		return fmt.Errorf("failed to update student help status: %w", err)
	}
	return nil
}

func UpdateStudentTutoringStat(tutoringStat string, lastUpdatedAt time.Time, studentID int, problemID int) error {
	studentStatus := StudentStatus{
		TutoringStat:  tutoringStat,
		LastUpdatedAt: lastUpdatedAt,
	}
	if err := DB.Model(&StudentStatus{}).
		Where("student_id = ? AND problem_id = ?", studentID, problemID).
		Updates(studentStatus).Error; err != nil {
		return fmt.Errorf("failed to update student tutoring status: %w", err)
	}

	return nil
}

func AddMessage(snapshotID int, message string, authorID int, authorRole string, givenAt time.Time, messageType int) (int, error) {
	msg := Message{
		SnapshotID: snapshotID,
		Message:    message,
		AuthorID:   authorID,
		AuthorRole: authorRole,
		GivenAt:    givenAt,
		Type:       messageType,
	}
	if err := DB.Create(&msg).Error; err != nil {
		return 0, fmt.Errorf("failed to add message: %w", err)
	}
	return msg.ID, nil
}

func AddMessageFeedback(messageID int, feedback string, authorID int, authorRole string, givenAt time.Time) (int, error) {
	feedbackRecord := MessageFeedback{
		MessageID:  messageID,
		Feedback:   feedback,
		AuthorID:   authorID,
		AuthorRole: authorRole,
		GivenAt:    givenAt,
	}
	if err := DB.Create(&feedbackRecord).Error; err != nil {
		return 0, fmt.Errorf("failed to add message feedback: %w", err)
	}
	return feedbackRecord.ID, nil
}

func AddProblemStatistics(problemID int) error {
	statsRecord := ProblemStatistics{
		ProblemID:       problemID,
		Active:          0,
		Submission:      0,
		HelpRequest:     0,
		GradedCorrect:   0,
		GradedIncorrect: 0,
	}
	if err := DB.Create(&statsRecord).Error; err != nil {
		return fmt.Errorf("failed to add problem statistics: %w", err)
	}
	return nil
}

func IncrementProblemStatActive(problemID int) error {
	result := DB.Model(&ProblemStatistics{}).
		Where("problem_id = ?", problemID).
		UpdateColumn("active", gorm.Expr("active + ?", 1))
	if result.Error != nil {
		return fmt.Errorf("failed to increment active stat for problem %d: %w", problemID, result.Error)
	}
	return nil
}

func IncrementProblemStatSubmission(problemID int) error {
	result := DB.Model(&ProblemStatistics{}).
		Where("problem_id = ?", problemID).
		UpdateColumn("submission", gorm.Expr("submission + ?", 1))
	if result.Error != nil {
		return fmt.Errorf("failed to increment submission stat for problem %d: %w", problemID, result.Error)
	}
	return nil
}

func IncrementProblemStatHelp(problemID int) error {
	result := DB.Model(&ProblemStatistics{}).
		Where("problem_id = ?", problemID).
		UpdateColumn("help_request", gorm.Expr("help_request + ?", 1))
	if result.Error != nil {
		return fmt.Errorf("failed to increment help request stat for problem %d: %w", problemID, result.Error)
	}
	return nil
}

func IncrementProblemStatGradedCorrect(problemID int) error {
	result := DB.Model(&ProblemStatistics{}).
		Where("problem_id = ?", problemID).
		UpdateColumn("graded_correct", gorm.Expr("graded_correct + ?", 1))
	if result.Error != nil {
		return fmt.Errorf("failed to increment graded correct stat for problem %d: %w", problemID, result.Error)
	}
	return nil
}

func IncrementProblemStatGradedIncorrect(problemID int) error {
	result := DB.Model(&ProblemStatistics{}).
		Where("problem_id = ?", problemID).
		UpdateColumn("graded_incorrect", gorm.Expr("graded_incorrect + ?", 1))
	if result.Error != nil {
		return fmt.Errorf("failed to increment graded incorrect stat for problem %d: %w", problemID, result.Error)
	}
	return nil
}

func AddMessageBackFeedback(messageFeedbackID, authorID int, authorRole, useful string) error {
	messageBackFeedback := MessageBackFeedback{
		MessageFeedbackID: messageFeedbackID,
		AuthorID:          authorID,
		AuthorRole:        authorRole,
		Useful:            useful,
		GivenAt:           time.Now(),
	}
	result := DB.Create(&messageBackFeedback)
	if result.Error != nil {
		return fmt.Errorf("failed to add message back feedback: %w", result.Error)
	}
	return nil
}

func UpdateMessageBackFeedback(useful string, givenAt time.Time, feedbackID, authorID int, authorRole string) error {
	result := DB.Model(&MessageBackFeedback{}).
		Where("message_feedback_id = ? AND author_id = ? AND author_role = ?", feedbackID, authorID, authorRole).
		UpdateColumns(map[string]interface{}{
			"useful":   useful,
			"given_at": givenAt,
		})
	if result.Error != nil {
		return fmt.Errorf("failed to update message back feedback: %w", result.Error)
	}
	return nil
}

func GetSubmissions() ([]Submission, error) {
	var submissions []Submission
	if err := DB.Model(&Submission{}).Select("problem_id", "student_id", "code_submitted_at").Find(&submissions).Error; err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	return submissions, nil
}

func GetProblemUploadTime(pid int) (time.Time, error) {
	var problem struct {
		ProblemUploadedAt time.Time `gorm:"column:problem_uploaded_at"`
	}
	if err := DB.Model(&Problem{}).
		Where("id = ?", pid).
		Select("problem_uploaded_at").
		First(&problem).Error; err != nil {
		return time.Time{}, fmt.Errorf("failed to retrieve upload time for problem %d: %w", pid, err)
	}
	return problem.ProblemUploadedAt, nil
}

func GetSubmissionsByProblemID(pid int) ([]Submission, error) {
	var submissions []Submission
	if err := DB.Model(&Submission{}).
		Where("problem_id = ?", pid).
		Select("student_id", "submission_category", "code_submitted_at", "completed").
		Find(&submissions).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve submissions for problem %d: %w", pid, err)
	}
	return submissions, nil
}

func GetStudentName(studentID int) (string, error) {
	var student Student
	if err := DB.Model(&Student{}).
		Where("id = ?", studentID).
		Select("name").
		First(&student).Error; err != nil {
		return "", fmt.Errorf("failed to retrieve student name for student ID %d: %w", studentID, err)
	}
	return student.Name, nil
}

func GetCodeSnapshot(snapshotID int) (*CodeSnapshot, error) {
	var codeSnapshot CodeSnapshot
	if err := DB.Table("code_snapshots cs").
		Joins("join problems p on cs.problem_id = p.id").
		Where("cs.id = ?", snapshotID).
		Select("cs.student_id, cs.problem_id, cs.code, p.filename").
		First(&codeSnapshot).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve code snapshot for snapshot ID %d: %w", snapshotID, err)
	}
	return &codeSnapshot, nil
}

func GetCodeSnapshotMessageDetails(messageID int) ([]CodeSnapshotMessageDetails, error) {
	var details []CodeSnapshotMessageDetails
	if err := DB.Table("code_snapshots cs").
		Select("cs.student_id AS student_id, cs.problem_id AS problem_id, cs.code AS code, p.filename AS filename, m.type AS message_type").
		Joins("JOIN problems p ON cs.problem_id = p.id").
		Joins("JOIN messages m ON m.snapshot_id = cs.id").
		Where("m.id = ?", messageID).
		Find(&details).Error; err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	return details, nil
}

func GetVoteCount(feedbackID int, voteType string) (int64, error) {
	var count int64
	if err := DB.Model(&SnapshotBackFeedback{}).
		Where("is_helpful = ? AND snapshot_feedback_id = ?", voteType, feedbackID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}
	return count, nil
}

func GetNumberOfReply(snapshotID int) (int, error) {
	var count int64
	if err := DB.Model(&SnapshotFeedback{}).
		Where("snapshot_id = ?", snapshotID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to retrieve the count of replies for snapshot ID %d: %w", snapshotID, err)
	}
	return int(count), nil
}

func GetAllTeachers() ([]Teacher, error) {
	var teachers []Teacher
	if err := DB.Model(&Teacher{}).Find(&teachers).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve teachers: %w", err)
	}
	return teachers, nil
}

func GetScoreDetails(problemID, studentID int) (*Score, error) {
	var score Score
	if err := DB.Model(&Score{}).
		Select("id, score, graded_submission_number, teacher_id").
		Where("problem_id = ? AND student_id = ?", problemID, studentID).
		First(&score).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No record found
		}
		return nil, fmt.Errorf("failed to retrieve score details for problem ID %d and student ID %d: %w", problemID, studentID, err)
	}
	return &score, nil
}

func GetProblemDetails(problemID int) (*Problem, error) {
	var problem Problem
	if err := DB.Model(&Problem{}).
		Select("merit, effort").
		Where("id = ?", problemID).
		First(&problem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No record found
		}
		return nil, fmt.Errorf("failed to retrieve problem details for problem ID %d: %w", problemID, err)
	}
	return &problem, nil
}

func GetStudentStatus(studentID, problemID int) (*StudentStatus, error) {
	var status StudentStatus
	if err := DB.Model(&StudentStatus{}).
		Where("student_id = ? AND problem_id = ?", studentID, problemID).
		First(&status).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No record found
		}
		return nil, fmt.Errorf("failed to retrieve student status: %w", err)
	}
	return &status, nil
}

func GetStudentByID(studentID int) (*Student, error) {
	var student Student
	if err := DB.Model(&Student{}).
		Where("id = ?", studentID).
		First(&student).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No record found
		}
		return nil, fmt.Errorf("failed to retrieve student: %w", err)
	}
	return &student, nil
}

func GetProblemIDByFilename(filename string) (int, error) {
	var problem Problem
	if err := DB.Model(&Problem{}).
		Select("id").
		Where("filename = ?", filename).
		First(&problem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil // No record found
		}
		return 0, fmt.Errorf("failed to retrieve problem ID: %w", err)
	}
	return problem.ID, nil
}

func GetTestCasesByProblemID(problemID int) ([]string, error) {
	var testCases []string
	if err := DB.Model(&TestCase{}).
		Select("test_cases").
		Where("problem_id = ?", problemID).
		Pluck("test_cases", &testCases).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve test cases: %w", err)
	}
	return testCases, nil
}

func GetCurrentStudents() []int {
	var studentIDs []int
	date := time.Now().Format("2006-01-02") // Correct format for GORM's DATE query
	if err := DB.Model(&Attendance{}).
		Select("student_id").
		Where("DATE(attendance_at) = ?", date).
		Pluck("student_id", &studentIDs).Error; err != nil {
		return nil
	}
	return studentIDs
}

func GetAllStudents() map[int]string {
	var students []Student
	studentMap := make(map[int]string)

	if err := DB.Model(&Student{}).
		Select("id, name").
		Find(&students).Error; err != nil {
		return nil
	}

	for _, student := range students {
		studentMap[student.ID] = student.Name
	}

	return studentMap
}

func GetProblemStats(problemID int) (int, int, int, int, int) {
	var stats ProblemStatistics

	if err := DB.Where("problem_id = ?", problemID).First(&stats).Error; err != nil {
		return 0, 0, 0, 0, 0
	}

	ungraded := stats.Submission - stats.GradedCorrect - stats.GradedIncorrect

	return stats.Active, stats.HelpRequest, ungraded, stats.GradedCorrect, stats.GradedIncorrect
}

func GetLatestSubmissionTime(problemID int) map[int]time.Time {
	var latestSubmissions = make(map[int]time.Time)
	var submissions []Submission

	if err := DB.Model(&Submission{}).
		Select("student_id, max(code_submitted_at) as code_submitted_at").
		Where("problem_id = ?", problemID).
		Group("student_id").
		Find(&submissions).Error; err != nil {
		log.Fatal(err)
	}

	for _, submission := range submissions {
		latestSubmissions[submission.StudentID] = submission.CodeSubmittedAt
	}

	return latestSubmissions
}

func GetProblemNameFromID(problemID int) string {
	var problem Problem

	if err := DB.Model(&Problem{}).
		Select("filename").
		Where("id = ?", problemID).
		First(&problem).Error; err != nil {
		return ""
	}

	return problem.Filename
}

func GetCodeSnapshotsByProblemID(problemID int) ([]CodeSnapshot, error) {
	var codeSnapshots []CodeSnapshot
	err := DB.Model(&CodeSnapshot{}).
		Select("student_id, MAX(last_updated_at) as last_updated_at").
		Where("problem_id = ?", problemID).
		Group("student_id").
		Find(&codeSnapshots).Error
	if err != nil {
		log.Fatalf("Failed to fetch code snapshots: %v", err)
	}
	return codeSnapshots, err
}

func GetProblemDetail(problemID int) (string, time.Time, error) {
	var problem struct {
		Description    string
		ProblemEndedAt time.Time
	}
	err := DB.Table("problem").
		Select("problem_description as description, problem_ended_at").
		Where("id = ?", problemID).
		Scan(&problem).Error
	return problem.Description, problem.ProblemEndedAt, err
}

func GetStudentStatusesByProblemID(problemID int) ([]StudentStatus, error) {
	var statuses []StudentStatus
	err := DB.Where("problem_id = ?", problemID).Find(&statuses).Error
	return statuses, err
}

func GetAnswerStats(problemID int) ([]*AnswerStatInfo, error) {
	var stats []struct {
		Answer string
		Count  int
	}
	err := DB.Table("submission").
		Select("answer, COUNT(*) as count").
		Where("problem_id = ? AND answer IS NOT NULL AND LENGTH(answer) > 0", problemID).
		Group("answer").
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	var answerStats []*AnswerStatInfo
	var total int
	for _, stat := range stats {
		total += stat.Count
	}
	for _, stat := range stats {
		percent := float64(stat.Count) * 100.0 / float64(total)
		percent = math.Round(percent*100) / 100
		answerStats = append(answerStats, &AnswerStatInfo{
			Answer:  stat.Answer,
			Count:   stat.Count,
			Percent: percent,
		})
	}

	return answerStats, nil
}

// GetProblems fetches all problems from the database with the specified fields.
func GetProblems() ([]Problem, error) {
	var problems []Problem

	// Query the database to fetch the required fields for all problems
	if err := DB.Select("id, filename, problem_uploaded_at, problem_ended_at").Find(&problems).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch problems: %w", err)
	}

	return problems, nil
}

func GetTeacherByName(name string) (Teacher, error) {
	var teacher Teacher
	if err := DB.Where("name = ?", name).First(&teacher).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return teacher, fmt.Errorf("teacher not found")
		}
		return teacher, fmt.Errorf("failed to fetch teacher: %w", err)
	}
	return teacher, nil
}

func GetStudentByName(name string) (Student, error) {
	var student Student
	if err := DB.Where("name = ?", name).First(&student).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return student, fmt.Errorf("student not found")
		}
		return student, fmt.Errorf("failed to fetch student: %w", err)
	}
	return student, nil
}

func GetTags() (map[int]string, error) {
	var tags []Tag
	if err := DB.Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve tags: %w", err)
	}

	tagMap := make(map[int]string)
	for _, tag := range tags {
		tagMap[tag.ID] = tag.TopicDescription
	}

	return tagMap, nil
}

func GetStudentScores() (map[int]*ScoreEntry, error) {
	var scores []struct {
		Score                  int
		GradedSubmissionNumber int
		StudentID              int
		StudentName            string
	}

	if err := DB.Table("scores").
		Joins("join students on score.student_id = student.id").
		Select("scores.score, scores.graded_submission_number, scores.student_id, students.name").
		Find(&scores).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve scores: %w", err)
	}

	scoreEntries := make(map[int]*ScoreEntry)
	for _, score := range scores {
		if _, ok := scoreEntries[score.StudentID]; !ok {
			scoreEntries[score.StudentID] = &ScoreEntry{Name: score.StudentName}
		}
		scoreEntries[score.StudentID].Points += score.Score
		scoreEntries[score.StudentID].Attempts += score.GradedSubmissionNumber
		scoreEntries[score.StudentID].Count++
	}

	return scoreEntries, nil
}

func GetTagDescriptionByID(tagID int) (string, error) {
	var tag Tag
	if err := DB.Where("id = ?", tagID).First(&tag).Error; err != nil {
		return "", fmt.Errorf("failed to retrieve tag description: %w", err)
	}
	return tag.TopicDescription, nil
}

func GetProblemPerformanceByTagID(tagID int) (map[int]*ProblemPerformance, error) {
	var results []struct {
		Pid       int
		Merit     int
		At        time.Time
		Points    int
		StudentID int
	}

	if err := DB.Table("problems").
		Joins("join scores on problems.id = scores.problem_id").
		Joins("join students on scores.student_id = students.id").
		Where("problems.tag = ?", tagID).
		Select("problems.id, problems.merit, problems.at, scores.points, scores.student_id").
		Find(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve problem performance: %w", err)
	}

	record := make(map[int]*ProblemPerformance)
	for _, result := range results {
		if _, ok := record[result.Pid]; !ok {
			record[result.Pid] = &ProblemPerformance{
				Pid:       result.Pid,
				Timestamp: result.At.UnixNano(),
				Correct:   0,
				Incorrect: 0,
				Activity:  0,
				PC:        Passcode,
			}
		}
		if result.Merit == result.Points {
			record[result.Pid].Correct++
		} else {
			record[result.Pid].Incorrect++
		}
		record[result.Pid].Activity += 1.0
	}

	return record, nil
}

func GetStudentCount() (float32, error) {
	var count int64
	if err := DB.Table("student").Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to get student count: %w", err)
	}
	return float32(count), nil
}

func GetLatestProblemID() (int, error) {
	var problem Problem
	// Retrieve the latest problem based on the highest ID
	if err := DB.Order("id desc").First(&problem).Error; err != nil {
		return 0, fmt.Errorf("failed to retrieve the latest problem ID: %w", err)
	}
	return problem.ID, nil
}

func GetProblemStatistics(pid int) (map[int]int, string, time.Time, map[string]int, map[string][]float64, error) {
	var results []struct {
		StudentID          int
		Score              int
		Attempts           int
		ProblemUploadedAt  time.Time
		ProblemDescription string
		SubmissionID       int
		CodeSubmittedAt    time.Time
		Completed          time.Time
	}

	participants := make(map[int]int)
	performance := make(map[string]int)
	durations := make(map[string][]float64)
	var probContent string
	var probAt time.Time

	err := DB.Table("scores").
		Joins("join problems on scores.problem_id = problems.id").
		Joins("join submissions on scores.problem_id = submissions.problem_id and scores.student_id = submissions.student_id").
		Where("problems.id = ?", pid).
		Order("submissions.id desc").
		Select(`scores.student_id, scores.score, scores.graded_submission_number, 
                problems.problem_uploaded_at, problems.problem_description, 
                submissions.id as submission_id, submissions.code_submitted_at, 
                submissions.completed`).
		Scan(&results).Error
	if err != nil {
		return nil, "", time.Time{}, nil, nil, fmt.Errorf("failed to retrieve problem statistics: %w", err)
	}

	for _, result := range results {
		if _, ok := participants[result.StudentID]; !ok {
			participants[result.StudentID] = result.SubmissionID
			probAt = result.ProblemUploadedAt
			probContent = result.ProblemDescription

			probDuration := result.CodeSubmittedAt.Sub(probAt).Minutes()
			key := fmt.Sprintf("%d points", result.Score)
			performance[key]++
			if _, ok := durations[key]; !ok {
				durations[key] = make([]float64, 0)
			}
			durations[key] = append(durations[key], probDuration)
		}
	}

	return participants, probContent, probAt, performance, durations, nil
}

func GetAttendanceByDate(theDate string) (map[int]int, error) {
	var attendances []Attendance
	attendants := make(map[int]int)

	err := DB.Where("DATE(attendance_at) = ?", theDate).
		Find(&attendances).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve attendance: %w", err)
	}

	for _, att := range attendances {
		attendants[att.StudentID] = 0
	}

	return attendants, nil
}

func GetAttendanceByStudentID(uid int) ([]Attendance, error) {
	var attendances []Attendance
	err := DB.Where("student_id = ?", uid).
		Select("attendance_at").
		Find(&attendances).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve attendance for student ID %d: %w", uid, err)
	}
	return attendances, nil
}

func GetCurrentUserVote(feedbackID int, userID int, userRole string) string {
	var feedback MessageBackFeedback

	err := DB.Where("message_feedback_id = ? AND author_id = ? AND author_role = ?", feedbackID, userID, userRole).
		First(&feedback).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ""
		}
		return ""
	}

	return feedback.Useful
}

func GetBackFeedbackCount(feedbackID int, backFeedbackType string) int {
	var count int64

	err := DB.Model(&MessageBackFeedback{}).
		Where("useful = ? AND message_feedback_id = ?", backFeedbackType, feedbackID).
		Count(&count).Error
	if err != nil {
		log.Fatal(err)
	}

	return int(count)
}

func GetMessageFeedbacksByMessageID(messageID int) ([]MessageFeedback, error) {
	var messageFeedbacks []MessageFeedback
	err := DB.Where("message_id = ?", messageID).Find(&messageFeedbacks).Error
	if err != nil {
		return nil, fmt.Errorf("Error retrieving message feedbacks: %v", err)
	}
	return messageFeedbacks, nil
}

func GetLatestSnapshot(studentID int, problemID int) (*Snapshot, error) {
	var snapshot CodeSnapshot

	err := DB.Where("student_id = ? AND problem_id = ?", studentID, problemID).
		Order("last_updated_at DESC").
		First(&snapshot).Error
	if err != nil {
		return nil, fmt.Errorf("Error retrieving latest snapshot for student %d, problem %d: %v", studentID, problemID, err)
	}

	return &Snapshot{
		ID:          snapshot.ID,
		StudentID:   snapshot.StudentID,
		ProblemID:   snapshot.ProblemID,
		ProblemName: GetProblemNameFromID(problemID),
		Code:        snapshot.Code,
		LastUpdated: snapshot.LastUpdatedAt,
	}, nil
}

func GetTeacherName(authorID int) string {
	var teacher Teacher

	err := DB.Where("id = ?", authorID).First(&teacher).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ""
		}
		log.Printf("Error retrieving teacher with ID %d: %v", authorID, err)
		return ""
	}

	return teacher.Name
}

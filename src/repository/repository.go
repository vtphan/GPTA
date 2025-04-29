package repository

import (
	"errors"
	"fmt"

	"github.com/GPTA/src/models"
	"gorm.io/gorm"
	"log"
	"math"
	"sort"
	"time"
)

func AddStudent(name string, password string) (models.Student, error) {
	student := models.Student{
		Name:     name,
		Password: password,
	}
	if err := models.DB.Create(&student).Error; err != nil {
		return student, fmt.Errorf("failed to insert student: %w", err)
	}
	return student, nil
}

func AddTeacher(name string, password string) (models.Teacher, error) {
	teacher := models.Teacher{
		Name:     name,
		Password: password,
	}
	if err := models.DB.Create(&teacher).Error; err != nil {
		return teacher, fmt.Errorf("failed to insert student: %w", err)
	}
	return teacher, nil
}

func FindStudentByName(name string) (models.Student, error) {
	var student models.Student
	err := models.DB.Where("name = ?", name).First(&student).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return student, fmt.Errorf("failed to query student: %w", err)
	}
	if err == gorm.ErrRecordNotFound {
		return student, fmt.Errorf(models.RecordNotFound)
	}
	return student, nil
}

func FindTeacherByName(name string) (models.Teacher, error) {
	var teacher models.Teacher
	err := models.DB.Where("name = ?", name).First(&teacher).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return teacher, fmt.Errorf("failed to query student: %w", err)
	}
	if err == gorm.ErrRecordNotFound {
		return teacher, fmt.Errorf(models.RecordNotFound)
	}
	return teacher, nil
}

func AddProblem(teacherID int, problemDescription, answer, filename string, merit, effort, attempts, topicID, tag int) (models.Problem, error) {
	problem := models.Problem{
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
		CourseID:           models.CourseId,
	}
	if err := models.DB.Create(&problem).Error; err != nil {
		return problem, fmt.Errorf("failed to insert problem: %w", err)
	}
	return problem, nil
}

func AddSubmission(problemID, studentID int, studentCode string, submissionCategory, attemptNumber int, codeSubmittedAt time.Time, snapshotID int, answer string) (models.SubmissionTable, error) {
	submission := models.SubmissionTable{
		ProblemID:          problemID,
		StudentID:          studentID,
		StudentCode:        studentCode,
		SubmissionCategory: submissionCategory,
		AttemptNumber:      attemptNumber,
		CodeSubmittedAt:    codeSubmittedAt,
		SnapshotID:         snapshotID,
		Answer:             answer,
	}
	if err := models.DB.Create(&submission).Error; err != nil {
		return submission, fmt.Errorf("failed to insert submission: %w", err)
	}
	return submission, nil
}

func AddSubmissionComplete(problemID, studentID int, studentCode string, submissionCategory, attemptNumber int, codeSubmittedAt, completedAt time.Time, snapshotID int, answer string) (models.SubmissionTable, error) {
	submission := models.SubmissionTable{
		ProblemID:          problemID,
		StudentID:          studentID,
		StudentCode:        studentCode,
		SubmissionCategory: submissionCategory,
		AttemptNumber:      attemptNumber,
		CodeSubmittedAt:    codeSubmittedAt,
		Completed:          &completedAt,
		SnapshotID:         snapshotID,
		Answer:             answer,
	}
	if err := models.DB.Create(&submission).Error; err != nil {
		return submission, fmt.Errorf("failed to insert submission: %w", err)
	}
	return submission, nil
}

func CompleteSubmission(completed time.Time, verdict string, submissionID int) error {
	if err := models.DB.Model(&models.SubmissionTable{}).
		Where("id = ?", submissionID).
		Updates(map[string]interface{}{
			"completed": completed,
			"verdict":   verdict,
		}).Error; err != nil {
		return fmt.Errorf("failed to update submission: %w", err)
	}

	return nil
}

func AddScore(problemID, studentID, teacherID, score, gradedSubmissionNumber int, scoreGivenAt time.Time) (models.Score, error) {
	scoreRecord := models.Score{
		ProblemID:              problemID,
		StudentID:              studentID,
		TeacherID:              teacherID,
		Score:                  score,
		GradedSubmissionNumber: gradedSubmissionNumber,
		ScoreGivenAt:           &scoreGivenAt,
	}
	if err := models.DB.Create(&scoreRecord).Error; err != nil {
		return scoreRecord, fmt.Errorf("failed to insert score: %w", err)
	}

	return scoreRecord, nil
}

func UpdateScore(scoreID, teacherID, score, gradedSubmissionNumber int) error {
	if err := models.DB.Model(&models.Score{}).
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

func AddFeedback(teacherID, studentID int, feedback string, feedbackGivenAt time.Time, submissionID int) (models.Feedback, error) {
	feedbackRecord := models.Feedback{
		TeacherID:       teacherID,
		StudentID:       studentID,
		Feedback:        feedback,
		FeedbackGivenAt: feedbackGivenAt,
		SubmissionID:    submissionID,
	}

	// Insert the feedback record into the database using GORM
	if err := models.DB.Create(&feedbackRecord).Error; err != nil {
		return feedbackRecord, fmt.Errorf("failed to insert feedback: %w", err)
	}

	return feedbackRecord, nil
}

func AddAttendance(studentID int, attendanceAt time.Time) (models.Attendance, error) {
	attendanceRecord := models.Attendance{
		StudentID:    studentID,
		AttendanceAt: attendanceAt,
	}
	if err := models.DB.Create(&attendanceRecord).Error; err != nil {
		return attendanceRecord, fmt.Errorf("failed to insert attendance: %w", err)
	}

	return attendanceRecord, nil
}

func AddTag(topicDescription string) (models.Tag, error) {
	tagRecord := models.Tag{
		TopicDescription: topicDescription,
	}
	if err := models.DB.Create(&tagRecord).Error; err != nil {
		return tagRecord, fmt.Errorf("failed to insert tag: %w", err)
	}

	return tagRecord, nil
}

func AddTestCase(problemID, studentID int, testCases string, addedAt time.Time) (models.TestCase, error) {
	testCaseRecord := models.TestCase{
		ProblemID: problemID,
		StudentID: studentID,
		TestCases: testCases,
		AddedAt:   addedAt,
	}
	if err := models.DB.Create(&testCaseRecord).Error; err != nil {
		return testCaseRecord, fmt.Errorf("failed to insert test case: %w", err)
	}
	return testCaseRecord, nil
}

func UpdateTestCase(testCases string, addedAt time.Time, testCaseID int) error {
	if err := models.DB.Model(&models.TestCase{}).
		Where("id = ?", testCaseID).
		Updates(map[string]interface{}{
			"test_cases": testCases,
			"added_at":   addedAt,
		}).Error; err != nil {
		return fmt.Errorf("failed to update test case with ID %d: %w", testCaseID, err)
	}

	return nil
}

func AddHelpMessage(codeExplanationID int, studentID int, message string, givenAt time.Time) (models.HelpMessage, error) {
	helpMessage := models.HelpMessage{
		CodeExplanationID: codeExplanationID,
		StudentID:         studentID,
		Message:           message,
		GivenAt:           givenAt,
	}
	if err := models.DB.Create(&helpMessage).Error; err != nil {
		return helpMessage, fmt.Errorf("failed to add help message: %w", err)
	}
	return helpMessage, nil
}

func UpdateHelpMessage(useful string, updatedAt time.Time, helpMessageID int) error {
	if err := models.DB.Model(&models.HelpMessage{}).
		Where("id = ?", helpMessageID).
		Updates(map[string]interface{}{
			"useful":     useful,
			"updated_at": updatedAt,
		}).Error; err != nil {
		return fmt.Errorf("failed to update help message with ID %d: %w", helpMessageID, err)
	}
	return nil
}

func AddCodeSnapshot(studentID int, problemID int, code string, status int, lastUpdatedAt time.Time, event string) (models.CodeSnapshot, error) {
	codeSnapshot := models.CodeSnapshot{
		StudentID:     studentID,
		ProblemID:     problemID,
		Code:          code,
		LastUpdatedAt: lastUpdatedAt,
		Status:        status,
		Event:         event,
	}
	if err := models.DB.Create(&codeSnapshot).Error; err != nil {
		return codeSnapshot, fmt.Errorf("failed to insert code snapshot: %w", err)
	}
	return codeSnapshot, nil
}

func UpdateProblemEndTime(problemEndedAt time.Time, problemID int) error {
	if err := models.DB.Model(&models.Problem{}).
		Where("id = ?", problemID).
		Update("problem_ended_at", problemEndedAt).Error; err != nil {
		return fmt.Errorf("failed to update problem end time for problem ID %d: %w", problemID, err)
	}
	return nil
}

func AddHelpEligible(problemID int, studentID int, becameEligibleAt time.Time) (models.HelpEligible, error) {
	helpEligible := models.HelpEligible{
		ProblemID:        problemID,
		StudentID:        studentID,
		BecameEligibleAt: becameEligibleAt,
	}
	if err := models.DB.Create(&helpEligible).Error; err != nil {
		return helpEligible, fmt.Errorf("failed to insert help eligible: %w", err)
	}
	return helpEligible, nil
}

func AddUserEventLog(name string, userID int, userType string, eventType string, referralInfo string, eventTime time.Time) (models.UserEventLog, error) {
	userEventLog := models.UserEventLog{
		Name:         name,
		UserID:       userID,
		UserType:     userType,
		EventType:    eventType,
		ReferralInfo: referralInfo,
		EventTime:    eventTime,
	}
	if err := models.DB.Create(&userEventLog).Error; err != nil {
		return userEventLog, fmt.Errorf("failed to insert user event log: %w", err)
	}

	return userEventLog, nil
}

func AddStudentStatus(studentID int, problemID int, codingStat string, helpStat string, submissionStat string, lastUpdatedAt time.Time) (models.StudentStatus, error) {
	studentStatus := models.StudentStatus{
		StudentID:      studentID,
		ProblemID:      problemID,
		CodingStat:     codingStat,
		HelpStat:       helpStat,
		SubmissionStat: submissionStat,
		LastUpdatedAt:  lastUpdatedAt,
	}
	if err := models.DB.Create(&studentStatus).Error; err != nil {
		return studentStatus, fmt.Errorf("failed to insert student status: %w", err)
	}
	return studentStatus, nil
}

func UpdateStudentCodingStat(codingStat string, lastUpdatedAt time.Time, studentID int, problemID int) error {
	studentStatus := models.StudentStatus{
		CodingStat:    codingStat,
		LastUpdatedAt: lastUpdatedAt,
	}
	if err := models.DB.Model(&models.StudentStatus{}).
		Where("student_id = ? AND problem_id = ?", studentID, problemID).
		Updates(studentStatus).Error; err != nil {
		return fmt.Errorf("failed to update student coding status: %w", err)
	}
	return nil
}

func UpdateStudentPercentStat(percent int, explanation string, lastUpdatedAt time.Time, studentID int, problemID int) error {
	studentStatus := models.StudentStatus{
		Percentage:    percent,
		LastUpdatedAt: lastUpdatedAt,
		Explanation:   explanation,
	}
	if err := models.DB.Model(&models.StudentStatus{}).
		Where("student_id = ? AND problem_id = ?", studentID, problemID).
		Updates(studentStatus).Error; err != nil {
		return fmt.Errorf("failed to update student coding status: %w", err)
	}
	return nil
}

func UpdateStudentSubmissionStat(submissionStat string, lastUpdatedAt time.Time, studentID int, problemID int) error {
	studentStatus := models.StudentStatus{
		SubmissionStat: submissionStat,
		LastUpdatedAt:  lastUpdatedAt,
	}
	if err := models.DB.Model(&models.StudentStatus{}).
		Where("student_id = ? AND problem_id = ?", studentID, problemID).
		Updates(studentStatus).Error; err != nil {
		return fmt.Errorf("failed to update student submission status: %w", err)
	}
	return nil
}

func UpdateStudentHelpStat(helpStat string, lastUpdatedAt time.Time, studentID int, problemID int) error {
	studentStatus := models.StudentStatus{
		HelpStat:      helpStat,
		LastUpdatedAt: lastUpdatedAt,
	}
	if err := models.DB.Model(&models.StudentStatus{}).
		Where("student_id = ? AND problem_id = ?", studentID, problemID).
		Updates(studentStatus).Error; err != nil {
		return fmt.Errorf("failed to update student help status: %w", err)
	}
	return nil
}

func AddMessage(snapshotID int, message string, authorID int, authorRole string, givenAt time.Time, messageType int) (int, error) {
	msg := models.Message{
		SnapshotID: snapshotID,
		Message:    message,
		AuthorID:   authorID,
		AuthorRole: authorRole,
		GivenAt:    givenAt,
		Type:       messageType,
	}
	if err := models.DB.Create(&msg).Error; err != nil {
		return 0, fmt.Errorf("failed to add message: %w", err)
	}
	return msg.ID, nil
}

func AddMessageFeedback(messageID int, feedback string, authorID int, authorRole string, givenAt time.Time) (int, error) {
	feedbackRecord := models.MessageFeedback{
		MessageID:  messageID,
		Feedback:   feedback,
		AuthorID:   authorID,
		AuthorRole: authorRole,
		GivenAt:    givenAt,
	}
	if err := models.DB.Create(&feedbackRecord).Error; err != nil {
		return 0, fmt.Errorf("failed to add message feedback: %w", err)
	}
	return feedbackRecord.ID, nil
}

func AddProblemStatistics(problemID int) error {
	statsRecord := models.ProblemStatistics{
		ProblemID:       problemID,
		Active:          0,
		Submission:      0,
		HelpRequest:     0,
		GradedCorrect:   0,
		GradedIncorrect: 0,
	}
	if err := models.DB.Create(&statsRecord).Error; err != nil {
		return fmt.Errorf("failed to add problem statistics: %w", err)
	}
	return nil
}

func IncrementProblemStatActive(problemID int) error {
	result := models.DB.Model(&models.ProblemStatistics{}).
		Where("problem_id = ?", problemID).
		UpdateColumn("active", gorm.Expr("active + ?", 1))
	if result.Error != nil {
		return fmt.Errorf("failed to increment active stat for problem %d: %w", problemID, result.Error)
	}
	return nil
}

func IncrementProblemStatSubmission(problemID int) error {
	result := models.DB.Model(&models.ProblemStatistics{}).
		Where("problem_id = ?", problemID).
		UpdateColumn("submission", gorm.Expr("submission + ?", 1))
	if result.Error != nil {
		return fmt.Errorf("failed to increment submission stat for problem %d: %w", problemID, result.Error)
	}
	return nil
}

func IncrementProblemStatHelp(problemID int) error {
	result := models.DB.Model(&models.ProblemStatistics{}).
		Where("problem_id = ?", problemID).
		UpdateColumn("help_request", gorm.Expr("help_request + ?", 1))
	if result.Error != nil {
		return fmt.Errorf("failed to increment help request stat for problem %d: %w", problemID, result.Error)
	}
	return nil
}

func IncrementProblemStatGradedCorrect(problemID int) error {
	result := models.DB.Model(&models.ProblemStatistics{}).
		Where("problem_id = ?", problemID).
		UpdateColumn("graded_correct", gorm.Expr("graded_correct + ?", 1))
	if result.Error != nil {
		return fmt.Errorf("failed to increment graded correct stat for problem %d: %w", problemID, result.Error)
	}
	return nil
}

func IncrementProblemStatGradedIncorrect(problemID int) error {
	result := models.DB.Model(&models.ProblemStatistics{}).
		Where("problem_id = ?", problemID).
		UpdateColumn("graded_incorrect", gorm.Expr("graded_incorrect + ?", 1))
	if result.Error != nil {
		return fmt.Errorf("failed to increment graded incorrect stat for problem %d: %w", problemID, result.Error)
	}
	return nil
}

func AddMessageBackFeedback(messageFeedbackID, authorID int, authorRole, useful string) error {
	messageBackFeedback := models.MessageBackFeedback{
		MessageFeedbackID: messageFeedbackID,
		AuthorID:          authorID,
		AuthorRole:        authorRole,
		Useful:            useful,
		GivenAt:           time.Now(),
	}
	result := models.DB.Create(&messageBackFeedback)
	if result.Error != nil {
		return fmt.Errorf("failed to add message back feedback: %w", result.Error)
	}
	return nil
}

func UpdateMessageBackFeedback(useful string, givenAt time.Time, feedbackID, authorID int, authorRole string) error {
	result := models.DB.Model(&models.MessageBackFeedback{}).
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

func GetSubmissions() ([]models.SubmissionTable, error) {
	var submissions []models.SubmissionTable
	if err := models.DB.Model(&models.SubmissionTable{}).Select("problem_id", "student_id", "code_submitted_at").Find(&submissions).Error; err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	return submissions, nil
}

func GetProblemUploadTime(pid int) (time.Time, error) {
	var problem struct {
		ProblemUploadedAt time.Time `gorm:"column:problem_uploaded_at"`
	}
	if err := models.DB.Model(&models.Problem{}).
		Where("id = ?", pid).
		Select("problem_uploaded_at").
		First(&problem).Error; err != nil {
		return time.Time{}, fmt.Errorf("failed to retrieve upload time for problem %d: %w", pid, err)
	}
	return problem.ProblemUploadedAt, nil
}

func GetProblemDescription(problemID int) (string, error) {
	var problemDescription string
	// Query the database to get the problem description
	err := models.DB.Model(&models.Problem{}).
		Where("id = ?", problemID).
		Select("problem_description").
		First(&problemDescription).Error

	// Return the result or an error if something went wrong
	if err != nil {
		return "", fmt.Errorf("failed to fetch problem description: %w", err)
	}

	return problemDescription, nil
}

func GetLatestFeedbackForStudentAndProblem(studentID, problemID int) (string, error) {
	var feedback string

	err := models.DB.
		Table("message_feedbacks MF").
		Select("MF.feedback").
		Joins("JOIN messages M ON MF.message_id = M.id").
		Joins("JOIN code_snapshots C ON M.snapshot_id = C.id").
		Where("C.student_id = ? AND C.problem_id = ?", studentID, problemID).
		Order("MF.given_at DESC").
		Limit(1).
		Scan(&feedback).Error

	if err != nil {
		return "", fmt.Errorf("failed to fetch latest feedback: %v", err)
	}

	return feedback, nil
}

// GetClassFeedbackByProblemID retrieves the class feedback for a specific problem
func GetClassFeedbackByProblemID(problemID int) ([]models.ClassFeedback, error) {
	var feedbacks []models.ClassFeedback

	// Find all feedback records for the given problem ID
	if err := models.DB.Where("problem_id = ?", problemID).Find(&feedbacks).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve class feedback: %w", err)
	}

	return feedbacks, nil
}

func GetLatestClassFeedbackByProblemID(problemID int) (*models.ClassFeedback, error) {
	var feedback models.ClassFeedback

	// Find the latest feedback record for the given problem ID, ordered by FeedbackTime descending
	if err := models.DB.Where("problem_id = ?", problemID).
		Order("feedback_time desc").
		First(&feedback).Error; err != nil {
		// If no feedback is found, return nil, nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		// If another error occurs, return the error
		return nil, fmt.Errorf("failed to retrieve latest class feedback: %w", err)
	}

	return &feedback, nil
}

func GetLatestScaffold(studentID int, problemID int) (string, error) {
	var feedback models.ScaffoldingFeedback

	err := models.DB.
		Where("student_id = ? AND problem_id = ?", studentID, problemID).
		Order("feedback_time DESC").
		First(&feedback).Error

	if err != nil {
		return "", fmt.Errorf("failed to fetch latest scaffolding feedback: %w", err)
	}

	return feedback.Scaffold, nil
}

func GetClassFeedbackByFeedbackID(feedbackID int) ([]models.ClassFeedback, error) {
	var feedbacks []models.ClassFeedback

	// Find all feedback records for the given problem ID
	if err := models.DB.Where("id = ?", feedbackID).Find(&feedbacks).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve class feedback: %w", err)
	}

	return feedbacks, nil
}

func SaveStudentProgress(progress models.StudentProgress) error {

	// Create a new StudentProgress record
	studentProgress := models.StudentProgress{
		ProblemID:    progress.ProblemID,
		StudentID:    progress.StudentID,
		Explanation:  progress.Explanation,
		Percentage:   progress.Percentage,
		FeedbackTime: time.Now(), // Automatically capture the timestamp
	}

	// Save the student progress to the database
	if err := models.DB.Create(&studentProgress).Error; err != nil {
		return fmt.Errorf("failed to save student progress: %w", err)
	}

	return nil
}

func SaveClassFeedback(problemID int, feedback string) error {
	// Check if the problem exists
	var problem models.Problem
	if err := models.DB.First(&problem, problemID).Error; err != nil {
		return fmt.Errorf("problem not found: %w", err)
	}

	// Create a new ClassFeedback record
	classFeedback := models.ClassFeedback{
		ProblemID:    problemID,
		Feedback:     feedback,
		FeedbackTime: time.Now(), // Automatically capture the timestamp
	}

	// Save the feedback to the database
	if err := models.DB.Create(&classFeedback).Error; err != nil {
		return fmt.Errorf("failed to save class feedback: %w", err)
	}

	return nil
}

func SaveScaffoldingFeedback(studentID int, problemID int, scaffold string) error {

	// Create a new ScaffoldingFeedback record
	feedback := models.ScaffoldingFeedback{
		StudentID:    studentID,
		ProblemID:    problemID,
		Scaffold:     scaffold,
		FeedbackTime: time.Now(),
	}

	// Save to the database
	if err := models.DB.Create(&feedback).Error; err != nil {
		return fmt.Errorf("failed to save scaffolding feedback: %w", err)
	}

	return nil
}

func GetSubmissionsByProblemID(pid int) ([]models.SubmissionTable, error) {
	var submissions []models.SubmissionTable
	if err := models.DB.Model(&models.SubmissionTable{}).
		Where("problem_id = ?", pid).
		Select("student_id", "submission_category", "code_submitted_at", "completed").
		Find(&submissions).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve submissions for problem %d: %w", pid, err)
	}
	return submissions, nil
}

func GetStudentName(studentID int) string {
	var student models.Student
	if err := models.DB.Model(&models.Student{}).
		Where("id = ?", studentID).
		Select("name").
		First(&student).Error; err != nil {
		return ""
	}
	return student.Name
}

func GetCodeSnapshot(snapshotID int) (*models.CodeSnapshot, error) {
	var codeSnapshot models.CodeSnapshot
	if err := models.DB.Table("code_snapshots cs").
		Joins("join problems p on cs.problem_id = p.id").
		Where("cs.id = ?", snapshotID).
		Select("cs.student_id, cs.problem_id, cs.code, p.filename").
		First(&codeSnapshot).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve code snapshot for snapshot ID %d: %w", snapshotID, err)
	}
	return &codeSnapshot, nil
}

func GetCodeSnapshotMessageDetails(messageID int) ([]models.CodeSnapshotMessageDetails, error) {
	var details []models.CodeSnapshotMessageDetails
	if err := models.DB.Table("code_snapshots cs").
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
	if err := models.DB.Model(&models.SnapshotBackFeedback{}).
		Where("is_helpful = ? AND snapshot_feedback_id = ?", voteType, feedbackID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}
	return count, nil
}

func GetNumberOfReply(snapshotID int) (int, error) {
	var count int64
	if err := models.DB.Model(&models.SnapshotFeedback{}).
		Where("snapshot_id = ?", snapshotID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to retrieve the count of replies for snapshot ID %d: %w", snapshotID, err)
	}
	return int(count), nil
}

func GetAllTeachers() ([]models.Teacher, error) {
	var teachers []models.Teacher
	if err := models.DB.Model(&models.Teacher{}).Find(&teachers).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve teachers: %w", err)
	}
	return teachers, nil
}

func GetAllTeacherClasses() ([]models.TeacherClass, error) {
	var classes []models.TeacherClass
	if err := models.DB.Model(&models.TeacherClass{}).Find(&classes).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve teachers: %w", err)
	}
	return classes, nil
}

func GetAllStudentClasses() ([]models.StudentClass, error) {
	var classes []models.StudentClass
	if err := models.DB.Model(&models.StudentClass{}).Find(&classes).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve student classes: %w", err)
	}
	return classes, nil
}

func GetScoreDetails(problemID, studentID int) (*models.Score, error) {
	var score models.Score
	if err := models.DB.Model(&models.Score{}).
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

func GetProblemDetails(problemID int) (*models.Problem, error) {
	var problem models.Problem
	if err := models.DB.Model(&models.Problem{}).
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

func GetStudentStatus(studentID, problemID int) (*models.StudentStatus, error) {
	var status models.StudentStatus
	if err := models.DB.Model(&models.StudentStatus{}).
		Where("student_id = ? AND problem_id = ?", studentID, problemID).
		First(&status).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No record found
		}
		return nil, fmt.Errorf("failed to retrieve student status: %w", err)
	}
	return &status, nil
}

func GetStudentByID(studentID int) (*models.Student, error) {
	var student models.Student
	if err := models.DB.Model(&models.Student{}).
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
	var problem models.Problem
	if err := models.DB.Model(&models.Problem{}).
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
	if err := models.DB.Model(&models.TestCase{}).
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
	if err := models.DB.Model(&models.Attendance{}).
		Select("student_id").
		Where("DATE(attendance_at) = ?", date).
		Pluck("student_id", &studentIDs).Error; err != nil {
		return nil
	}
	return studentIDs
}

func GetAllStudents() map[int]string {
	var students []models.Student
	studentMap := make(map[int]string)

	if err := models.DB.Model(&models.Student{}).
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
	var stats models.ProblemStatistics

	if err := models.DB.Where("problem_id = ?", problemID).First(&stats).Error; err != nil {
		return 0, 0, 0, 0, 0
	}

	ungraded := stats.Submission - stats.GradedCorrect - stats.GradedIncorrect

	return stats.Active, stats.HelpRequest, ungraded, stats.GradedCorrect, stats.GradedIncorrect
}

func GetLatestSubmissionTime(problemID int) map[int]time.Time {
	var latestSubmissions = make(map[int]time.Time)
	var submissions []models.SubmissionTable

	if err := models.DB.Model(&models.SubmissionTable{}).
		Select("student_id, max(code_submitted_at) as code_submitted_at").
		Where("problem_id = ?", problemID).
		Group("student_id").
		Find(&submissions).Error; err != nil {
		log.Println(err)
		return latestSubmissions
	}

	for _, submission := range submissions {
		latestSubmissions[submission.StudentID] = submission.CodeSubmittedAt
	}

	return latestSubmissions
}

func GetProblemNameFromID(problemID int) string {
	var problem models.Problem

	if err := models.DB.Model(&models.Problem{}).
		Select("filename").
		Where("id = ?", problemID).
		First(&problem).Error; err != nil {
		return ""
	}

	return problem.Filename
}

func GetCodeSnapshotsByProblemID(problemID int) ([]models.CodeSnapshot, error) {
	var codeSnapshots []models.CodeSnapshot
	err := models.DB.Model(&models.CodeSnapshot{}).
		Select("student_id, MAX(last_updated_at) as last_updated_at").
		Where("problem_id = ?", problemID).
		Group("student_id").
		Find(&codeSnapshots).Error
	if err != nil {
		log.Printf("Failed to fetch code snapshots: %v", err)
		return codeSnapshots, err
	}
	return codeSnapshots, err
}
func GetLatestCodeSnapshots(problemID int) ([]models.CodeSnapshot, error) {
	var codeSnapshots []models.CodeSnapshot

	// Subquery to get the latest timestamp for each student
	subquery := models.DB.Model(&models.CodeSnapshot{}).
		Select("student_id, MAX(last_updated_at) as latest_timestamp"). // Changed timestamp to last_updated_at
		Where("problem_id = ?", problemID).
		Group("student_id")

	// Main query to get the code snapshots with the latest timestamp for each student
	err := models.DB.Model(&models.CodeSnapshot{}).
		Joins("JOIN (?) AS latest ON code_snapshots.student_id = latest.student_id AND code_snapshots.last_updated_at = latest.latest_timestamp", subquery). // Changed timestamp to last_updated_at
		Where("code_snapshots.problem_id = ?", problemID).
		Order("CASE WHEN event = 'at_submission' THEN 1 ELSE 2 END, last_updated_at DESC"). // Changed timestamp to last_updated_at
		Find(&codeSnapshots).Error

	if err != nil {
		log.Printf("Failed to fetch latest code snapshots: %v", err)
		return nil, err
	}

	return codeSnapshots, nil
}
func GetLatestCodeSnapshotForStudent(problemID int, studentID int) (*models.CodeSnapshot, error) {
	var codeSnapshot models.CodeSnapshot

	err := models.DB.
		Where("problem_id = ? AND student_id = ?", problemID, studentID).
		Order("last_updated_at DESC").
		Limit(1).
		Find(&codeSnapshot).Error

	if err != nil {
		log.Printf("Failed to fetch latest code snapshot for student %d: %v", studentID, err)
		return nil, err
	}

	return &codeSnapshot, nil
}

func GetLatestScaffoldingText(studentID, problemID int) (string, error) {
	var scaffolding string

	result := models.DB.
		Table("assigned_scaffoldings").
		Select("scaffolding").
		Where("student_id = ? AND problem_id = ?", studentID, problemID).
		Order("assigned_at DESC").
		Limit(1).
		Scan(&scaffolding)

	if result.Error != nil {
		return "", fmt.Errorf("failed to fetch latest scaffolding: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return "", nil // No scaffolding found
	}

	return scaffolding, nil
}

func GetProblemDetail(problemID int) (string, time.Time, error) {
	var problem struct {
		Description    string
		ProblemEndedAt time.Time
	}
	err := models.DB.Table("problems").
		Select("problem_description as description, problem_ended_at").
		Where("id = ?", problemID).
		Scan(&problem).Error
	return problem.Description, problem.ProblemEndedAt, err
}

func GetStudentStatusesByProblemID(problemID int) ([]models.StudentStatus, error) {
	var statuses []models.StudentStatus
	err := models.DB.Where("problem_id = ?", problemID).Find(&statuses).Error
	return statuses, err
}

func GetAnswerStats(problemID int) ([]*models.AnswerStatInfo, error) {
	var stats []struct {
		Answer string
		Count  int
	}
	err := models.DB.Table("submissions").
		Select("answer, COUNT(*) as count").
		Where("problem_id = ? AND answer IS NOT NULL AND LENGTH(answer) > 0", problemID).
		Group("answer").
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	var answerStats []*models.AnswerStatInfo
	var total int
	for _, stat := range stats {
		total += stat.Count
	}
	for _, stat := range stats {
		percent := float64(stat.Count) * 100.0 / float64(total)
		percent = math.Round(percent*100) / 100
		answerStats = append(answerStats, &models.AnswerStatInfo{
			Answer:  stat.Answer,
			Count:   stat.Count,
			Percent: percent,
		})
	}

	return answerStats, nil
}

func GetProblems(courseId string) ([]models.Problem, error) {
	var problems []models.Problem

	// Query the database to fetch problems filtered by CourseId
	if err := models.DB.Where("course_id = ?", courseId).
		Find(&problems).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch problems for course %s: %w", courseId, err)
	}

	return problems, nil
}

func GetTeacherByName(name string) (models.Teacher, error) {
	var teacher models.Teacher
	if err := models.DB.Where("name = ?", name).First(&teacher).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return teacher, fmt.Errorf("teacher not found")
		}
		return teacher, fmt.Errorf("failed to fetch teacher: %w", err)
	}
	return teacher, nil
}

func GetStudentByName(name string) (models.Student, error) {
	var student models.Student
	if err := models.DB.Where("name = ?", name).First(&student).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return student, fmt.Errorf("student not found")
		}
		return student, fmt.Errorf("failed to fetch student: %w", err)
	}
	return student, nil
}

func GetTags() (map[int]string, error) {
	var tags []models.Tag
	if err := models.DB.Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve tags: %w", err)
	}

	tagMap := make(map[int]string)
	for _, tag := range tags {
		tagMap[tag.ID] = tag.TopicDescription
	}

	return tagMap, nil
}

func GetStudentScores() (map[int]*models.ScoreEntry, error) {
	var scores []struct {
		Score                  int
		GradedSubmissionNumber int
		StudentID              int
		StudentName            string
	}

	if err := models.DB.Table("scores").
		Joins("join students on score.student_id = student.id").
		Select("scores.score, scores.graded_submission_number, scores.student_id, students.name").
		Find(&scores).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve scores: %w", err)
	}

	scoreEntries := make(map[int]*models.ScoreEntry)
	for _, score := range scores {
		if _, ok := scoreEntries[score.StudentID]; !ok {
			scoreEntries[score.StudentID] = &models.ScoreEntry{Name: score.StudentName}
		}
		scoreEntries[score.StudentID].Points += score.Score
		scoreEntries[score.StudentID].Attempts += score.GradedSubmissionNumber
		scoreEntries[score.StudentID].Count++
	}

	return scoreEntries, nil
}

func GetTagDescriptionByID(tagID int) (string, error) {
	var tag models.Tag
	if err := models.DB.Where("id = ?", tagID).First(&tag).Error; err != nil {
		return "", fmt.Errorf("failed to retrieve tag description: %w", err)
	}
	return tag.TopicDescription, nil
}

func GetProblemPerformanceByTagID(tagID int) (map[int]*models.ProblemPerformance, error) {
	var results []struct {
		Pid       int
		Merit     int
		At        time.Time
		Points    int
		StudentID int
	}

	if err := models.DB.Table("problems").
		Joins("join scores on problems.id = scores.problem_id").
		Joins("join students on scores.student_id = students.id").
		Where("problems.tag = ?", tagID).
		Select("problems.id, problems.merit, problems.at, scores.points, scores.student_id").
		Find(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve problem performance: %w", err)
	}

	record := make(map[int]*models.ProblemPerformance)
	for _, result := range results {
		if _, ok := record[result.Pid]; !ok {
			record[result.Pid] = &models.ProblemPerformance{
				Pid:       result.Pid,
				Timestamp: result.At.UnixNano(),
				Correct:   0,
				Incorrect: 0,
				Activity:  0,
				PC:        models.Passcode,
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
	if err := models.DB.Table("student").Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to get student count: %w", err)
	}
	return float32(count), nil
}

func GetLatestProblemID() (int, error) {
	var problem models.Problem
	// Retrieve the latest problem based on the highest ID
	if err := models.DB.Order("id desc").First(&problem).Error; err != nil {
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

	err := models.DB.Table("scores").
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
	var attendances []models.Attendance
	attendants := make(map[int]int)

	err := models.DB.Where("DATE(attendance_at) = ?", theDate).
		Find(&attendances).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve attendance: %w", err)
	}

	for _, att := range attendances {
		attendants[att.StudentID] = 0
	}

	return attendants, nil
}

func GetAttendanceByStudentID(uid int) ([]models.Attendance, error) {
	var attendances []models.Attendance
	err := models.DB.Where("student_id = ?", uid).
		Select("attendance_at").
		Find(&attendances).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve attendance for student ID %d: %w", uid, err)
	}
	return attendances, nil
}

func GetCurrentUserVote(feedbackID int, userID int, userRole string) string {
	var feedback models.MessageBackFeedback

	err := models.DB.Where("message_feedback_id = ? AND author_id = ? AND author_role = ?", feedbackID, userID, userRole).
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

	err := models.DB.Model(&models.MessageBackFeedback{}).
		Where("useful = ? AND message_feedback_id = ?", backFeedbackType, feedbackID).
		Count(&count).Error
	if err != nil {
		log.Println(err)
		return 0
	}

	return int(count)
}

func GetMessageFeedbacksByMessageID(messageID int) ([]models.MessageFeedback, error) {
	var messageFeedbacks []models.MessageFeedback
	err := models.DB.Where("message_id = ?", messageID).Find(&messageFeedbacks).Error
	if err != nil {
		return nil, fmt.Errorf("Error retrieving message feedbacks: %v", err)
	}
	return messageFeedbacks, nil
}

func GetLatestSnapshot(studentID int, problemID int) (*models.Snapshot, error) {
	var snapshot models.CodeSnapshot

	_ = models.DB.Where("student_id = ? AND problem_id = ?", studentID, problemID).
		Order("last_updated_at DESC").
		First(&snapshot).Error

	return &models.Snapshot{
		ID:          snapshot.ID,
		StudentID:   snapshot.StudentID,
		ProblemID:   snapshot.ProblemID,
		ProblemName: GetProblemNameFromID(problemID),
		Code:        snapshot.Code,
		LastUpdated: snapshot.LastUpdatedAt,
	}, nil
}

func GetTeacherName(authorID int) string {
	var teacher models.Teacher

	err := models.DB.Where("id = ?", authorID).First(&teacher).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ""
		}
		log.Printf("Error retrieving teacher with ID %d: %v", authorID, err)
		return ""
	}

	return teacher.Name
}

func FetchExistingMessageBackFeedback(feedbackID, authorID int, authorRole string) (*models.MessageBackFeedback, error) {
	var feedback models.MessageBackFeedback
	err := models.DB.Table("message_back_feedbacks").
		Where("message_feedback_id = ? AND author_id = ? AND author_role = ?", feedbackID, authorID, authorRole).
		First(&feedback).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to fetch existing message back feedback: %w", err)
	}

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &feedback, nil
}

func FetchExistingTestCase(studentID, problemID int) (int, error) {
	var testCaseID int
	err := models.DB.Table("test_cases").
		Where("student_id = ? AND problem_id = ?", studentID, problemID).
		Select("id").
		Scan(&testCaseID).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, fmt.Errorf("failed to fetch existing test case: %w", err)
	}

	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}

	return testCaseID, nil
}

func FetchTagIDByDescription(description string) (int64, error) {
	var tagID int64
	err := models.DB.Table("tags").
		Where("topic_description = ?", description).
		Select("id").
		Scan(&tagID).Error // Use Scan instead of First

	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, fmt.Errorf("failed to fetch tag ID: %w", err)
	}

	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}

	return tagID, nil
}

func CheckMessageBackFeedback(feedbackID, authorID int, authorRole string) (bool, error) {
	var count int64
	err := models.DB.Model(&models.MessageBackFeedback{}).
		Where("message_feedback_id = ? AND author_id = ? AND author_role = ?", feedbackID, authorID, authorRole).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check message back feedback: %w", err)
	}

	return count > 0, nil
}

func GetStudentIDByMessageID(messageID int) (int, error) {
	var studentID int
	err := models.DB.Model(&models.HelpMessage{}).
		Select("student_id").
		Where("id = ?", messageID).
		Scan(&studentID).Error

	if err != nil {
		return 0, fmt.Errorf("failed to retrieve student ID: %w", err)
	}

	return studentID, nil
}

func FetchScores() (map[string]map[int]int, error) {
	type ScoreData struct {
		StudentID int
		Filename  string
		Score     int
	}

	var results []ScoreData
	err := models.DB.Table("problems as P").
		Select("S.student_id, P.filename, S.score").
		Joins("join scores as S on P.id = S.problem_id").
		Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch scores: %w", err)
	}

	// Organize data into a nested map: filename -> studentID -> score
	data := make(map[string]map[int]int)
	for _, result := range results {
		if data[result.Filename] == nil {
			data[result.Filename] = make(map[int]int)
		}
		data[result.Filename][result.StudentID] = result.Score
	}

	return data, nil
}

func FetchStudentReport(uid int) ([]*models.StudentReport, error) {
	type ReportData struct {
		Points   int
		Date     time.Time
		Filename string
	}

	var results []ReportData
	err := models.DB.Table("scores").
		Select("scores.score as points, scores.score_given_at as date, problems.filename").
		Joins("join problems on problems.id = scores.problem_id").
		Where("student_id = ?", uid).
		Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch student report: %w", err)
	}

	// Convert to []*StudentReport
	report := make([]*models.StudentReport, len(results))
	for i, result := range results {
		report[i] = &models.StudentReport{
			Points:   result.Points,
			Filename: result.Filename,
			Date:     result.Date.Unix(),
		}
	}

	return report, nil
}

func FetchStudentStatus(problemID, studentID int) (*models.DashBoardStudentInfo, error) {
	var studentStatus models.StudentStatus
	err := models.DB.Where("problem_id = ? AND student_id = ?", problemID, studentID).First(&studentStatus).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &models.DashBoardStudentInfo{}, nil // No record found
		}
		return &models.DashBoardStudentInfo{}, fmt.Errorf("failed to fetch student status: %w", err)
	}

	return &models.DashBoardStudentInfo{
		CodingStat:     studentStatus.CodingStat,
		HelpStat:       studentStatus.HelpStat,
		SubmissionStat: studentStatus.SubmissionStat,
		Percentage:     studentStatus.Percentage,
		Explanation:    studentStatus.Explanation,
	}, nil
}

func FetchSubmissions(problemID, studentID int) ([]*models.SubmissionInfo, error) {
	var submissionInfos = make([]*models.SubmissionInfo, 0)
	var submissions []models.SubmissionTable
	err := models.DB.Where("student_id = ? AND problem_id = ?", studentID, problemID).Find(&submissions).Error
	if err != nil {
		return submissionInfos, fmt.Errorf("failed to fetch submissions: %w", err)
	}

	// Map the fetched submissions to the SubmissionInfo format
	for _, submission := range submissions {
		submissionInfos = append(submissionInfos, &models.SubmissionInfo{
			ID:          submission.ID,
			SnapshotID:  submission.SnapshotID,
			Code:        submission.StudentCode,
			Grade:       submission.Verdict,
			SubmittedAt: submission.CodeSubmittedAt,
		})
	}

	// Sort submissions according to the grading and submission time logic
	sort.SliceStable(submissionInfos, func(i, j int) bool {
		if submissionInfos[i].Grade == "" && submissionInfos[j].Grade == "" {
			return submissionInfos[i].SubmittedAt.After(submissionInfos[j].SubmittedAt)
		}
		if submissionInfos[i].Grade == "" {
			return true
		}
		if submissionInfos[j].Grade == "" {
			return false
		}
		return submissionInfos[i].SubmittedAt.After(submissionInfos[j].SubmittedAt)
	})

	return submissionInfos, nil
}

func FetchMessagesAndSnapshots(problemID, studentID int) ([]models.MessageWithSnapshot, error) {
	var result []models.MessageWithSnapshot

	err := models.DB.Table("messages M").
		Select("M.id, M.snapshot_id, M.message, M.author_id, M.author_role, M.given_at, M.type, C.Code, C.event").
		Joins("join code_snapshots C on M.snapshot_id = C.id").
		Where("C.problem_id = ? AND C.student_id = ?", problemID, studentID).
		Scan(&result).Error

	return result, err
}

func FetchSubmission(studentID, problemID int) ([]models.SubmissionTable, error) {
	var submissions []models.SubmissionTable
	err := models.DB.Where("student_id = ? AND problem_id = ?", studentID, problemID).
		Order("verdict, code_submitted_at ASC").
		Find(&submissions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch submissions: %w", err)
	}

	return submissions, nil
}

func FetchStudentStatuses(problemID, studentID int) (*models.DashBoardStudentInfo, error) {
	var result models.StudentStatus

	// Execute the query using the StudentStatus struct
	err := models.DB.Where("problem_id = ? AND student_id = ?", problemID, studentID).
		First(&result).Error // Use `First` to get the first matching result

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Convert the result into DashBoardStudentInfo format
	studentStats := &models.DashBoardStudentInfo{
		CodingStat:     result.CodingStat,
		HelpStat:       result.HelpStat,
		SubmissionStat: result.SubmissionStat,
	}

	return studentStats, nil
}

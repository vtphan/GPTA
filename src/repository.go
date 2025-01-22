package main

import (
	"fmt"
	"gorm.io/gorm"
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

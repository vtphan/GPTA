package main

import (
	"gorm.io/gorm"
	"time"
)

func AddStudent(name string, password string) error {
	student := Student{Name: name, Password: password}
	if err := Database.Create(&student).Error; err != nil {
		return err
	}
	return nil
}

func AddTeacher(name string, password string) error {
	teacher := Teacher{Name: name, Password: password}
	if err := Database.Create(&teacher).Error; err != nil {
		return err
	}
	return nil
}

func AddProblem(teacherID int, problemDescription, answer, filename string, merit, effort, attempts int, topicID int, tag int, problemUploadedAt time.Time) (Problem, error) {
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
	if err := Database.Create(&problem).Error; err != nil {
		return problem, err
	}

	return problem, nil
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
	if err := Database.Create(&submission).Error; err != nil {
		return err
	}
	return nil
}

func CompleteSubmission(id int, completed time.Time, verdict string) error {
	if err := Database.Model(&Submission{}).Where("id = ?", id).Updates(map[string]interface{}{
		"Completed": completed,
		"Verdict":   verdict,
	}).Error; err != nil {
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
	if err := Database.Create(&scoreEntry).Error; err != nil {
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
	if err := Database.Create(&feedbackEntry).Error; err != nil {
		return err
	}
	return nil
}

func UpdateScore(id int, teacherID int, score float64, gradedSubmissionNumber int) error {
	if err := Database.Model(&Score{}).Where("id = ?", id).Updates(map[string]interface{}{
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
	if err := Database.Create(&attendance).Error; err != nil {
		return err
	}
	return nil
}

func AddTag(topicDescription string) (Tag, error) {
	tag := Tag{TopicDescription: topicDescription}
	if err := Database.Create(&tag).Error; err != nil {
		return tag, err
	}
	return tag, nil
}

func AddTestCase(problemID, studentID int, testCases string, addedAt time.Time) error {
	testCase := TestCase{
		ProblemID: problemID,
		StudentID: studentID,
		TestCases: testCases,
		AddedAt:   addedAt,
	}
	if err := Database.Create(&testCase).Error; err != nil {
		return err
	}
	return nil
}

func UpdateTestCase(id int, testCases string, addedAt time.Time) error {
	if err := Database.Model(&TestCase{}).Where("id = ?", id).Updates(map[string]interface{}{
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
	if err := Database.Create(&helpSubmission).Error; err != nil {
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
	if err := Database.Create(&helpMessage).Error; err != nil {
		return err
	}
	return nil
}

func UpdateHelpMessage(id int, useful bool, updatedAt time.Time) error {
	if err := Database.Model(&HelpMessage{}).Where("id = ?", id).Updates(map[string]interface{}{
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
	if err := Database.Create(&codeSnapshot).Error; err != nil {
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
	if err := Database.Create(&snapshotFeedback).Error; err != nil {
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
	if err := Database.Create(&snapshotBackFeedback).Error; err != nil {
		return err
	}
	return nil
}

func UpdateSnapshotBackFeedback(snapshotFeedbackID int, isHelpful bool, givenAt time.Time) error {
	if err := Database.Model(&SnapshotBackFeedback{}).Where("snapshot_feedback_id = ?", snapshotFeedbackID).Updates(map[string]interface{}{
		"IsHelpful": isHelpful,
		"GivenAt":   givenAt,
	}).Error; err != nil {
		return err
	}
	return nil
}

func UpdateProblemEndTime(problemEndedAt time.Time, problemID int) error {
	if err := Database.Model(&Problem{}).Where("id = ?", problemID).Update("ProblemEndedAt", problemEndedAt).Error; err != nil {
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
	if err := Database.Create(&helpEligible).Error; err != nil {
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
	if err := Database.Create(&userEventLog).Error; err != nil {
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
	if err := Database.Create(&studentStatus).Error; err != nil {
		return err
	}
	return nil
}

func UpdateStudentCodingStat(studentID, problemID int, codingStat string, lastUpdatedAt time.Time) error {
	if err := Database.Model(&StudentStatus{}).Where("student_id = ? AND problem_id = ?", studentID, problemID).Update("CodingStat", codingStat).Update("LastUpdatedAt", lastUpdatedAt).Error; err != nil {
		return err
	}
	return nil
}

func UpdateStudentHelpStat(studentID, problemID int, helpStat string, lastUpdatedAt time.Time) error {
	if err := Database.Model(&StudentStatus{}).Where("student_id = ? AND problem_id = ?", studentID, problemID).Update("HelpStat", helpStat).Update("LastUpdatedAt", lastUpdatedAt).Error; err != nil {
		return err
	}
	return nil
}

func UpdateStudentSubmissionStat(studentID, problemID int, submissionStat string, lastUpdatedAt time.Time) error {
	if err := Database.Model(&StudentStatus{}).Where("student_id = ? AND problem_id = ?", studentID, problemID).Update("SubmissionStat", submissionStat).Update("LastUpdatedAt", lastUpdatedAt).Error; err != nil {
		return err
	}
	return nil
}

func UpdateStudentTutoringStat(studentID, problemID int, tutoringStat string, lastUpdatedAt time.Time) error {
	if err := Database.Model(&StudentStatus{}).Where("student_id = ? AND problem_id = ?", studentID, problemID).Update("TutoringStat", tutoringStat).Update("LastUpdatedAt", lastUpdatedAt).Error; err != nil {
		return err
	}
	return nil
}

func AddMessage(snapshotID int, message string, authorID int, authorRole string, givenAt time.Time, messageType int) (Message, error) {
	messageEntry := Message{
		SnapshotID: snapshotID,
		Message:    message,
		AuthorID:   authorID,
		AuthorRole: authorRole,
		GivenAt:    givenAt,
		Type:       messageType,
	}
	if err := Database.Create(&messageEntry).Error; err != nil {
		return messageEntry, err
	}
	return messageEntry, nil
}

func AddMessageFeedback(messageID int, feedback string, authorID int, authorRole string, givenAt time.Time) error {
	messageFeedback := MessageFeedback{
		MessageID:  messageID,
		Feedback:   feedback,
		AuthorID:   authorID,
		AuthorRole: authorRole,
		GivenAt:    givenAt,
	}
	if err := Database.Create(&messageFeedback).Error; err != nil {
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
	if err := Database.Create(&messageBackFeedback).Error; err != nil {
		return err
	}
	return nil
}

func UpdateMessageBackFeedback(messageFeedbackID int, useful bool, givenAt time.Time) error {
	if err := Database.Model(&MessageBackFeedback{}).Where("message_feedback_id = ?", messageFeedbackID).Updates(map[string]interface{}{
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
	if err := Database.Create(&problemStats).Error; err != nil {
		return err
	}
	return nil
}

func IncProblemStatActive(problemID int) error {
	if err := Database.Model(&ProblemStatistics{}).Where("problem_id = ?", problemID).UpdateColumn("Active", gorm.Expr("active + ?", 1)).Error; err != nil {
		return err
	}
	return nil
}

func IncProblemStatSubmission(problemID int) error {
	if err := Database.Model(&ProblemStatistics{}).Where("problem_id = ?", problemID).UpdateColumn("Submission", gorm.Expr("submission + ?", 1)).Error; err != nil {
		return err
	}
	return nil
}

func IncProblemStatHelp(problemID int) error {
	if err := Database.Model(&ProblemStatistics{}).Where("problem_id = ?", problemID).UpdateColumn("HelpRequest", gorm.Expr("help_request + ?", 1)).Error; err != nil {
		return err
	}
	return nil
}

func IncProblemStatGradedCorrect(problemID int) error {
	if err := Database.Model(&ProblemStatistics{}).Where("problem_id = ?", problemID).UpdateColumn("GradedCorrect", gorm.Expr("graded_correct + ?", 1)).Error; err != nil {
		return err
	}
	return nil
}

func IncProblemStatGradedIncorrect(problemID int) error {
	if err := Database.Model(&ProblemStatistics{}).Where("problem_id = ?", problemID).UpdateColumn("GradedIncorrect", gorm.Expr("graded_incorrect + ?", 1)).Error; err != nil {
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
	if err := Database.Create(&submission).Error; err != nil {
		return err
	}
	return nil
}

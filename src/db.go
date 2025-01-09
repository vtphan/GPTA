// Author: Vinhthuy Phan, 2018
package main

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func execSQL(s string) {
	if err := DB.Exec(s).Error; err != nil {
		log.Fatalf("failed to execute SQL: %v", err)
	}
}

// TODO - See migrations for pervios entries...that is change the table names from student to students
func create_tables() {
	execSQL("create table if not exists students (id INT AUTO_INCREMENT NOT NULL, name VARCHAR(100) unique, password VARCHAR(100), PRIMARY KEY (`id`))")
	execSQL("create table if not exists teachers (id INT AUTO_INCREMENT NOT NULL, name VARCHAR(100) unique, password VARCHAR(100), PRIMARY KEY (`id`))")
	execSQL("create table if not exists attendances (id INT AUTO_INCREMENT NOT NULL, student_id INT NOT NULL, attendance_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists tags (id INT AUTO_INCREMENT NOT NULL, topic_description VARCHAR(200) unique, PRIMARY KEY (`id`))")
	execSQL("create table if not exists problems (id INT AUTO_INCREMENT NOT NULL, teacher_id INT, problem_description text, answer text, filename text, merit INT, effort INT, attempts INT, topic_id INT, tag INT, problem_uploaded_at timestamp, problem_ended_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists submissions (id INT AUTO_INCREMENT NOT NULL, problem_id INT NOT NULL, student_id INT NOT NULL, student_code text, snapshot_id INT default 0, submission_category INT, code_submitted_at timestamp, completed timestamp, verdict text, attempt_number INT, answer text, PRIMARY KEY (`id`))")
	execSQL("create table if not exists scores (id INT AUTO_INCREMENT NOT NULL, problem_id INT NOT NULL, student_id INT, teacher_id INT, score INT, graded_submission_number INT, score_given_at timestamp, unique(problem_id,student_id), PRIMARY KEY (`id`))")
	execSQL("create table if not exists feedbacks (id INT AUTO_INCREMENT NOT NULL, teacher_id INT, student_id INT, feedback text, feedback_given_at timestamp, submission_id INT, PRIMARY KEY (`id`))")
	execSQL("create table if not exists test_cases (id INT AUTO_INCREMENT NOT NULL, problem_id INT, student_id INT, test_cases text, added_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists code_explanations (id INT AUTO_INCREMENT NOT NULL, problem_id INT, student_id INT, snapshot_id INT, trying_what text, need_help_with text, code_submitted_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists help_messages (id INT AUTO_INCREMENT NOT NULL, code_explanation_id INT, student_id INT, message text, given_at timestamp, useful text, updated_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists code_snapshots (id INT AUTO_INCREMENT NOT NULL, student_id INT, problem_id INT, code text, last_updated_at timestamp, status int default 0, event VARCHAR(50), PRIMARY KEY (`id`))") // 0 = not submitted, 1 = submitted but not graded, 2 = submitted and incorrect, 3 = submitted and correct
	execSQL("create table if not exists snapshot_feedbacks (id INT AUTO_INCREMENT NOT NULL, snapshot_id INT, feedback text, author_id INT, author_role VARCHAR(50), given_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists snapshot_back_feedbacks (id INT AUTO_INCREMENT NOT NULL, snapshot_feedback_id INT, author_id INT, author_role VARCHAR(50), is_helpful VARCHAR(50), given_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists messages (id INT AUTO_INCREMENT NOT NULL, snapshot_id INT, message text, author_id INT, author_role VARCHAR(50), given_at timestamp, type INT, PRIMARY KEY (`id`))")
	execSQL("create table if not exists message_feedbacks (id INT AUTO_INCREMENT NOT NULL, message_id INT, feedback text, author_id INT, author_role VARCHAR(50), given_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists message_back_feedbacks (id INT AUTO_INCREMENT NOT NULL, message_feedback_id INT, author_id INT, author_role VARCHAR(50), useful VARCHAR(50), given_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists help_eligibles (id INT AUTO_INCREMENT NOT NULL, problem_id INT, student_id INT, became_eligible_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists user_event_logs (id INT AUTO_INCREMENT NOT NULL, name VARCHAR(50), user_id INT, user_type VARCHAR(50), event_type VARCHAR(50), referral_info VARCHAR(50), event_time timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists student_statuses (id INT AUTO_INCREMENT NOT NULL, student_id INT, problem_id INT, coding_stat VARCHAR(50), help_stat VARCHAR(50), submission_stat VARCHAR(50), tutoring_stat VARCHAR(50), last_updated_at timestamp, PRIMARY KEY (`id`))")
	execSQL("create table if not exists problem_statistics (id INT AUTO_INCREMENT NOT NULL, problem_id INT not null, active INT default 0, submission INT default 0, help_request INT default 0, graded_correct INT default 0, graded_incorrect INT default 0, PRIMARY KEY (`id`))")
	// foreign key example: http://www.sqlitetutorial.net/sqlite-foreign-key/
}

// -----------------------------------------------------------------
func init_database(db_name string, username string, pass string, server string) {
	var err error

	// Prepare DSN (Data Source Name) for GORM
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/?parseTime=true", username, pass, server)

	// Open a connection to MySQL using GORM
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Create the database if it doesn't exist
	err = DB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", db_name)).Error
	if err != nil {
		log.Fatal("Failed to create database: ", err)
	}

	// Switch to the selected database
	dsn = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", username, pass, server, db_name)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the selected database: ", err)
	}

	// Set connection pool settings (like MaxLifetime) if necessary
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get raw SQL database object: ", err)
	}
	sqlDB.SetConnMaxLifetime(time.Minute * 3)
	create_tables()
	Passcode = RandStringRunes(12)
	Students[0] = &StudenInfo{
		Boards: make([]*Board, 0),
	}
}

// -----------------------------------------------------------------
// Add or update score based on a decision. If decision is "correct"
// a new problem, if there's one, is added to student's board.
// -----------------------------------------------------------------
func addOrUpdateScore(decision string, pid, studentID, teacherID, partialCredits int) string {
	var score Score
	var problem Problem
	var message string

	// Retrieve score information for this student and problem
	if err := DB.Where("problem_id = ? AND student_id = ?", pid, studentID).First(&score).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Sprintf("Unable to retrieve score: %v", err)
	}

	// Retrieve merit and effort points for the problem
	if err := DB.First(&problem, pid).Error; err != nil {
		return fmt.Sprintf("Unable to retrieve problem: %v", err)
	}

	// Determine points for this student
	points := 0
	teacher := teacherID

	if decision == "correct" {
		points = problem.Merit
		message = "Answer is correct."
	} else {
		if partialCredits < problem.Merit {
			points = partialCredits
		} else {
			points = problem.Effort
		}

		// Ensure points are not reduced if previously graded correct
		if points < score.Score {
			points = score.Score
			teacher = score.TeacherID
		}
		message = "Answer is incorrect."
	}
	currentTime := time.Now()
	// Update or create a new score record
	if score.ID == 0 {
		newScore := Score{
			ProblemID:              pid,
			StudentID:              studentID,
			TeacherID:              teacher,
			Score:                  points,
			GradedSubmissionNumber: score.GradedSubmissionNumber + 1,
			ScoreGivenAt:           &currentTime,
		}
		if err := DB.Create(&newScore).Error; err != nil {
			return fmt.Sprintf("Unable to add score: %v", err)
		}
	} else {
		if err := DB.Model(&score).Updates(Score{
			TeacherID:              teacher,
			Score:                  points,
			GradedSubmissionNumber: score.GradedSubmissionNumber + 1,
		}).Error; err != nil {
			return fmt.Sprintf("Unable to update score: %v", err)
		}
	}
	return message
}

func addOrUpdateStudentStatus(studentID, problemID int, codingStat, helpStat, submissionStat, tutoringStat string) {
	var studentStatus StudentStatus

	if err := DB.Where("student_id = ? AND problem_id = ?", studentID, problemID).First(&studentStatus).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newStatus := StudentStatus{
				StudentID:      studentID,
				ProblemID:      problemID,
				CodingStat:     codingStat,
				HelpStat:       helpStat,
				SubmissionStat: submissionStat,
				TutoringStat:   tutoringStat,
				LastUpdatedAt:  time.Now(),
			}
			if err := DB.Create(&newStatus).Error; err != nil {
				log.Fatalf("Unable to add student status: %v", err)
			}
		} else {
			log.Fatalf("Error retrieving student status: %v", err)
		}
	} else {
		updates := map[string]interface{}{
			"updated_at": time.Now(),
		}
		if codingStat != "" {
			updates["coding_stat"] = codingStat
		}
		if helpStat != "" {
			updates["help_stat"] = helpStat
		}
		if submissionStat != "" {
			updates["submission_stat"] = submissionStat
		}
		if tutoringStat != "" {
			updates["tutoring_stat"] = tutoringStat
		}
		if err := DB.Model(&studentStatus).Updates(updates).Error; err != nil {
			log.Fatalf("Unable to update student status: %v", err)
		}
	}
}

// -----------------------------------------------------------------
func init_teacher(id int, name string, password string) {
	Teachers[id] = password
	TeacherPass[name] = password
	TeacherNameToId[name] = id
	TeacherIdToName[id] = name
	SeenHelpSubmissions[id] = map[int]bool{}
}

// -----------------------------------------------------------------
// initialize once per session
// -----------------------------------------------------------------
func init_student(student_id int, name string, password string) {
	err := AddAttendance(student_id, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	BoardsSem.Lock()
	defer BoardsSem.Unlock()

	Students[student_id] = &StudenInfo{
		Name:                  name,
		Password:              password,
		Boards:                make([]*Board, 0),
		SubmissionStatus:      make([]*StudentSubmissionStatus, 0),
		SnapShotFeedbackQueue: make([]*SnapShotFeedback, 0),
		ThankStatus:           0,
	}

	// Student[student_id] = password
	// MessageBoards[student_id] = ""
	// Boards[student_id] = make([]*Board, 0)

	for i := 0; i < len(Students[0].Boards); i++ {
		b := &Board{
			Content:      Students[0].Boards[i].Content,
			Answer:       Students[0].Boards[i].Answer,
			Attempts:     Students[0].Boards[i].Attempts,
			Filename:     Students[0].Boards[i].Filename,
			Pid:          Students[0].Boards[i].Pid,
			StartingTime: time.Now(),
		}
		Students[student_id].Boards = append(Students[student_id].Boards, b)
	}
	StudentSnapshot[student_id] = map[int]int{}
}

// -----------------------------------------------------------------
func loadAndAuthorizeStudent(studentID int, password string) bool {
	var student Student

	if err := DB.First(&student, studentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
		log.Fatalf("Error retrieving student: %v", err)
	}

	if student.Password != password {
		return false
	}

	init_student(student.ID, student.Name, password)
	return true
}

// -----------------------------------------------------------------
func LoadTeachers() {
	var teachers []Teacher
	if err := DB.Find(&teachers).Error; err != nil {
		log.Fatalf("Unable to load teachers: %v", err)
	}

	for _, teacher := range teachers {
		Teachers[teacher.ID] = teacher.Password
		TeacherPass[teacher.Name] = teacher.Password
		TeacherNameToId[teacher.Name] = teacher.ID
		TeacherIdToName[teacher.ID] = teacher.Name
	}
	Passcode = RandStringRunes(20)
}

// Author: Vinhthuy Phan, 2018
package repository

import (
	"fmt"
	"github.com/GPTA/src/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

func execSQL(s string) {
	if err := models.DB.Exec(s).Error; err != nil {
		log.Fatalf("failed to execute SQL: %v", err)
	}
}

func create_tables() {
	// Function to check if a table exists and rename it if necessary
	renameTableIfExists := func(oldName, newName string) {
		// Check if the table with the old name exists
		if checkIfTableExists(oldName) {
			// Check if the new table already exists to avoid renaming to an existing table
			if !checkIfTableExists(newName) {
				// Rename the old table to the new name
				models.DB.Exec(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", oldName, newName))
			}
		}
	}

	// Rename singular tables to plural names if they exist
	renameTableIfExists("teacher", "teachers")
	renameTableIfExists("student", "students")
	renameTableIfExists("attendance", "attendances")
	renameTableIfExists("tag", "tags")
	renameTableIfExists("problem", "problems")
	renameTableIfExists("submission", "submissions")
	renameTableIfExists("score", "scores")
	renameTableIfExists("feedback", "feedbacks")
	renameTableIfExists("test_case", "test_cases")
	renameTableIfExists("code_explanation", "code_explanations")
	renameTableIfExists("help_message", "help_messages")
	renameTableIfExists("code_snapshot", "code_snapshots")
	renameTableIfExists("snapshot_feedback", "snapshot_feedbacks")
	renameTableIfExists("snapshot_back_feedback", "snapshot_back_feedbacks")
	renameTableIfExists("message", "messages")
	renameTableIfExists("message_feedback", "message_feedbacks")
	renameTableIfExists("message_back_feedback", "message_back_feedbacks")
	renameTableIfExists("help_eligible", "help_eligibles")
	renameTableIfExists("user_event_log", "user_event_logs")
	renameTableIfExists("student_status", "student_statuses")
	renameTableIfExists("problem_statistic", "problem_statistics")

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

func checkIfTableExists(tableName string) bool {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = '%s'", tableName)
	models.DB.Raw(query).Scan(&count)
	return count > 0
}

// -----------------------------------------------------------------
func InitDatabase(db_name string, username string, pass string, server string) {
	var err error

	// Prepare DSN (Data Source Name) for GORM
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/?parseTime=true", username, pass, server)

	// Open a connection to MySQL using GORM
	models.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Create the database if it doesn't exist
	err = models.DB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", db_name)).Error
	if err != nil {
		log.Fatal("Failed to create database: ", err)
	}

	// Switch to the selected database
	dsn = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", username, pass, server, db_name)
	models.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the selected database: ", err)
	}

	// Set connection pool settings (like MaxLifetime) if necessary
	sqlDB, err := models.DB.DB()
	if err != nil {
		log.Fatal("Failed to get raw SQL database object: ", err)
	}
	sqlDB.SetConnMaxLifetime(time.Minute * 3)
	create_tables()
	models.Passcode = models.RandStringRunes(12)
	models.Students[0] = &models.StudenInfo{
		Boards: make([]*models.Board, 0),
	}
}

// -----------------------------------------------------------------
// Add or update score based on a decision. If decision is "correct"
// a new problem, if there's one, is added to student's board.
// -----------------------------------------------------------------
func AddOrUpdateScore(decision string, pid, student_id, teacher_id, partial_credits int) string {
	mesg := ""

	// Find score information for this student (student_id) for this problem (pid)

	score_id, current_points, current_attempts, current_tid := 0, 0, 0, 0
	score, _ := GetScoreDetails(pid, student_id)
	if score != nil {
		// Assign values from the returned Score object
		score_id = score.ID
		current_points = score.Score
		current_attempts = score.GradedSubmissionNumber
		current_tid = score.TeacherID
	}
	// Find merit points and effort points for this problem (pid)
	merit, effort := 0, 0
	problem, _ := GetProblemDetails(pid)
	if problem != nil {
		merit = problem.Merit
		effort = problem.Effort
	}
	// Determine points for this student
	points, teacher := 0, teacher_id
	if decision == "correct" {
		points = merit
		mesg = "Answer is correct."
	} else {
		if partial_credits < merit {
			points = partial_credits
		} else {
			points = effort
		}

		// If the problem was previously graded correct, this submission
		// does not reduce it.  Grading is asynchronous.
		if points < current_points {
			points = current_points
			teacher = current_tid
		}
		mesg = "Answer is incorrect."
	}
	// m := add_next_problem_to_board(pid, student_id, decision)
	// mesg = mesg + m

	// Add a new score or update a current score for this student & problem
	if score_id == 0 {
		_, err := AddScore(pid, student_id, teacher_id, points, current_attempts+1, time.Now())
		if err != nil {
			mesg = fmt.Sprintf("Unable to add score: %d %d %d", pid, student_id, teacher_id)
			models.WriteLog(models.Config.LogFile, mesg)
			return mesg
		}
	} else {
		err := UpdateScore(teacher, points, current_attempts+1, score_id)
		if err != nil {
			mesg = fmt.Sprintf("Unable to update score: %d %d", teacher, score_id)
			models.WriteLog(models.Config.LogFile, mesg)
			return mesg
		}
	}
	return mesg
}

func AddOrUpdateStudentStatus(studentID int, problemID int, codingStat, helpStat, submissionStat, tutoringStat string) {
	status, err := GetStudentStatus(studentID, problemID)
	if err != nil {
		log.Fatalf("Error retrieving student status: %v", err)
	}

	now := time.Now()
	if status != nil {
		// Record exists: Update fields
		if codingStat != "" {
			err = UpdateStudentCodingStat(codingStat, now, studentID, problemID)
			if err != nil {
				log.Fatalf("Error updating coding stat: %v", err)
			}
		}
		if helpStat != "" {
			err = UpdateStudentHelpStat(helpStat, now, studentID, problemID)
			if err != nil {
				log.Fatalf("Error updating help stat: %v", err)
			}
		}
		if submissionStat != "" {
			err = UpdateStudentSubmissionStat(submissionStat, now, studentID, problemID)
			if err != nil {
				log.Fatalf("Error updating submission stat: %v", err)
			}
		}
		if tutoringStat != "" {
			err = UpdateStudentTutoringStat(tutoringStat, now, studentID, problemID)
			if err != nil {
				log.Fatalf("Error updating tutoring stat: %v", err)
			}
		}
	} else {
		// No record exists: Insert new record
		_, err = AddStudentStatus(studentID, problemID, codingStat, helpStat, submissionStat, tutoringStat, now)
		if err != nil {
			log.Fatalf("Error adding student status: %v", err)
		}
	}
}

// -----------------------------------------------------------------
func InitTeacher(id int, name string, password string) {
	models.TeacherMap[id] = password
	models.TeacherPass[name] = password
	models.TeacherNameToId[name] = id
	models.TeacherIdToName[id] = name
	models.SeenHelpSubmissions[id] = map[int]bool{}
}

// -----------------------------------------------------------------
// initialize once per session
// -----------------------------------------------------------------
func InitStudent(student_id int, name string, password string) {
	_, err := AddAttendance(student_id, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	models.BoardsSem.Lock()
	defer models.BoardsSem.Unlock()

	models.Students[student_id] = &models.StudenInfo{
		Name:                  name,
		Password:              password,
		Boards:                make([]*models.Board, 0),
		SubmissionStatus:      make([]*models.StudentSubmissionStatus, 0),
		SnapShotFeedbackQueue: make([]*models.SnapShotFeedback, 0),
		ThankStatus:           0,
	}

	// Student[student_id] = password
	// MessageBoards[student_id] = ""
	// Boards[student_id] = make([]*Board, 0)

	for i := 0; i < len(models.Students[0].Boards); i++ {
		b := &models.Board{
			Content:      models.Students[0].Boards[i].Content,
			Answer:       models.Students[0].Boards[i].Answer,
			Attempts:     models.Students[0].Boards[i].Attempts,
			Filename:     models.Students[0].Boards[i].Filename,
			Pid:          models.Students[0].Boards[i].Pid,
			StartingTime: time.Now(),
		}
		models.Students[student_id].Boards = append(models.Students[student_id].Boards, b)
	}
	models.StudentSnapshot[student_id] = map[int]int{}
}

// -----------------------------------------------------------------
func LoadAndAuthorizeStudent(studentID int, password string) bool {
	student, err := GetStudentByID(studentID)
	if err != nil {
		log.Fatalf("Error retrieving student: %v", err)
	}
	if student == nil || student.Password != password {
		return false
	}
	InitStudent(studentID, student.Name, password)
	return true
}

// -----------------------------------------------------------------
func LoadTeachers() {
	teachers, err := GetAllTeachers()
	if err != nil {
		log.Fatalf("Error loading teachers: %v", err)
	}

	for _, teacher := range teachers {
		models.TeacherMap[teacher.ID] = teacher.Password
		models.TeacherPass[teacher.Name] = teacher.Password
		models.TeacherNameToId[teacher.Name] = teacher.ID
		models.TeacherIdToName[teacher.ID] = teacher.Name
	}

	models.Passcode = models.RandStringRunes(20)
}

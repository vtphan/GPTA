// Author: Vinhthuy Phan, 2018
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os/signal"
	"strings"
	"syscall"

	"github.com/GPTA/src/models"
	"github.com/GPTA/src/openAI"
	"github.com/GPTA/src/repository"
	"github.com/GPTA/src/restHandlers"
	"gorm.io/gorm"

	_ "github.com/mattn/go-sqlite3"

	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"
)

// -----------------------------------------------------------------
func init_handlers() {
	http.HandleFunc("/test", restHandlers.TestHandler)

	// Analytics
	// http.HandleFunc("/learning_report", learning_reportHandler)
	http.HandleFunc("/analyze_submissions", restHandlers.AnalyzeSubmissionsHandler)
	http.HandleFunc("/view_activities", restHandlers.ViewActivitiesHandler)
	http.HandleFunc("/report", restHandlers.ReportHandler)
	http.HandleFunc("/report_tag", restHandlers.ReportTagHandler)
	http.HandleFunc("/view_answers", restHandlers.ViewAnswersHandler)
	http.HandleFunc("/statistics", restHandlers.StatisticsHandler)

	// Others
	http.HandleFunc("/student_periodic_update", Authorize(restHandlers.StudentPeriodicUpdateHandler, "student"))

	http.HandleFunc("/student_gets_report", Authorize(restHandlers.StudentGetsReportHandler, "student"))
	http.HandleFunc("/student_checks_in", Authorize(restHandlers.StudentChecksInHandler, "student"))
	http.HandleFunc("/student_shares", Authorize(restHandlers.StudentSharesHandler, "student"))
	http.HandleFunc("/student_gets", Authorize(restHandlers.StudentGetsHandler, "student"))
	http.HandleFunc("/view_bulletin_board", restHandlers.ViewBulletinBoardHandler)
	http.HandleFunc("/remove_bulletin_page", restHandlers.RemoveBulletinPageHandler)
	http.HandleFunc("/bulletin_board_data", restHandlers.BulletinBoardDataHandler)
	http.HandleFunc("/complete_registration", restHandlers.CompleteRegistrationHandler)

	http.HandleFunc("/student_ask_help", Authorize(restHandlers.StudentAskHelpHandler, "student"))
	http.HandleFunc("/student_get_help_code", Authorize(restHandlers.StudentGetHelpCode, "student"))
	http.HandleFunc("/student_return_without_feedback", Authorize(restHandlers.StudentReturnWithoutFeedbackHandler, "student"))
	http.HandleFunc("/student_send_help_message", Authorize(restHandlers.StudentSendHelpMessageHandler, "student"))
	http.HandleFunc("/student_send_thank_you", Authorize(restHandlers.SendThankYouHandler, "student"))

	http.HandleFunc("/teacher_get_help_code", Authorize(restHandlers.TeacherGetHelpCode, "teacher"))
	http.HandleFunc("/teacher_return_without_feedback", Authorize(restHandlers.TeacherReturnWithoutFeedbackHandler, "teacher"))
	http.HandleFunc("/teacher_send_help_message", Authorize(restHandlers.TeacherSendHelpMessageHandler, "teacher"))

	http.HandleFunc("/teacher_gets_queue", Authorize(restHandlers.TeacherGetsQueueHandler, "teacher"))
	http.HandleFunc("/teacher_adds_bulletin_page", Authorize(restHandlers.TeacherAddsBulletinPageHandler, "teacher"))
	http.HandleFunc("/teacher_clears_submissions", Authorize(restHandlers.TeacherClearsSubmissionsHandler, "teacher"))
	http.HandleFunc("/teacher_deactivates_problems", Authorize(restHandlers.TeacherDeactivatesProblemsHandler, "teacher"))
	http.HandleFunc("/teacher_grades", Authorize(restHandlers.TeacherGradesHandler, "teacher"))
	http.HandleFunc("/teacher_puts_back", Authorize(restHandlers.TeacherPutsBackHandler, "teacher"))
	http.HandleFunc("/teacher_gets", Authorize(restHandlers.TeacherGetsHandler, "teacher"))
	http.HandleFunc("/teacher_broadcasts", Authorize(restHandlers.TeacherBroadcastsHandler, "teacher"))
	http.HandleFunc("/teacher_gets_passcode", Authorize(restHandlers.TeacherGetsPasscodeHandler, "teacher"))
	http.HandleFunc("/student_gets_passcode", Authorize(restHandlers.StudentGetsPasscodeHandler, "student"))
	http.Handle("/scaffolding-guidance-html.html", http.FileServer(http.Dir("/Users/shashwatdadhich/go/src/github.com/GPTA")))
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "pong") })

	http.HandleFunc("/get_testcase", Authorize(restHandlers.Testcase_getsHandler, ""))

	http.HandleFunc("/code_snapshot", Authorize(restHandlers.CodeSnapshotHandler, ""))
	http.HandleFunc("/get_global_info", Authorize(restHandlers.GlobalInfoHandler, ""))

	// http.HandleFunc("/get_codespace", Authorize(codespaceHandler))
	// http.HandleFunc("/get_snapshot", Authorize(getCodeSnapshotHandler))
	http.HandleFunc("/save_snapshot_feedback", Authorize(restHandlers.CodeSnapshotFeedbackHandler, ""))
	http.HandleFunc("/get_snapshot_feedback", Authorize(restHandlers.GetSnapshotFeedbackHandler, ""))
	http.HandleFunc("/save_snapshot_back_feedback", Authorize(restHandlers.StudentSendBackFeedbackHandler, ""))
	// http.HandleFunc("/student_views_feedback", Authorize(studentViewsFeedbackHandler))
	// http.HandleFunc("/teacher_views_feedback", Authorize(teacherViewsFeedbackHandler))

	// http.HandleFunc("/help_requests", Authorize(helpRequestListHandler))
	// http.HandleFunc("/view_help_request", Authorize(viewHelpRequestHandler))
	http.HandleFunc("/set_peer_tutor", Authorize(restHandlers.SetPeerTutorHandler, "teacher"))

	http.HandleFunc("/view_exercises", Authorize(restHandlers.ProblemListHandler, ""))
	http.HandleFunc("/view_feedback", Authorize(restHandlers.ExerciseListHandler, ""))
	http.HandleFunc("/view_user_feedback", Authorize(restHandlers.StudentFeedbackProvisionHandler, ""))
	http.HandleFunc("/problem_dashboard", Authorize(restHandlers.ProblemDashboardHandler, ""))
	http.HandleFunc("/scaffolding_dashboard", Authorize(restHandlers.ScaffoldingDashboardHandler, ""))
	http.HandleFunc("/student_dashboard_feedback_provision", Authorize(restHandlers.StudentDashboardFeedbackProvisionHandler, ""))
	http.HandleFunc("/save_message_feedback", Authorize(restHandlers.MessageFeedbackHandler, ""))
	http.HandleFunc("/student_dashboard_submissions", Authorize(restHandlers.StudentDashboardSubmissionHandler, ""))
	http.HandleFunc("/has_message_feedback", Authorize(restHandlers.HasMessageBackFeedbackHandler, ""))
	http.HandleFunc("/teacher_signin_complete", restHandlers.TeacherSigninCompleteHandler)
	http.HandleFunc("/admin_signin", restHandlers.AdminSigninHandler)
	http.HandleFunc("/admin_dashboard", restHandlers.AdminDashboardHandler)
	http.HandleFunc("/teacher_signin", restHandlers.TeacherSigninHandler)
	http.HandleFunc("/settings_view", restHandlers.SettingsViewHandler)
	http.HandleFunc("/ai_settings_view", restHandlers.AISettingsViewHandler)
	http.HandleFunc("/prompt_view", restHandlers.PromptViewHandler)
	http.HandleFunc("/vi_view", restHandlers.ViViewHandler)
	http.HandleFunc("/sc_view", restHandlers.ScViewHandler)
	http.HandleFunc("/assign_scaffold", restHandlers.AssignScaffoldHandler)
	http.HandleFunc("/teacher_web_broadcast", Authorize(restHandlers.TeacherWebBroadcastHandler, "teacher"))
	http.HandleFunc("/student_dashboard_code_snapshot", Authorize(restHandlers.StudentDashboardCodeSpaceHandler, ""))
	http.HandleFunc("/teacher_exports_point", Authorize(restHandlers.ExportPointsHandler, "teacher"))
	http.HandleFunc("/", restHandlers.IndexHandler)
	http.HandleFunc("/index", restHandlers.IndexHandler)
	http.HandleFunc("/peer_tutoring", Authorize(restHandlers.PeerTutorHandler, "student"))
	http.HandleFunc("/process_code_with_prompt", openAI.ProcessCodeWithPromptHandler)
	http.HandleFunc("/summarize_class_performance", openAI.SummarizeClassPerformance)
	http.HandleFunc("/process_scaffolding", openAI.SummarizeScaffolding)
	http.HandleFunc("/get_scaffolding", openAI.GetScaffolding)
	http.HandleFunc("/get_scaffolding_material", openAI.GetScaffoldingMaterial)
	http.HandleFunc("/summarize_student_progress", openAI.SummarizeStudentProgress)
	http.HandleFunc("/get_class_feedback", openAI.GetLatestFeedbackByProblemID)
	http.HandleFunc("/get_feedback_by_id", openAI.GetFeedbackByFeedbackID)
	http.HandleFunc("/get_courses", restHandlers.GetCoursesHandler)
	http.HandleFunc("/add_course", restHandlers.AddCourseHandler)
	http.HandleFunc("/add_teacher", restHandlers.AddTeacherHandler)
	http.HandleFunc("/add_students", restHandlers.AddStudentsHandler)
	http.HandleFunc("/get_feedback_list", openAI.ListFeedbackHistoryByProblemID)
	http.HandleFunc("/add_api_key", openAI.AddAPIKey)
	http.HandleFunc("/logout", LogoutHandler)
}

// -----------------------------------------------------------------
func informIPAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}
	// for _, a := range addrs {
	// if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() {
	// if ipnet.IP.To4() != nil {
	// return ipnet.IP.String()
	// }
	// }
	// }
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() {
			ip4 := ipnet.IP.To4()
			if ip4 != nil {
				switch {
				// case ip4[0] == 10:
				case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
				case ip4[0] == 192 && ip4[1] == 168:
				default:
					return ip4.String()
				}
			}
		}
	}
	return ""
}

// -----------------------------------------------------------------
func init_config(filename string) *models.Configuration {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	config := &models.Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	if config.IP == "" {
		config.IP = informIPAddress()
	}
	config.Address = fmt.Sprintf("%s:%d", config.IP, config.Port)
	// if config.PeerTutor == 1 {
	// 	PeerTutorAllowed = true
	// } else {
	// 	PeerTutorAllowed = false
	// }
	return config
}

// -----------------------------------------------------------------
func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	rand.Seed(time.Now().UnixNano())
	config_file, teacher_file, student_file, course_file := "/Users/ajay/MS/GPTA/Examples/gem_config.json", "/Users/ajay/MS/GPTA/Examples/teachers.txt", "/Users/ajay/MS/GPTA/Examples/students.txt", "/Users/ajay/MS/GPTA/Examples/courses.txt"
	ai_prompts_file := "/Users/ajay/MS/GPTA/src/prompts"
	flag.StringVar(&config_file, "c", config_file, "json-formatted configuration file.")
	flag.StringVar(&teacher_file, "add_teachers", teacher_file, "teacher file.")
	flag.StringVar(&student_file, "add_students", student_file, "student file.")
	flag.StringVar(&course_file, "add_courses", course_file, "courses file.") // New flag for courses
	flag.Parse()
	if config_file == "" {
		flag.Usage()
		os.Exit(1)
	}
	models.Config = init_config(config_file)
	repository.InitDatabase(models.Config.Database, models.Config.DBUserName, models.Config.DBPassWord, models.Config.DBServerIP)
	setupGracefulShutdown()
	ReloadGlobalMaps()
	if course_file != "" {
		restHandlers.AddMultipleCourses(course_file) // Call AddMultipleCourses to add courses
	}
	if teacher_file != "" {
		restHandlers.AddMultiple(teacher_file, "teacher")
	}
	if student_file != "" {
		restHandlers.AddMultiple(student_file, "student")
	}
	if ai_prompts_file != "" {
		restHandlers.LoadAIPromptsFromFiles(ai_prompts_file)
	}
	init_handlers()
	courses, err := repository.LoadTeachers()
	if err != nil {
		log.Fatal("Failed to load teachers: ", err)
	}

	// Join unique courses with commas
	courseList := strings.Join(courses, ", ")
	fmt.Println("**************************************************")

	fmt.Printf("*   Course List:      %s\n", courseList)
	fmt.Printf("*   Serving at:     http://%s\n", models.Config.Address)
	fmt.Printf("*   GEM %s\n", VERSION)
	fmt.Println("**************************************************\n")
	err = http.ListenAndServe(models.Config.Address, nil)
	if err != nil {
		log.Fatal("Unable to serve gem server at " + models.Config.Address)
	}
}

func setupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-c
		fmt.Println("Shutting down gracefully...")
		err := StoreGlobalMaps()
		if err != nil {
			log.Println(err)
		}
		os.Exit(0)
	}()
}

type GlobalMap struct {
	ID                   int    `gorm:"primaryKey;autoIncrement"`
	TeacherMap           string `gorm:"type:text"` // JSON string for map[int]string
	TeacherPass          string `gorm:"type:text"` // JSON string for map[string]string
	TeacherNameToId      string `gorm:"type:text"` // JSON string for map[string]int
	TeacherIdToName      string `gorm:"type:text"` // JSON string for map[int]string
	Students             string `gorm:"type:text"` // JSON string for map[int]*StudentInfo
	BulletinBoard        string `gorm:"type:text"` // JSON string for []string
	WorkingSubs          string `gorm:"type:text"` // JSON string for []*Submission
	Submissions          string `gorm:"type:text"` // JSON string for map[int]*Submission
	WorkingHelpSubs      string `gorm:"type:text"` // JSON string for []*HelpSubmission
	HelpSubmissions      string `gorm:"type:text"` // JSON string for map[int]*HelpSubmission
	ActiveProblems       string `gorm:"type:text"` // JSON string for map[string]*ActiveProblem
	HelpEligibleStudents string `gorm:"type:text"` // JSON string for map[int]map[int]bool
	SeenHelpSubmissions  string `gorm:"type:text"` // JSON string for map[int]map[int]bool
	Snapshots            string `gorm:"type:text"` // JSON string for []*Snapshot
	StudentSnapshot      string `gorm:"type:text"` // JSON string for map[int]map[int]int
}

func StoreGlobalMaps() error {
	// Convert maps to JSON
	teacherMapJSON, err := json.Marshal(models.TeacherMap)
	if err != nil {
		return fmt.Errorf("failed to marshal TeacherMap: %w", err)
	}
	teacherPassJSON, err := json.Marshal(models.TeacherPass)
	if err != nil {
		return fmt.Errorf("failed to marshal TeacherPass: %w", err)
	}
	teacherNameToIdJSON, err := json.Marshal(models.TeacherNameToId)
	if err != nil {
		return fmt.Errorf("failed to marshal TeacherNameToId: %w", err)
	}
	teacherIdToNameJSON, err := json.Marshal(models.TeacherIdToName)
	if err != nil {
		return fmt.Errorf("failed to marshal TeacherIdToName: %w", err)
	}
	studentsJSON, err := json.Marshal(models.Students)
	if err != nil {
		return fmt.Errorf("failed to marshal Students: %w", err)
	}
	bulletinBoardJSON, err := json.Marshal(models.BulletinBoard)
	if err != nil {
		return fmt.Errorf("failed to marshal BulletinBoard: %w", err)
	}
	workingSubsJSON, err := json.Marshal(models.WorkingSubs)
	if err != nil {
		return fmt.Errorf("failed to marshal WorkingSubs: %w", err)
	}
	submissionsJSON, err := json.Marshal(models.Submissions)
	if err != nil {
		return fmt.Errorf("failed to marshal Submissions: %w", err)
	}
	workingHelpSubsJSON, err := json.Marshal(models.WorkingHelpSubs)
	if err != nil {
		return fmt.Errorf("failed to marshal WorkingHelpSubs: %w", err)
	}
	helpSubmissionsJSON, err := json.Marshal(models.HelpSubmissions)
	if err != nil {
		return fmt.Errorf("failed to marshal HelpSubmissions: %w", err)
	}
	activeProblemsJSON, err := json.Marshal(models.ActiveProblems)
	if err != nil {
		return fmt.Errorf("failed to marshal ActiveProblems: %w", err)
	}
	helpEligibleStudentsJSON, err := json.Marshal(models.HelpEligibleStudents)
	if err != nil {
		return fmt.Errorf("failed to marshal HelpEligibleStudents: %w", err)
	}
	seenHelpSubmissionsJSON, err := json.Marshal(models.SeenHelpSubmissions)
	if err != nil {
		return fmt.Errorf("failed to marshal SeenHelpSubmissions: %w", err)
	}
	snapshotsJSON, err := json.Marshal(models.Snapshots)
	if err != nil {
		return fmt.Errorf("failed to marshal Snapshots: %w", err)
	}
	studentSnapshotJSON, err := json.Marshal(models.StudentSnapshot)
	if err != nil {
		return fmt.Errorf("failed to marshal StudentSnapshot: %w", err)
	}

	// Insert or update the data in the database
	data := GlobalMap{
		ID:                   1, // Single row
		TeacherMap:           string(teacherMapJSON),
		TeacherPass:          string(teacherPassJSON),
		TeacherNameToId:      string(teacherNameToIdJSON),
		TeacherIdToName:      string(teacherIdToNameJSON),
		Students:             string(studentsJSON),
		BulletinBoard:        string(bulletinBoardJSON),
		WorkingSubs:          string(workingSubsJSON),
		Submissions:          string(submissionsJSON),
		WorkingHelpSubs:      string(workingHelpSubsJSON),
		HelpSubmissions:      string(helpSubmissionsJSON),
		ActiveProblems:       string(activeProblemsJSON),
		HelpEligibleStudents: string(helpEligibleStudentsJSON),
		SeenHelpSubmissions:  string(seenHelpSubmissionsJSON),
		Snapshots:            string(snapshotsJSON),
		StudentSnapshot:      string(studentSnapshotJSON),
	}
	return models.DB.Save(&data).Error
}

func ReloadGlobalMaps() error {
	var data GlobalMap

	// Try to fetch the first record
	err := models.DB.First(&data, 1).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No data found, skip processing
			return nil
		}
		return fmt.Errorf("failed to load teacher data: %w", err)
	}

	// Unmarshal JSON into maps
	err = json.Unmarshal([]byte(data.TeacherMap), &models.TeacherMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal TeacherMap: %w", err)
	}
	err = json.Unmarshal([]byte(data.TeacherPass), &models.TeacherPass)
	if err != nil {
		return fmt.Errorf("failed to unmarshal TeacherPass: %w", err)
	}
	err = json.Unmarshal([]byte(data.TeacherNameToId), &models.TeacherNameToId)
	if err != nil {
		return fmt.Errorf("failed to unmarshal TeacherNameToId: %w", err)
	}
	err = json.Unmarshal([]byte(data.TeacherIdToName), &models.TeacherIdToName)
	if err != nil {
		return fmt.Errorf("failed to unmarshal TeacherIdToName: %w", err)
	}
	err = json.Unmarshal([]byte(data.Students), &models.Students)
	if err != nil {
		return fmt.Errorf("failed to unmarshal Students: %w", err)
	}
	err = json.Unmarshal([]byte(data.BulletinBoard), &models.BulletinBoard)
	if err != nil {
		return fmt.Errorf("failed to unmarshal BulletinBoard: %w", err)
	}
	err = json.Unmarshal([]byte(data.Submissions), &models.Submissions)
	if err != nil {
		return fmt.Errorf("failed to unmarshal Submissions: %w", err)
	}
	err = json.Unmarshal([]byte(data.ActiveProblems), &models.ActiveProblems)
	if err != nil {
		return fmt.Errorf("failed to unmarshal ActiveProblems: %w", err)
	}
	err = json.Unmarshal([]byte(data.StudentSnapshot), &models.StudentSnapshot)
	if err != nil {
		return fmt.Errorf("failed to unmarshal StudentSnapshot: %w", err)
	}
	err = json.Unmarshal([]byte(data.WorkingSubs), &models.WorkingSubs)
	if err != nil {
		return fmt.Errorf("failed to unmarshal WorkingSubs: %w", err)
	}
	err = json.Unmarshal([]byte(data.WorkingHelpSubs), &models.WorkingHelpSubs)
	if err != nil {
		return fmt.Errorf("failed to unmarshal WorkingHelpSubs: %w", err)
	}
	err = json.Unmarshal([]byte(data.HelpSubmissions), &models.HelpSubmissions)
	if err != nil {
		return fmt.Errorf("failed to unmarshal HelpSubmissions: %w", err)
	}
	err = json.Unmarshal([]byte(data.HelpEligibleStudents), &models.HelpEligibleStudents)
	if err != nil {
		return fmt.Errorf("failed to unmarshal HelpEligibleStudents: %w", err)
	}
	err = json.Unmarshal([]byte(data.SeenHelpSubmissions), &models.SeenHelpSubmissions)
	if err != nil {
		return fmt.Errorf("failed to unmarshal SeenHelpSubmissions: %w", err)
	}
	err = json.Unmarshal([]byte(data.Snapshots), &models.Snapshots)
	if err != nil {
		return fmt.Errorf("failed to unmarshal Snapshots: %w", err)
	}

	return nil
}

// Author: Vinhthuy Phan, 2018
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"github.com/GPTA/src/restHandlers"

	_ "github.com/mattn/go-sqlite3"

	"io/ioutil"
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
	http.HandleFunc("/problem_dashboard", Authorize(restHandlers.ProblemDashboardHandler, ""))
	http.HandleFunc("/student_dashboard_feedback_provision", Authorize(restHandlers.StudentDashboardFeedbackProvisionHandler, ""))
	http.HandleFunc("/save_message_feedback", Authorize(restHandlers.MessageFeedbackHandler, ""))
	http.HandleFunc("/student_dashboard_submissions", Authorize(restHandlers.StudentDashboardSubmissionHandler, ""))
	http.HandleFunc("/has_message_feedback", Authorize(restHandlers.HasMessageBackFeedbackHandler, ""))
	http.HandleFunc("/teacher_signin_complete", restHandlers.TeacherSigninCompleteHandler)
	http.HandleFunc("/teacher_signin", restHandlers.TeacherSigninHandler)
	http.HandleFunc("/teacher_web_broadcast", Authorize(restHandlers.TeacherWebBroadcastHandler, "teacher"))
	http.HandleFunc("/student_dashboard_code_snapshot", Authorize(restHandlers.StudentDashboardCodeSpaceHandler, ""))
	http.HandleFunc("/teacher_exports_point", Authorize(restHandlers.ExportPointsHandler, "teacher"))
	http.HandleFunc("/", restHandlers.IndexHandler)
	http.HandleFunc("/index", restHandlers.IndexHandler)
	http.HandleFunc("/peer_tutoring", Authorize(restHandlers.PeerTutorHandler, "student"))
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
func inform_name_server() {
	nameserver := fmt.Sprintf("%s/tell?who=%s&address=%s", models.Config.NameServer, models.Config.CourseId, models.Config.Address)
	_, err := http.Get(nameserver)
	if err != nil {
		fmt.Println("Error", err)
		log.Fatal("Unable to contact with name server.")
	}
}

func get_course_specific_address(nameserver string, course string) {
	resp, err := http.Get(fmt.Sprintf("%s/ask?who=%s", nameserver, course))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	address := string(bodyBytes)
	fmt.Printf("* TeacherMap Login: %s/teacher_signin\n", address)
}

// -----------------------------------------------------------------
func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	rand.Seed(time.Now().UnixNano())
	config_file, teacher_file, student_file := "/Users/shashwatdadhich/go/src/github.com/GPTA/Examples/gem_config.json", "/Users/shashwatdadhich/go/src/github.com/GPTA/Examples/teachers.txt", "/Users/shashwatdadhich/go/src/github.com/GPTA/Examples/students.txt"
	flag.StringVar(&config_file, "c", config_file, "json-formatted configuration file.")
	flag.StringVar(&teacher_file, "add_teachers", teacher_file, "teacher file.")
	flag.StringVar(&student_file, "add_students", student_file, "student file.")
	flag.Parse()
	if config_file == "" {
		flag.Usage()
		os.Exit(1)
	}
	models.Config = init_config(config_file)
	if models.Config.NameServer != "" {
		inform_name_server()
	}
	repository.InitDatabase(models.Config.Database, models.Config.DBUserName, models.Config.DBPassWord, models.Config.DBServerIP)
	if teacher_file != "" {
		restHandlers.AddMultiple(teacher_file, "teacher")
	}
	if student_file != "" {
		restHandlers.AddMultiple(student_file, "student")
	}
	init_handlers()
	repository.LoadTeachers()
	fmt.Println("**************************************************")

	fmt.Printf("*   Course id:      %s\n", models.Config.CourseId)
	if models.Config.NameServer != "" {
		fmt.Printf("*   Server address: %s\n", models.Config.NameServer)
	} else {
		fmt.Printf("*   Serving at:     %s\n", models.Config.Address)
	}
	fmt.Printf("*   GEM %s\n", VERSION)
	fmt.Println("**************************************************\n")
	get_course_specific_address(models.Config.NameServer, models.Config.CourseId)
	err := http.ListenAndServe(models.Config.Address, nil)
	if err != nil {
		log.Fatal("Unable to serve gem server at " + models.Config.Address)
	}
}

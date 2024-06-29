// Author: Vinhthuy Phan, 2018
package main

import (
	"encoding/json"
	"flag"

	_ "github.com/mattn/go-sqlite3"

	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"
)

//-----------------------------------------------------------------
func initHandlers(app App) {
	http.HandleFunc("/test", testHandler)

	// Analytics
	// http.HandleFunc("/learning_report", learning_reportHandler)
	// http.HandleFunc("/analyze_submissions", analyze_submissionsHandler)
	// http.HandleFunc("/view_activities", view_activitiesHandler)
	// http.HandleFunc("/report", reportHandler)
	// http.HandleFunc("/report_tag", report_tagHandler)
	// http.HandleFunc("/view_answers", view_answersHandler)
	// http.HandleFunc("/statistics", statisticsHandler)

	// // Others
	// http.HandleFunc("/student_periodic_update", Authorize(student_periodic_updateHandler, "student"))

	// http.HandleFunc("/student_gets_report", Authorize(student_gets_reportHandler, "student"))
	// http.HandleFunc("/student_checks_in", Authorize(student_checks_inHandler, "student"))
	// http.HandleFunc("/student_shares", Authorize(student_sharesHandler, "student"))
	// http.HandleFunc("/student_gets", Authorize(student_getsHandler, "student"))
	// http.HandleFunc("/view_bulletin_board", view_bulletin_boardHandler)
	// http.HandleFunc("/remove_bulletin_page", remove_bulletin_pageHandler)
	// http.HandleFunc("/bulletin_board_data", bulletin_board_dataHandler)
	// http.HandleFunc("/complete_registration", complete_registrationHandler)

	// http.HandleFunc("/student_ask_help", Authorize(studentAskHelpHandler, "student"))
	// http.HandleFunc("/student_get_help_code", Authorize(studentGetHelpCode, "student"))
	// http.HandleFunc("/student_return_without_feedback", Authorize(student_return_without_feedbackHandler, "student"))
	// http.HandleFunc("/student_send_help_message", Authorize(student_send_help_messageHandler, "student"))
	// http.HandleFunc("/student_send_thank_you", Authorize(sendThankYouHandler, "student"))

	// http.HandleFunc("/teacher_get_help_code", Authorize(teacherGetHelpCode, "teacher"))
	// http.HandleFunc("/teacher_return_without_feedback", Authorize(teacher_return_without_feedbackHandler, "teacher"))
	// http.HandleFunc("/teacher_send_help_message", Authorize(teacher_send_help_messageHandler, "teacher"))

	// http.HandleFunc("/teacher_gets_queue", Authorize(teacher_gets_queueHandler, "teacher"))
	// http.HandleFunc("/teacher_adds_bulletin_page", Authorize(teacher_adds_bulletin_pageHandler, "teacher"))
	// http.HandleFunc("/teacher_clears_submissions", Authorize(teacher_clears_submissionsHandler, "teacher"))
	// http.HandleFunc("/teacher_deactivates_problems", Authorize(teacher_deactivates_problemsHandler, "teacher"))
	// http.HandleFunc("/teacher_grades", Authorize(teacher_gradesHandler, "teacher"))
	// http.HandleFunc("/teacher_puts_back", Authorize(teacher_puts_backHandler, "teacher"))
	// http.HandleFunc("/teacher_gets", Authorize(teacher_getsHandler, "teacher"))
	// http.HandleFunc("/teacher_broadcasts", Authorize(teacher_broadcastsHandler, "teacher"))
	// http.HandleFunc("/teacher_gets_passcode", Authorize(teacher_gets_passcodeHandler, "teacher"))
	// http.HandleFunc("/student_gets_passcode", Authorize(student_gets_passcodeHandler, "student"))

	// http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "pong") })

	// http.HandleFunc("/get_testcase", Authorize(testcase_getsHandler, ""))

	// http.HandleFunc("/code_snapshot", Authorize(codeSnapshotHandler, ""))
	// http.HandleFunc("/get_global_info", Authorize(globalInfoHandler, ""))

	// // http.HandleFunc("/get_codespace", Authorize(codespaceHandler))
	// // http.HandleFunc("/get_snapshot", Authorize(getCodeSnapshotHandler))
	// http.HandleFunc("/save_snapshot_feedback", Authorize(codeSnapshotFeedbackHandler, ""))
	// http.HandleFunc("/get_snapshot_feedback", Authorize(getSnapshotFeedbackHandler, ""))
	// http.HandleFunc("/save_snapshot_back_feedback", Authorize(studentSendBackFeedbackHandler, ""))
	// // http.HandleFunc("/student_views_feedback", Authorize(studentViewsFeedbackHandler))
	// // http.HandleFunc("/teacher_views_feedback", Authorize(teacherViewsFeedbackHandler))

	// // http.HandleFunc("/help_requests", Authorize(helpRequestListHandler))
	// // http.HandleFunc("/view_help_request", Authorize(viewHelpRequestHandler))
	// http.HandleFunc("/set_peer_tutor", Authorize(setPeerTutorHandler, "teacher"))

	// http.HandleFunc("/view_exercises", Authorize(problemListHandler, ""))
	// http.HandleFunc("/problem_dashboard", Authorize(problemDashboardHandler, ""))
	// http.HandleFunc("/student_dashboard_feedback_provision", Authorize(studentDashboardFeedbackProvisionHandler, ""))
	// http.HandleFunc("/save_message_feedback", Authorize(messageFeedbackHandler, ""))
	// http.HandleFunc("/student_dashboard_submissions", Authorize(studentDashboardSubmissionHandler, ""))
	// http.HandleFunc("/has_message_feedback", Authorize(hasMessageBackFeedbackHandler, ""))
	http.HandleFunc("/signin_complete", app.signinCompleteHandler)
	http.HandleFunc("/signin", app.signinHandler)
	// http.HandleFunc("/teacher_web_broadcast", Authorize(teacherWebBroadcastHandler, "teacher"))
	// http.HandleFunc("/student_dashboard_code_snapshot", Authorize(studentDashboardCodeSpaceHandler, ""))
	// http.HandleFunc("/teacher_exports_point", Authorize(exportPointsHandler, "teacher"))
	http.HandleFunc("/", app.indexHandler)
	http.HandleFunc("/index", app.indexHandler)
	// http.HandleFunc("/peer_tutoring", Authorize(peerTutorHandler, "student"))
}

//-----------------------------------------------------------------
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

//-----------------------------------------------------------------
func initConfig(filename string) *Configuration {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	config := &Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	rand.Seed(time.Now().UnixNano())
	configFile := ""
	flag.StringVar(&configFile, "c", configFile, "json-formatted configuration file.")
	flag.Parse()
	if configFile == "" {
		flag.Usage()
		os.Exit(1)
	}
	Config = initConfig(configFile)
	app := App{}
	app.migrateSchema(*Config)
	initHandlers(app)

	err := http.ListenAndServe(Config.Address, nil)
	if err != nil {
		log.Fatal("Unable to serve gem server at " + Config.Address)
	}
}

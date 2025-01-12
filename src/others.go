// Author: Vinhthuy Phan, 2018
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// -----------------------------------------------------------------------------------
func teacher_gets_passcodeHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	fmt.Fprintf(w, Passcode)
}

func student_gets_passcodeHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	fmt.Fprintf(w, Passcode)
}

// -----------------------------------------------------------------------------------
func testHandler(w http.ResponseWriter, r *http.Request) {
	// Show content of boards
	fmt.Println("Students:", len(Students))
	for _, st := range Students {
		// fmt.Printf("Uid: %d has %d pages. Status: %d\n", uid, len(st.Boards), st.SubmissionStatus)
		for i := 0; i < len(st.Boards); i++ {
			b := st.Boards[i]
			fmt.Printf("Attempts: %d, Filename: %s, Pid: %d, Answer: %s, len of content: %d\n",
				b.Attempts, b.Filename, b.Pid, b.Answer, len(b.Content))
		}
	}

	fmt.Printf("WorkingSubs: %d entries", len(WorkingSubs))
	for i := 0; i < len(WorkingSubs); i++ {
		fmt.Println(WorkingSubs[i].Sid, WorkingSubs[i].Uid, WorkingSubs[i].Pid, WorkingSubs[i].Priority)
		fmt.Println(WorkingSubs[i].Content)
	}
	fmt.Println()

	fmt.Println("ActiveProblems:", ActiveProblems)
	for fname, v := range ActiveProblems {
		fmt.Println(fname, v.Active, "Answers:", v.Answers, "Attempts:", v.Attempts)
		fmt.Println(fname, v.Info.Pid, v.Info.Merit, v.Info.Effort, v.Info.Attempts, v.Info.ExactAnswer, v.Info.Answer)
	}
	fmt.Fprintf(w, Passcode)
}

//-----------------------------------------------------------------------------------

func testcase_getsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	filename := r.FormValue("file_name")

	// Get the problem ID by filename
	var problem Problem
	err := Database.Where("filename = ?", filename).First(&problem).Error
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Problem not found", http.StatusNotFound)
		return
	}

	// Get the test cases for the problem ID
	var testCases []TestCase
	err = Database.Where("problem_id = ?", problem.ID).Find(&testCases).Error
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Test cases not found", http.StatusNotFound)
		return
	}

	// Concatenate all the test cases
	var testCasesStr string
	for _, testCase := range testCases {
		if testCase.TestCases != "" {
			testCasesStr += testCase.TestCases[1 : len(testCase.TestCases)-1] // Remove surrounding quotes if any
		}
	}

	// Respond with the concatenated test cases in JSON format
	fmt.Fprintf(w, "[%s]", testCasesStr)
}

func logEvent(eventName string, userID int, userType, eventType, otherInfo string) {
	err := AddUserEventLog(userID, eventName, userType, eventType, otherInfo, time.Now())
	if err != nil {
		log.Println("Error logging event:", err)
	}
}

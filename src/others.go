// Author: Vinhthuy Phan, 2018
package main

import (
	"fmt"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"log"
	"net/http"
	"time"
)

// -----------------------------------------------------------------------------------
func teacher_gets_passcodeHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	fmt.Fprintf(w, models.Passcode)
}

func student_gets_passcodeHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	fmt.Fprintf(w, models.Passcode)
}

// -----------------------------------------------------------------------------------
func testHandler(w http.ResponseWriter, r *http.Request) {
	// Show content of boards
	fmt.Println("Students:", len(models.Students))
	for _, st := range models.Students {
		// fmt.Printf("Uid: %d has %d pages. Status: %d\n", uid, len(st.Boards), st.SubmissionStatus)
		for i := 0; i < len(st.Boards); i++ {
			b := st.Boards[i]
			fmt.Printf("Attempts: %d, Filename: %s, Pid: %d, Answer: %s, len of content: %d\n",
				b.Attempts, b.Filename, b.Pid, b.Answer, len(b.Content))
		}
	}

	fmt.Printf("WorkingSubs: %d entries", len(models.WorkingSubs))
	for i := 0; i < len(models.WorkingSubs); i++ {
		fmt.Println(models.WorkingSubs[i].Sid, models.WorkingSubs[i].Uid, models.WorkingSubs[i].Pid, models.WorkingSubs[i].Priority)
		fmt.Println(models.WorkingSubs[i].Content)
	}
	fmt.Println()

	fmt.Println("ActiveProblems:", models.ActiveProblems)
	for fname, v := range models.ActiveProblems {
		fmt.Println(fname, v.Active, "Answers:", v.Answers, "Attempts:", v.Attempts)
		fmt.Println(fname, v.Info.Pid, v.Info.Merit, v.Info.Effort, v.Info.Attempts, v.Info.ExactAnswer, v.Info.Answer)
	}
	fmt.Fprintf(w, models.Passcode)
}

//-----------------------------------------------------------------------------------

func testcase_getsHandler(w http.ResponseWriter, r *http.Request, who string, uid int) {
	filename := r.FormValue("file_name")

	// Get Problem ID
	problemID, err := repository.GetProblemIDByFilename(filename)
	if err != nil {
		log.Fatalf("Error fetching problem ID: %v", err)
	}
	if problemID == 0 {
		http.Error(w, "Problem not found", http.StatusNotFound)
		return
	}

	// Get Test Cases
	testCases, err := repository.GetTestCasesByProblemID(problemID)
	if err != nil {
		log.Fatalf("Error fetching test cases: %v", err)
	}

	// Format Test Cases
	var formattedTestCases string
	for _, tc := range testCases {
		if tc != "" {
			formattedTestCases += tc[1 : len(tc)-1]
		}
	}

	// Respond with Test Cases
	fmt.Fprintf(w, "["+formattedTestCases+"]")
}

func logEvent(eventName string, userID int, userType, eventType, otherInfo string) {
	_, _ = repository.AddUserEventLog(eventName, userID, userType, eventType, otherInfo, time.Now())
}

package openAI

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/GPTA/src/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

//func InstructionsWithExampleHandler(w http.ResponseWriter, r *http.Request) {
//	c, _ := gin.CreateTestContext(w)
//	c.Request = r
//	EffectiveFeedbackInstructionWithExample(c)
//}

func ProcessCodeWithPromptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ProcessCodeWithPromptC(w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func SummarizeClassPerformance(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ProcessStudentCodeSubmissions(w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

//func EffectiveFeedbackInstructionWithExample(c *gin.Context) {
//	processHandlers(c, getFeedbackWithInstructionsWithExample)
//}

//func processHandlers(c *gin.Context, f getPromptFunc) {
//	var requestBody RequestBody
//	if err := c.BindJSON(&requestBody); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	requestData := createRequestData(f(requestBody.Course, requestBody.Duration), requestBody)
//	makeRequest(c, requestData)
//}

//
//func ProcessCodeWithPrompt(w http.ResponseWriter, r *http.Request) {
//	var requestData struct {
//		Code   string `json:"code"`
//		UserID int    `json:"user_id"`
//	}
//
//	// Parse the incoming JSON body
//	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	NewInstructions := "I am a student learning how to code in Python. Here is my solution to a problem. The specification of the problem is in the comment of the problem. I want to improve my programming skills. Please:\n\n1. First explain what parts of my code show good programming practices\n2. Then identify 2-3 specific areas where I could improve\n3. For each area of improvement, explain why it matters and show me a small example of better code\n4. Give me a similar but slightly different programming exercise to practice these concepts\n5. Don't use any Markdown formatting like ** for bold or ### for headings. Reply in plain text only, without any special formatting. Give the reply in a clear and understandable format."
//
//	prompt := fmt.Sprintf(
//		"Instructions:  %s\nCode:\n%s\n",
//		NewInstructions, requestData.Code,
//	)
//
//	// Create OpenAI request message
//	messages := []goopenai.Message{
//		{Role: "user", Content: prompt},
//	}
//
//	// Create a Gin context from http.ResponseWriter and http.Request
//	c, _ := gin.CreateTestContext(w)
//	c.Request = r
//
//	// Call the existing makeRequest function
//	makeRequest2(c, messages)
//}

func ProcessCodeWithPromptC(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Code   string `json:"code"`
		UserID int    `json:"user_id"`
	}

	// Parse JSON request body
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	NewInstructions := "I am a student learning how to code in Python. Here is my solution to a problem. The specification of the problem is in the comment of the problem. I want to improve my programming skills. Please:\n\n1. First explain what parts of my code show good programming practices\n2. Then identify 2-3 specific areas where I could improve\n3. For each area of improvement, explain why it matters and show me a small example of better code\n4. Give me a similar but slightly different programming exercise to practice these concepts\n5. Don't use any Markdown formatting like ** for bold or ### for headings. Reply in plain text only, without any special formatting. Give the reply in a clear and understandable format."

	prompt := fmt.Sprintf(
		"Instructions:  %s\nCode:\n%s\n",
		NewInstructions, requestData.Code,
	)

	// Create Claude request
	messages := []map[string]string{
		{"role": "user", "content": prompt},
	}

	// Create a Gin context from http.ResponseWriter and http.Request
	c, _ := gin.CreateTestContext(w)
	c.Request = r

	// Call Claude API request function
	makeRequestClaude(c, messages)
}

func ProcessStudentCodeSubmissions(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		ProblemID string `json:"problem_id"`
	}

	//Parse JSON request body
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert ProblemID from string to int
	problemID, err := strconv.Atoi(requestData.ProblemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid problem_id format: %v", err), http.StatusBadRequest)
		return
	}

	problemDescription, err := repository.GetProblemDescription(problemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch problem description: %v", err), http.StatusInternalServerError)
		return
	}

	// Get the latest code snapshots from the database using the problem_id
	latestCodeSnapshots, err := repository.GetLatestCodeSnapshots(problemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch latest code snapshots: %v", err), http.StatusInternalServerError)
		return
	}

	// Construct a detailed prompt for error analysis
	NewInstructions := "You are given multiple code submissions from students. Your task is to:\n\n" +
		"1. Identify recurring errors, misconceptions, and logic mistakes across these submissions.\n" +
		"2. Explain why these mistakes occur and how they can be corrected.\n" +
		"3. Suggest teaching strategies or lesson improvements to address the underlying misconceptions.\n" +
		"4. Categorize similar mistakes together and provide concise, clear explanations.\n" +
		"5. Focus on syntax, logical, and conceptual errors.\n" +
		"6. Conclude with a 'Key Takeaways' section summarizing the most important issues and recommendations.\n\n"

	// Combine all code snapshots into the prompt
	codeData := ""
	for _, snapshot := range latestCodeSnapshots {
		codeData += fmt.Sprintf("Student ID: %d\nTimestamp: %s\nCode:\n%s\n\n", snapshot.StudentID, snapshot.LastUpdatedAt, snapshot.Code)
	}

	// Build the full prompt
	prompt := fmt.Sprintf(
		"Problem Description:\n%s\n\nInstructions: %s\nCode Submissions:\n%s",
		problemDescription, NewInstructions, codeData,
	)

	// Create Claude request
	messages := []map[string]string{
		{"role": "user", "content": prompt},
	}

	// Create a Gin context from http.ResponseWriter and http.Request
	c, _ := gin.CreateTestContext(w)
	c.Request = r

	// Call Claude API request function
	makeRequestClaude2(c, messages, problemID)
}

func GetClassFeedback(w http.ResponseWriter, r *http.Request) {
	// Parse the problem_id from the request body
	var requestData struct {
		ProblemID string `json:"problem_id"`
	}

	// Decode the request JSON
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	problemID, err := strconv.Atoi(requestData.ProblemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid problem_id format: %v", err), http.StatusBadRequest)
		return
	}

	// Fetch the class feedback from the database
	feedback, err := repository.GetClassFeedback(problemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching class feedback: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare the response
	response := struct {
		Feedback string `json:"feedback"`
	}{
		Feedback: feedback,
	}

	// Set the response content type and encode the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

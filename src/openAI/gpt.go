package openAI

import (
	"encoding/json"
	"fmt"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
	"time"
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

func SummarizeScaffolding(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ProcessScaffolding(w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func SummarizeStudentProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ProcessStudentProgress(w, r)
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

	// Convert FeedbackID from string to int
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

func GetLatestFeedbackByProblemID(w http.ResponseWriter, r *http.Request) {
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
	feedbacks, err := repository.GetClassFeedbackByProblemID(problemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching class feedback: %v", err), http.StatusInternalServerError)
		return
	}

	// Sort feedbacks by FeedbackTime in descending order (latest first)
	sort.Slice(feedbacks, func(i, j int) bool {
		return feedbacks[i].FeedbackTime.After(feedbacks[j].FeedbackTime)
	})

	// Pick the latest feedback (first element after sorting)
	var latestFeedback string
	if len(feedbacks) > 0 {
		latestFeedback = feedbacks[0].Feedback
	}

	response := struct {
		Feedback string `json:"feedback"`
	}{
		Feedback: latestFeedback,
	}

	// Set the response content type and encode the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

func ListFeedbackHistoryByProblemID(w http.ResponseWriter, r *http.Request) {
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
	feedbacks, err := repository.GetClassFeedbackByProblemID(problemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching class feedback: %v", err), http.StatusInternalServerError)
		return
	}

	type FeedbackResponse struct {
		Feedbacks []struct {
			FeedbackID   int       `json:"feedback_id"`
			FeedbackTime time.Time `json:"feedback_time"`
		} `json:"feedbacks"`
	}

	response := FeedbackResponse{
		Feedbacks: make([]struct {
			FeedbackID   int       `json:"feedback_id"`
			FeedbackTime time.Time `json:"feedback_time"`
		}, 0, len(feedbacks)),
	}

	for _, fb := range feedbacks {
		response.Feedbacks = append(response.Feedbacks, struct {
			FeedbackID   int       `json:"feedback_id"`
			FeedbackTime time.Time `json:"feedback_time"`
		}{
			FeedbackID:   fb.ID,
			FeedbackTime: fb.FeedbackTime,
		})
	}

	// Set the response content type and encode the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

func GetFeedbackByFeedbackID(w http.ResponseWriter, r *http.Request) {
	// Parse the problem_id from the request body
	var requestData struct {
		FeedbackID string `json:"feedback_id"`
	}

	// Decode the request JSON
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	feedbackID, err := strconv.Atoi(requestData.FeedbackID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid problem_id format: %v", err), http.StatusBadRequest)
		return
	}

	// Fetch the class feedback from the database
	feedbacks, err := repository.GetClassFeedbackByFeedbackID(feedbackID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching class feedback: %v", err), http.StatusInternalServerError)
		return
	}

	// Sort feedbacks by FeedbackTime in descending order (latest first)
	sort.Slice(feedbacks, func(i, j int) bool {
		return feedbacks[i].FeedbackTime.After(feedbacks[j].FeedbackTime)
	})

	// Pick the latest feedback (first element after sorting)
	var latestFeedback string
	if len(feedbacks) > 0 {
		latestFeedback = feedbacks[0].Feedback
	}

	response := struct {
		Feedback string `json:"feedback"`
	}{
		Feedback: latestFeedback,
	}

	// Set the response content type and encode the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

func ProcessStudentProgress(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		ProblemID string `json:"problem_id"`
	}

	//Parse JSON request body
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert FeedbackID from string to int
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

	NewInstructions := "NewInstructions := `You are given multiple code submissions from students. Your task is to:\n1. Analyze each student's submission and compare it with the correct solution.\n2. Provide a one-line explanation of the student's feedback.\n3. Calculate a percentage (0-100) indicating how close their submission is to the correct solution.\n4. Return the response strictly in the following JSON format:\n[\n  {\n    \"student_id\": 1,\n    \"explanation\": \"one line explanation for the student feedback\",\n    \"percentage\": 25\n  },\n  {\n    \"student_id\": 2,\n    \"explanation\": \"one line explanation for the student feedback\",\n    \"percentage\": 85\n  }\n]\nOnly return valid JSON. Do not include any additional text or formatting outside the JSON structure. dont even give messages like Here is the JSON response with the analysis of each student's submission:. I want just the array of json`\n\n"
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
	makeRequestClaude3(c, messages, problemID)
}

func ProcessScaffolding(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		ProblemID           string `json:"problem_id"`
		ScaffoldingStrategy string `json:"scaffolding_strategy"`
	}

	//Parse JSON request body
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert FeedbackID from string to int
	problemID, err := strconv.Atoi(requestData.ProblemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid problem_id format: %v", err), http.StatusBadRequest)
		return
	}

	scaffoldingID, err := strconv.Atoi(requestData.ScaffoldingStrategy)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid problem_id format: %v", err), http.StatusBadRequest)
		return
	}

	NewInstructions := ""

	switch scaffoldingID {
	case models.ScaffoldingStrategyFillInTheBlanks:
		NewInstructions = `Generate 3 variations of this problem using the fill-in-the-blanks scaffolding technique, tailored to different student skill levels: 
(struggling, developing, nearly_proficient).

Each variation must be a complete, self-contained problem description that students can understand without additional context.

For each variation, include the following:
1. A clear problem statement with an example of input and output.
2. A brief "Purpose" explaining why this skill matters and its real-world application.
3. The code with strategically placed blanks (e.g., ___1___, ___2___, etc.). The blanks should correspond to missing code sections that the student needs to fill in.
   - Start with a highly scaffolded version with small blanks and a word bank for support.
   - Gradually reduce support in later versions by removing the word bank and providing entire missing blocks of code.
4. Label each blank (e.g., ___1___) and provide appropriate hints to guide students.
5. Include assessment criteria to evaluate students' work, such as correctness, completeness, and clarity.
6. Add 1-2 reflection questions that encourage students to think about the reasoning behind their choices.
7. Clearly label each version with its corresponding skill level (struggling, developing, nearly_proficient) and provide an estimated time for completion.

Return the response strictly in **JSON format** as shown below:

{
  "struggling": "The generated scaffolding for struggling students",
  "developing": "The generated scaffolding for developing students",
  "nearly_proficient": "The generated scaffolding for nearly proficient students"
}

Ensure that the JSON is **valid**:
- Only include valid key-value pairs for "struggling", "developing", and "nearly_proficient".
- Do not include any additional text, formatting, or explanations outside of the JSON structure.
- The JSON keys must exactly match "strugglin", "developing", and "nearly_proficient".
- The output should be **pure JSON** with no extra characters, sentences, or formatting.
- dont even write "Here is my attempt at the scaffolded solutions in valid JSON format:" just give the json object`

	case models.ScaffoldingStrategyStepByStepTasks:
		NewInstructions = `Rewrite the problem as 3 progressively scaffolded task lists, each representing a different level of support (struggling, developing, nearly_proficient).

Each variation must be a complete, self-contained problem description with clear examples of inputs/outputs. 

For each level of support, provide the following:
1. A clear problem description with an example of input and output.
2. A brief "Learning Objectives" section explaining what skills students will practice.
3. A list of tasks with specific steps:
   - Early versions should break the problem into many explicit steps.
   - Later versions should use higher-level tasks requiring students to determine implementation details.

Each version should include:
- Clear success criteria.
- 2-3 checkpoint questions to assess understanding.
- A teacher note identifying potential stumbling points.
- Estimated completion time.

Label each version with the skill level (struggling, developing, nearly_proficient) and ensure the JSON keys match the exact skill levels.

Return the response strictly in **JSON format** as shown below:
{
  "struggling": "The generated scaffolding for struggling students",
  "developing": "The generated scaffolding for developing students",
  "nearly_proficient": "The generated scaffolding for nearly proficient students"
}

Ensure that the JSON is **valid**:
- Only include valid key-value pairs.
- Do not include any additional text, formatting, or explanations outside of the JSON structure.
- The keys must exactly match "struggling", "developing", and "nearly_proficient".
- Do not include any other messages or instructions outside the JSON structure.
- dont even write "Here is my attempt at the scaffolded solutions in valid JSON format:" just give the json object

The output must be **pure JSON** with no extra characters, sentences, or formatting. It should only contain the valid JSON object shown above, with the generated scaffolding in the corresponding values.`

	case models.ScaffoldingStrategyGuidedCodeWithHints:
		NewInstructions = `Create 3 scaffolded versions of a guided code solution using inline comments and hints: each representing a different level of support (struggling, developing, nearly_proficient).

For each variation, include the following:

1. A clear problem description with an example of input/output.
2. A brief explanation of real-world applications of the problem.
3. The code, including inline comments and hints:
   - High-support version: Provide nearly complete code with detailed comments and hints to guide students through the problem.
   - Medium-support version: Leave key logic missing, with guiding comments to help students figure out the next steps.
   - Low-support version: Provide only the skeleton structure of the code with minimal guidance, relying on students to fill in most of the logic.
   
4. Mark each student task with "# TODO" comments, clearly indicating where they need to fill in the missing parts.

For each version, include:
1. A simple assessment rubric to evaluate students' work.
2. 1-2 extension challenges that will encourage further exploration and deeper understanding.
3. A debugging checklist to help students identify common mistakes.

Clearly label each version with its difficulty level (struggling, developing, nearly_proficient) and provide an estimated time for completion.

Return the response strictly in **JSON format** as shown below:

{
  "struggling": "The generated scaffolding for struggling students",
  "developing": "The generated scaffolding for developing students",
  "nearly_proficient": "The generated scaffolding for nearly proficient students"
}

Ensure that the JSON is **valid**:
- Only include valid key-value pairs for "struggling", "developing", and "nearly_proficient".
- Do not include any additional text, formatting, or explanations outside of the JSON structure.
- The JSON keys must exactly match "struggling", "developing", and "nearly_proficient".
- The output should be **pure JSON** with no extra characters, sentences, or formatting.
- dont even write "Here is my attempt at the scaffolded solutions in valid JSON format:" just give the json object
`
	case models.ScaffoldingStrategyDebugThisCode:
		NewInstructions = `Write 3 variations of a buggy version of this problem for developing debugging skills: each representing a different level of support (struggling, developing, nearly_proficient).

For each variation, include the following:

1. A **clear problem description** explaining what the program should do when working correctly, with **specific input/output examples**. 
2. Provide the **buggy code** for each variation, along with:
   - A description of the **expected behavior** of the program.
   - A description of the **actual behavior** of the program with the bugs present.
   
3. The bugs should vary in difficulty:
   - High-support version: Contain **obvious syntax or semantic errors**, with **clear hints** to guide the students toward the solution.
   - Medium-support version: Contain **subtler logic errors**, with **fewer clues** to encourage more independent debugging.
   - Low-support version: Contain **complex issues**, requiring students to troubleshoot without significant hints or guidance.

For each variation, include:
1. **Bug severity ratings** to help students understand the importance of each issue.
2. **2-3 specific debugging strategies** that the students can practice (e.g., using print statements, checking for edge cases, reviewing variable types).
3. A **post-debugging reflection template** where students can note what strategies worked, what they learned, and how they would approach similar bugs in the future.

Clearly label each version with its **difficulty level** (struggling, developing, nearly_proficient) and provide an **estimated debugging time**.

Return the response strictly in **JSON format** as shown below:

{
  "struggling": "The generated scaffolding for struggling students",
  "developing": "The generated scaffolding for developing students",
  "nearly_proficient": "The generated scaffolding for nearly proficient students"
}

Ensure that the JSON is **valid**:
- Only include valid key-value pairs for "struggling", "developing", and "nearly_proficient".
- Do not include any additional text, formatting, or explanations outside of the JSON structure.
- The JSON keys must exactly match "struggling", "developing", and "nearly_proficient".
- The output should be **pure JSON** with no extra characters, sentences, or formatting.
- dont even write "Here is my attempt at the scaffolded solutions in valid JSON format:" just give the json object
`
	case models.ScaffoldingStrategyIncrementalFeatureImpl:
		NewInstructions = `Break down this problem into 4-5 incremental feature builds, progressing from basic to complete: each representing a different level of support (struggling, developing, nearly_proficient).

For each version, include the following:

1. A **complete problem description** that progressively builds toward the same end goal.
2. Start with a **simple core feature implementation** for the struggling version. Each subsequent version will add one new feature or increase complexity.
3. The **final version** should match the **complete original requirements**, demonstrating a fully implemented solution.

For each version:
1. **Working code** from the previous version, with all implemented features up to that point.
2. **Clear specification** of what is new in this version (e.g., adding a new feature, improving performance, handling edge cases).
3. **Implementation guidance**, providing step-by-step instructions or hints about how to proceed.
4. **Test cases** to ensure the new features work as expected.
5. **Refactoring suggestions** for improving code structure, performance, or readability.

Label each version with:
- **Difficulty level** (struggling, developing, nearly_proficient).
- **Estimated implementation time** for that version.

Ensure that each version is a complete, self-contained description that enables students to build progressively towards a fully functional implementation. This approach is best for multi-session projects where students can receive feedback before advancing to the next version.

Return the response strictly in **JSON format** as shown below:

{
  "struggling": "The generated scaffolding for struggling students",
  "developing": "The generated scaffolding for developing students",
  "nearly_proficient": "The generated scaffolding for nearly proficient students"
}

Ensure that the JSON is **valid**:
- Only include valid key-value pairs for "struggling", "developing", and "nearly_proficient".
- Do not include any additional text, formatting, or explanations outside of the JSON structure.
- The JSON keys must exactly match "struggling", "developing", and "nearly_proficient".
- The output should be **pure JSON** with no extra characters, sentences, or formatting.
- dont even write "Here is my attempt at the scaffolded solutions in valid JSON format:" just give the json object
`

	default:
		NewInstructions = "" // No scaffolding strategy found, fallback case
	}

	problemDescription, err := repository.GetProblemDescription(problemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch problem description: %v", err), http.StatusInternalServerError)
		return
	}

	// Build the full prompt
	prompt := fmt.Sprintf(
		"Problem Description:\n%s\n\nInstructions: %s\n",
		problemDescription, NewInstructions,
	)

	// Create Claude request
	messages := []map[string]string{
		{"role": "user", "content": prompt},
	}

	// Create a Gin context from http.ResponseWriter and http.Request
	c, _ := gin.CreateTestContext(w)
	c.Request = r

	// Call Claude API request function
	makeRequestClaude4(c, messages, problemID, scaffoldingID)
}

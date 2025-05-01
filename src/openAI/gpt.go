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
	"strings"
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

func ProcessStudentCodeSubmissions(w http.ResponseWriter, r *http.Request) string {

	problemIDStr := r.FormValue("problem_id")
	if problemIDStr == "" {
		http.Error(w, "Missing problem_id", http.StatusBadRequest)
		return ""
	}

	// Convert string to int
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		http.Error(w, "Invalid problem_id", http.StatusBadRequest)
		return ""
	}
	newStr := r.FormValue("new")
	isNew := false
	if newStr == "true" {
		isNew = true
	}

	if !isNew {
		summary, _ := repository.GetLatestClassFeedbackByProblemID(problemID)
		if summary != nil {
			return summary.Feedback
		}
	}

	problemDescription, err := repository.GetProblemDescription(problemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch problem description: %v", err), http.StatusInternalServerError)
		return ""
	}

	// Get the latest code snapshots from the database using the problem_id
	latestCodeSnapshots, err := repository.GetLatestCodeSnapshots(problemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch latest code snapshots: %v", err), http.StatusInternalServerError)
		return ""
	}

	// Construct a detailed prompt for error analysis
	NewInstructions := `
dont give Here is the analysis of the student code submissions for the word game score calculator problem:
json...just give me the whole json object as response 
# Student Code Submission Analysis Prompt

You are an expert in analyzing student code submissions for introductory computer science exercises (CS1/CS2 level). Given a collection of student code submissions for a specific programming problem and their corresponding assessment results, your task is to generate a structured summary that provides instructors with actionable insights into common student difficulties.

Your analysis should follow a staged approach, where each stage builds upon insights gained in previous stages. This will ensure a systematic, consistent, and comprehensive analysis of student performance patterns.

## Staged Analysis Process

### Stage 1: Data Processing and Performance Classification
First, process the raw submission data:
1. Calculate the total number of submissions.
2. Classify each submission into one of these performance levels using these specific criteria:
   - **Poor (Failing with major issues)**:
     * Submissions that fail most or all test cases
     * Code contains 3+ fundamental errors from different error categories
     * Code structure shows significant misunderstanding of core concepts
     * Code may not compile/run at all
   - **Struggling (Failing but showing potential)**:
     * Submissions that pass some test cases but fail others
     * Code contains 1-2 fundamental errors, but shows understanding of basic concepts
     * Core algorithm demonstrates partial correctness
     * Code runs but produces incorrect results for certain inputs
   - **Good Progress (Passing with minor mistakes)**:
     * Submissions that pass most test cases
     * Code has working core functionality with minor issues
     * Errors are limited to edge cases or efficiency concerns
     * Code produces correct results for most inputs but may have minor issues
   - **Strong (Fully passing)**:
     * Submissions that pass all test cases
     * Code is efficient, well-structured, and handles all edge cases
     * Solution demonstrates thorough understanding of concepts
     * Code is well-documented and follows best practices
3. Calculate counts and percentages for each performance level.
4. Write a concise summary of overall class performance based on these categories.

### Stage 2: Error Identification and Categorization
Using the performance classifications from Stage 1, particularly focusing on the "Poor" and "Struggling" submissions:
1. Identify all errors in each submission. Start with these common error categories:
   - **Logic Errors**: Fundamental mistakes in the algorithm or program logic.
   - **Syntax Errors**: Issues with the syntax of the programming language.
   - **Data Type Errors**: Incorrect use or manipulation of data types.
   - **Control Flow Errors**: Problems with the sequence of execution (e.g., incorrect loop conditions, branching).
   - **Function/Method Errors**: Incorrect use or implementation of functions/methods.
   - **Data Structure Errors**: Incorrect use or implementation of data structures.
   - **Edge Case Handling**: Failure to handle unusual or boundary input values.
   - **Runtime Errors**: Errors that occur during the execution of the program.
   - **Algorithm Selection Errors**: Choosing inappropriate algorithms or approaches for solving problems.
   - **Variable Scope Issues**: Problems with variable accessibility (local vs. global).
   - **Off-by-One Errors**: Errors related to boundary conditions in loops or indexing.

   You may identify additional error categories if you observe patterns that don't fit well into the above categories. Include any additional categories in the all_categories_used array in Stage 6.

2. Count the frequency of each error category.
3. Select the top 5 most common error categories.
4. For each common error category:
   a. Calculate occurrence count and percentage among failing submissions.
   b. Write a clear description of the error pattern.
   c. Select a representative code example that clearly demonstrates the issue, using no more than 5-7 lines of code.
   d. Include the student_ids array listing all students who exhibited this error.
   e. Choose examples that are:
      - Clear: Demonstrate the error with minimal extraneous code
      - Typical: Represent the most common manifestation of the error
      - Concise: Short but with enough context to understand the issue
      - Educational: Show cases where the fix would be instructive

### Stage 3: Correlation and Pattern Analysis
Building on the error categorizations from Stage 2:
1. Analyze which error categories frequently co-occur in the same submissions.
2. Calculate correlation percentages for pairs of errors.
3. Identify the 3-5 strongest error correlations.
4. For each correlated pair:
   a. Report the correlation strength (percentage of failing submissions containing both errors).
   b. Report the correlation count (number of students showing both errors).
   c. Formulate a hypothesis explaining why these errors might be related.
   d. Select an example that clearly demonstrates both errors together, using no more than 8 lines of code.
   e. Include the student_ids array of students who exhibited this correlation.
   f. Choose examples where:
      - Both errors are clearly visible in proximity
      - The relationship between the errors is evident
      - The example supports your hypothesis about why they correlate

### Stage 4: Misconception Inference
Based on the error patterns and correlations identified in Stages 2 and 3:
1. Infer 3-5 potential misconceptions that might explain the observed error patterns.
2. For each misconception:
   a. Link it to specific error categories from earlier stages using the related_error_categories array.
   b. Calculate how many submissions exhibit error patterns suggesting this misconception (occurrence_count).
   c. Calculate the percentage of submissions showing this misconception (occurrence_percentage).
   d. Write a clear explanation of the likely misunderstanding.
   e. Select a representative code example that demonstrates this misconception, using no more than 5-7 lines of code.
   f. Include the student_ids array of students who exhibited this misconception.
   g. Choose examples that:
      - Clearly illustrate the conceptual misunderstanding
      - Represent common patterns across multiple students
      - Show the root cause rather than just symptoms

### Stage 5: Intervention Development
For each misconception identified in Stage 4:
1. Develop 2-3 specific instructional interventions that address the misconception.
2. For strongly correlated error pairs from Stage 3, suggest interventions that address both issues simultaneously.
3. Prioritize interventions based on the frequency and severity of associated errors.
4. For each intervention, specify:
   a. The target_misconception it addresses.
   b. An array of related_errors categories it targets.
   c. An array of instructional_strategies (2-3 specific approaches).
   d. A priority_level of "High", "Medium", or "Low".
   e. An estimated_effort of "Low", "Medium", or "High".
   f. An implementation_phase of "Immediate", "Short-term", or "Long-term".
5. Make recommendations concrete and actionable, such as:
   - Specific in-class activities
   - Targeted assignments
   - Code examples to discuss
   - Conceptual explanations to emphasize

### Stage 6: Additional Insights
After completing Stages 1-5, include:
1. Common good practices observed in successful submissions:
   a. For each good practice, provide:
      - A description of the practice.
      - Example code demonstrating this good practice.
      - A correct_implementation flag set to true.
2. Error category distribution:
   a. List all error categories used in the analysis in all_categories_used array.
   b. List any additional categories identified but not included in the main analysis in additional_categories_identified array.
3. Submission patterns:
   a. Provide an observation about when/how students submitted their work.
   b. Include a potential_impact analysis of these submission patterns.
   c. Include a submission_timeline object with dynamically determined time intervals as keys and objects containing:
      - count: number of submissions on that date
      - passing_rate: percentage of passing submissions on that date (e.g., "75%")
4. Comparative analysis:
   a. Include a success_factors array listing what successful students do differently.
   b. Include a challenge_factors array listing common challenges for struggling students.

## Output Structure

Format your output as a JSON object with the following structure, ensuring all requested statistics (counts and percentages) are included and formatted consistently with the "%" suffix:

	{
		"stage_1_performance_analysis": {
		"total_submissions": 50,
			"performance_distribution": {
			"Poor": {"count": 8, "percentage": "16.00%"},
			"Struggling": {"count": 10, "percentage": "20.00%"},
			"Good_Progress": {"count": 15, "percentage": "30.00%"},
			"Strong": {"count": 17, "percentage": "34.00%"}
		},
		"overall_summary": "The class shows a mixed performance. A significant portion (36%) of students are still struggling with fundamental concepts, as indicated by the 'Poor' and 'Struggling' categories. However, a majority (64%) are showing 'Good Progress' or are 'Strong', suggesting a general grasp of the core concepts."
	},
		"stage_2_error_analysis": {
		"top_errors": [
	{
	"category": "Case Sensitivity Issues",
	"occurrence_count": 25,
	"occurrence_percentage": "71.43%",
	"description": "Students are treating uppercase and lowercase versions of the same word as distinct.",
	"example_code": [
	"word_counts = {}",
	"sentence = 'The quick the brown'",
	"for word in sentence.split():",
	"  word_counts[word] = word_counts.get(word, 0) + 1",
	"// Result: {'The': 1, 'quick': 1, 'the': 1, 'brown': 1}"
	],
	"student_ids": [101, 105, 110, 115, 120]
	}
	]
	},
	"stage_3_correlation_analysis": {
	"error_correlations": [
	{
	"correlated_errors": ["Case Sensitivity Issues", "Punctuation Handling"],
	"correlation_count": 15,
	"correlation_percentage": "42.86%",
	"hypothesis": "Students who don't normalize text case also tend to overlook punctuation removal, suggesting a broader misunderstanding of text preprocessing.",
	"example_code": [
	"sentence = 'The quick, the brown.'",
	"words = sentence.split()",
	"for word in words:",
	"  word_counts[word] = word_counts.get(word, 0) + 1",
	"// Result: {'The': 1, 'quick,': 1, 'the': 1, 'brown.': 1}"
	],
	"student_ids": [101, 102, 105, 97, 98]
	}
	]
	},
	"stage_4_misconception_analysis": {
	"potential_misconceptions": [
	{
	"misconception": "Lack of understanding of string normalization techniques",
	"related_error_categories": ["Case Sensitivity Issues", "Tokenization Errors"],
	"occurrence_count": 25,
	"occurrence_percentage": "71.43%",
	"explanation": "Students may not realize the need to convert text to a consistent case for accurate word counting.",
	"example_code": [
	"sentence = 'The quick the brown'",
	"word_counts = {}",
	"for word in sentence.split():",
	"  word_counts[word] += 1",
	"// Without lowercasing"
	],
	"student_ids": [101, 105, 110, 115, 120]
	}
	]
	},
	"stage_5_instructional_plan": {
	"interventions": [
	{
	"target_misconception": "Lack of understanding of string normalization techniques",
	"related_errors": ["Case Sensitivity Issues", "Tokenization Errors"],
	"instructional_strategies": [
	"Explicitly teach and demonstrate string methods like '.lower()' and '.upper()' with examples.",
	"Include exercises that require students to normalize text before processing.",
	"Have students trace through code with and without normalization to see the difference in results."
	],
	"priority_level": "High",
	"estimated_effort": "Low",
	"implementation_phase": "Immediate"
	}
	]
	},
	"stage_6_additional_insights": {
	"common_good_practices": [
	{
	"description": "Several students effectively used dictionary's '.get(key, 0)' method for concise counting.",
	"example_code": [
	"word_counts[word] = word_counts.get(word, 0) + 1"
	],
	"correct_implementation": true
	}
	],
	"error_category_distribution": {
	"all_categories_used": ["Case Sensitivity Issues", "Punctuation Handling", "Tokenization Errors", "Algorithm Logic Errors", "Data Structure Selection"],
	"additional_categories_identified": ["Regular Expression Misuse", "String Method Confusion"]
	},
	"submission_patterns": {
	"observation": "Student submissions appear to cluster around two distinct dates (March 15 and March 19), suggesting some may have waited until close to the deadline.",
	"potential_impact": "Students submitting closer to the deadline show a higher rate of fundamental errors, suggesting rushed work or inadequate preparation time.",
	"submission_timeline": {
	"March 15 10:00-10:02": {"count": 5, "passing_rate": "80%"},
	"March 15 10:02-10:04": {"count": 7, "passing_rate": "71%"},
	"March 15 10:04-10:06": {"count": 8, "passing_rate": "75%"},
	"March 15 10:06-10:08": {"count": 10, "passing_rate": "60%"},
	"March 15 10:08-10:10": {"count": 8, "passing_rate": "50%"},
	"March 15 10:10-10:12": {"count": 7, "passing_rate": "43%"},
	"March 15 10:12-10:14": {"count": 5, "passing_rate": "40%"}
	}
	},
	"comparative_analysis": {
	"success_factors": [
	"Proper dictionary initialization using empty braces {}",
	"Use of dict.get() method for clean frequency counting",
	"Consistent case normalization with lower() or upper()",
	"Breaking the problem into logical steps"
	],
	"challenge_factors": [
	"Confusion between list and dictionary data structures",
	"Misunderstanding of how to accumulate counts",
	"Failure to recognize the case sensitivity requirement",
	"Incomplete understanding of function return types"
	]
	}
	}
	}
dont even give Here is the analysis of the student code submissions for the temperature converter problem:
..just give me the whole json object as response ...nothing else.
dont even write Here is the analysis of the student code submissions for the even number counter problem:.
JUST THE JSON OBJECT!`

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
	text := makeRequestClaude2(c, messages, problemID)
	return text
}

func ProcessIndividualSubmissions(w http.ResponseWriter, r *http.Request) string {

	problemIDStr := r.FormValue("problem_id")
	if problemIDStr == "" {
		http.Error(w, "Missing problem_id", http.StatusBadRequest)
		return ""
	}

	// Convert string to int
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		http.Error(w, "Invalid problem_id", http.StatusBadRequest)
		return ""
	}

	problemDescription, err := repository.GetProblemDescription(problemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch problem description: %v", err), http.StatusInternalServerError)
		return ""
	}

	studentIDStr := r.FormValue("student_id")
	if studentIDStr == "" {
		http.Error(w, "Missing problem_id", http.StatusBadRequest)
		return ""
	}

	// Convert string to int
	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		http.Error(w, "Invalid problem_id", http.StatusBadRequest)
		return ""
	}

	newStr := r.FormValue("isNew")
	isNew := false
	if newStr == "true" {
		isNew = true
	}

	if !isNew {
		summary, _ := repository.GetLatestScaffold(studentID, problemID)
		if summary != "" {
			return summary
		}
	}

	// Get the latest code snapshots from the database using the problem_id
	latestCodeSnapshot, err := repository.GetLatestCodeSnapshotForStudent(problemID, studentID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch latest code snapshots: %v", err), http.StatusInternalServerError)
		return ""
	}

	latestScaffolding, err := repository.GetLatestScaffoldingText(studentID, problemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch latest code snapshots: %v", err), http.StatusInternalServerError)
		return ""
	}

	// Construct a detailed prompt for error analysis
	NewInstructions := `You are assisting instructors in real-time classroom coding exercises, each lasting only 10-15 minutes. Your role is to analyze a student's current code submission, identify relevant error patterns, select appropriate scaffolding strategies, and generate multiple scaffolds clearly varying in effort and actionability. Your scaffolds must explicitly support student autonomy by prompting students to actively engage in problem-solving, reflect upon their coding decisions, and maintain responsibility for their learning outcomes. Scaffolds should facilitate productive struggle, fostering independence and deeper conceptual understanding rather than providing overly explicit or prescriptive solutions. The scaffolds you produce will be reviewed by instructors, who will select the most suitable scaffolds based on student needs and classroom context.

<exercise_description>
{Detailed description of the exercise would be placed here}
</exercise_description>

<student_code>
{The student's current code would be placed here}
</student_code>

<previous_scaffold>
{Previous scaffold content placed here, or "None" if this is the first scaffold}
</previous_scaffold>

# Programming Language
The student is coding in Python. Tailor all scaffolding approaches to Python-specific concepts, common Python errors, and Pythonic solutions.

# Time Constraint Guidance
Each scaffold MUST be implementable within the 10-15 minute exercise window. Consider:
- For Low-Effort scaffolds: Student should need <3 minutes to understand and apply
- For Moderate-Effort scaffolds: Student should need 3-7 minutes to process and implement
- For High-Effort scaffolds: Student should need 7-12 minutes to understand and integrate the concept
Avoid scaffolds requiring extensive code rewriting or introducing entirely new approaches that cannot be reasonably implemented in the remaining time.

# Scaffolding Strategies
1. **Gap-Fill Prompts:** Templates with blanks requiring student input.
2. **Incremental Hints:** Short, progressively explicit guidance.
3. **Minimal Correction Prompts:** Highlight specific errors with minimal suggested fixes.
4. **Guided Questions:** Reflective questions prompting conceptual thinking.
5. **Targeted Code Comments:** Specific, embedded comments guiding improvements without explicit solutions.
6. **Partial Worked Examples:** Concise examples aligned with the student's current approach.
7. **Reference Examples:** Adaptable examples illustrating relevant coding patterns.
8. **Error Message Translation:** Clear, student-friendly interpretations of errors.

# Effort Levels (Positive Struggle)
- **Low (Explicit):** Direct instructions; minimal cognitive challenge.
- **Moderate (Reflective):** Prompts thoughtful reflection and moderate reasoning.
- **High (Conceptual):** Requires deeper conceptual understanding and integrative thinking.

# Actionability Levels
- **High:** Clearly specifies immediate steps for student action.
- **Moderate:** Provides directional guidance, requiring some interpretation.
- **Low:** Offers conceptual insights without explicit instructions, significant interpretation required.

# Instructions
Analyze the student's submitted code and clearly describe identified error patterns. Generate multiple scaffold suggestions spanning a range of effort (low to high) and actionability (high to low). Explicitly support student autonomy, productive struggle, and continuity.

Scaffolds must incorporate the student's current code structure and approach whenever possible. Scaffolds should build upon what the student has already created rather than suggesting entirely new implementations. This ensures continuity in the student's thought process, honors their existing work, and makes scaffolds immediately actionable. Only suggest significant restructuring when the current approach contains fundamental conceptual errors that cannot be remedied through incremental improvements.

For each scaffold clearly state:
- **Strategy used**
- **Effort level (Low, Moderate, High)**
- **Actionability level (High, Moderate, Low)**
- **Content:** Scaffold content must be brief, clear, and immediately actionable—ideally presented in 1–3 concise sentences or a succinct code snippet, easily understandable within a short time frame.
- **Learning support:** Explicitly describe how the scaffold addresses misconceptions and supports student learning.
- **Implementation success metrics:** Observable indicators demonstrating the student’s successful application of scaffolds (e.g., "Student correctly implements loop structures instead of repetitive statements," "Student independently resolves index errors by accurately adjusting loop boundaries.")
- **Recommendation reasoning:** Clearly justify why this scaffold occupies its specific position in your ordered list.

**IMPORTANT:**  
Order scaffolds clearly by recommendation level. The FIRST scaffold listed should be your HIGHEST recommendation, informed by:
1. Clearly identified common errors in the student's code.
2. Previous scaffolds provided (if applicable).
3. Optimal progression in cognitive complexity and student support.
4. Feasibility of implementation within the time constraint.

# Response Format (JSON)
Your response must follow exactly this JSON format:

{
    exerciseDescription: "Design a function that accepts a string input and returns a dictionary where each key is a unique word from the text and each value represents how many times that word appears. The function should process text in a case-insensitive manner, meaning \"Hello\" and \"hello\" are counted as the same word. All punctuation should be ignored when identifying and counting words, ensuring only actual text is processed. Empty strings should not be counted as words. For example, passing \"Hello, hello world!\" to this function would return {'hello': 2, 'world': 1}.",
    errorPatterns: "The function iterates over individual characters rather than words, so it counts letters instead of whole words. This approach also breaks words at punctuation and spaces, meaning 'hello,' and 'hello' are treated differently or lost. Overall, the code mis‑interprets the specification to return word‑level frequencies with punctuation removed and case ignored.",
    reasoning: "The student's primary misconception is at the granularity of iteration (characters vs. words). My scaffold sequence moves from an immediately actionable code tweak (ensuring quick success within tight time constraints) toward more reflective, concept‑heavy prompts that deepen understanding. Each scaffold explicitly preserves the student's existing structure to honor their work and promote autonomy, while success metrics give instructors rapid checks. Ordering prioritizes quick wins first, then increasing cognitive challenge to sustain productive struggle.",
    studentCode: "def word_counts(text):\n    result = {}\n    for ch in text:\n        ch = ch.lower()\n        if ch.isalpha():\n            result[ch] = result.get(ch, 0) + 1\n    return result",
    previousScaffold: "1. What does the 'for ch in text' loop actually iterate over?\n2. The task says unique word – which Python method lets you split a string into words quickly?",
    scaffolds: [
        {
            strategy: "Minimal Correction Prompt",
            effortLevel: "Low",
            actionabilityLevel: "High",
            content: "🔧 **Quick fix** – keep your dictionary but swap the loop:\n'''for word in text.lower().split():          # iterate by WORD, not char\n    clean = ''.join(c for c in word if c.isalpha())  # strip punctuation\n    if clean:\n        result[clean] = result.get(clean, 0) + 1\n'''",
            learningSupport: "Shows exactly where to adjust iteration and sanitization while preserving the student's structure, letting them focus on one clear change rather than rewriting everything.",
            implementationSuccessMetrics: "Running 'word_counts('Hello, hello world!')' now returns '{'hello': 2, 'world': 1}'; no individual letters appear.",
            recommendationReasoning: "Fastest path to a correct result within 3 minutes; leverages student's existing code with minimal cognitive load.",
            estimatedTime: "3 minutes"
        },
        {
            strategy: "Guided Questions",
            effortLevel: "Moderate",
            actionabilityLevel: "Moderate",
            content: "❓ Reflect:\n1. What does the 'for ch in text' loop actually iterate over?\n2. The task says *unique word* – which Python method lets you split a string into words quickly?\n3. After splitting, how can you remove non‑alphabetic characters **inside** each word so 'hello,' and 'hello' match?",
            learningSupport: "Prompts metacognition about iteration granularity and text cleaning without prescribing code, encouraging the student to diagnose and repair independently.",
            implementationSuccessMetrics: "Student rewrites the loop to iterate over words, adds a cleaning step, and explains their choice during review.",
            recommendationReasoning: "Encourages productive struggle while remaining feasible (3–7 min). Suitable if instructors want students to reason before seeing code.",
            estimatedTime: "5 minutes"
        },
        {
            strategy: "Gap‑Fill Prompt",
            effortLevel: "High",
            actionabilityLevel: "Moderate",
            content: "Fill the blanks to upgrade your loop:\n'''python\nfor ___ in text.lower().split():\n    temp = ''.join([c for c in ___ if c.isalpha()])\n    if temp:\n        result[___] = result.get(temp, 0) + 1\n'''",
            learningSupport: "Requires the student to choose correct variable names and ensure consistency, reinforcing understanding of loop variables and dictionary updates.",
            implementationSuccessMetrics: "Blanks are correctly replaced; function passes instructor‑provided tests on mixed‑case, punctuated input.",
            recommendationReasoning: "Demands deeper engagement and variable‑scope awareness (7‑10 min) yet still anchored in existing code.",
            estimatedTime: "8 minutes"
        },
        {
            strategy: "Partial Worked Example",
            effortLevel: "High",
            actionabilityLevel: "Low",
            content: "Example flow (not full solution):\n1️⃣ 'text = \"Dogs, dogs & cats!\"'\n2️⃣ 'words = ['dogs', 'dogs', 'cats']'\n3️⃣ Update counts → '{'dogs': 2, 'cats': 1}'\nReplicate this three‑step process in your own code.",
learningSupport: "Illustrates the desired transformation pipeline, pushing the student to map each step onto their program logic without giving code.",
implementationSuccessMetrics: "Student articulates each stage in comments and adapts their function to mirror the pipeline.",
recommendationReasoning: "Great for students who learn best from concrete examples but still leaves significant synthesis work; reserved as last recommendation in case earlier scaffolds are too prescriptive.",
estimatedTime: "10 minutes"
}
]
}
you can use backticks instead of '''python so that the code comes in python if required
provide the value of previousScaffold: exactly as provided. if it is nothing then give the value as "none""
just give me the whole json object as response ...nothing else.
dont even write Here is the analysis of the student code submissions for the even number counter problem:.
JUST THE JSON OBJECT!
`

	// Build the full prompt
	prompt := fmt.Sprintf(
		`Problem Description:
%s

Instructions: %s
Code Submissions:
%s

previousScaffold: %s`,
		problemDescription,
		NewInstructions,
		latestCodeSnapshot,
		latestScaffolding,
	)

	// Create Claude request
	messages := []map[string]string{
		{"role": "user", "content": prompt},
	}

	// Create a Gin context from http.ResponseWriter and http.Request
	c, _ := gin.CreateTestContext(w)
	c.Request = r

	// Call Claude API request function
	text := makeRequestClaudeSC(c, messages, studentID, problemID)
	return text
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
func GetScaffoldingMaterial(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		ProblemID           int    `json:"problem_id"`
		ScaffoldingLevel    int    `json:"scaffolding_level"`
		ScaffoldingStrategy string `json:"scaffolding_strategy"`
	}

	// Decode the request body into the requestData struct
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	// Query the database for scaffolding material
	var scaffolding models.Scaffolding
	if err := models.DB.Where("problem_id = ? AND scaffolding_level = ? AND scaffolding_strategy = ?",
		requestData.ProblemID, requestData.ScaffoldingLevel, requestData.ScaffoldingStrategy).
		First(&scaffolding).Error; err != nil {
		// Handle error if no scaffolding material found or DB error
		http.Error(w, "Failed to fetch scaffolding material", http.StatusInternalServerError)
		return
	}

	// Send the response with scaffolding material
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Explicitly return 200 OK status code
	json.NewEncoder(w).Encode(struct {
		ScaffoldingMaterial string `json:"scaffolding_material"`
	}{
		ScaffoldingMaterial: scaffolding.ScaffoldingMaterial,
	})
}

func GetScaffolding(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		ProblemID int `json:"problem_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	var scaffoldings []models.Scaffolding
	if err := models.DB.Where("problem_id = ?", requestData.ProblemID).Find(&scaffoldings).Error; err != nil {
		http.Error(w, "Failed to fetch scaffolding records", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scaffoldings)
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
- The output **must pass the following check**:

  if !json.Valid([]byte(extractedJSON)) {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Extracted JSON is not valid"})
      return
  }
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
- The output **must pass the following check**:

  if !json.Valid([]byte(extractedJSON)) {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Extracted JSON is not valid"})
      return
  }
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
- The output **must pass the following check**:

  if !json.Valid([]byte(extractedJSON)) {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Extracted JSON is not valid"})
      return
  }
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
- The output **must pass the following check**:

  if !json.Valid([]byte(extractedJSON)) {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Extracted JSON is not valid"})
      return
  }
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
- The output **must pass the following check**:

  if !json.Valid([]byte(extractedJSON)) {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Extracted JSON is not valid"})
      return
  }

- Only include valid key-value pairs for "struggling", "developing", and "nearly_proficient".
- Do not include any additional text, formatting, or explanations outside of the JSON structure.
- The JSON keys must exactly match "struggling", "developing", and "nearly_proficient".
- The output should be **pure JSON** with no extra characters, sentences, or formatting.
- Do **not** write "Here is my attempt at the scaffolded solutions in valid JSON format:"—just return the JSON object directly.
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

type AddAPIKeyRequest struct {
	ID     int    `json:"id"`
	APIKey string `json:"api_key"`
}

func AddAPIKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req AddAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Trim input
	req.APIKey = strings.TrimSpace(req.APIKey)

	if req.ID == 0 || req.APIKey == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Update the API key using the helper function
	err := repository.UpdateAPIKeyByID(req.ID, req.APIKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update API key: %v", err), http.StatusInternalServerError)
		return
	}

	// Success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "API key updated successfully"})
}

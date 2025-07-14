package ai

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/GPTA/src/restHandlers"

	"github.com/GPTA/src/Grade"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/openAI"
	"github.com/GPTA/src/repository"
	"github.com/franciscoescher/goopenai"
	"github.com/gin-gonic/gin"
)

type ProblemDescription struct {
	Timestamp          string `json:"timestamp"`
	ProblemDescription string `json:"problem_description"`
	IsMCQ              bool   `json:"is_mcq"`
}

type CodeSnapshot struct {
	StudentID  int    `json:"student_id"`
	Timestamp  string `json:"timestamp"`
	Content    string `json:"content"`
	Grade      string `json:"grade"`
	SnapshotId int    `json:"snapshot_id"`
}

type CodeSnapshots struct {
	Entries []CodeSnapshot `json:"entries"`
}

type SubmissionTimeEntry struct {
	StudentID int    `json:"student_id"`
	Timestamp string `json:"timestamp"`
}

type SubmissionTimes struct {
	SubmissionTimes []SubmissionTimeEntry `json:"submission_times"`
}

type TAIntervention struct {
	StudentID int    `json:"student_id"`
	Timestamp string `json:"timestamp"`
	HelpStat  string `json:"help_stat"`
	Message   string `json:"message"`
}

type TAInterventions struct {
	Interventions []TAIntervention `json:"interventions"`
}

type PerformanceCategory struct {
	Count      int    `json:"count"`
	Percentage string `json:"percentage"`
}

type PerformanceDistribution struct {
	Correct     PerformanceCategory `json:"correct"`
	Incorrect   PerformanceCategory `json:"incorrect"`
	NotAssessed PerformanceCategory `json:"not_assessed"`
}

type OverallAssessment struct {
	TotalEntries            int                     `json:"total_entries"`
	PerformanceDistribution PerformanceDistribution `json:"performance_distribution"`
}

type IndividualAssessment struct {
	StudentID        int    `json:"student_id"`
	PerformanceLevel string `json:"performance_level"`
}

type TopError struct {
	Category             string   `json:"category"`
	OccurrenceCount      int      `json:"occurrence_count"`
	OccurrencePercentage string   `json:"occurrence_percentage"`
	Description          string   `json:"description"`
	ExampleCode          []string `json:"example_code"`
	StudentIDs           []int    `json:"student_ids"`
}

type ErrorCorrelation struct {
	CorrelatedErrors      []string `json:"correlated_errors"`
	CorrelationCount      int      `json:"correlation_count"`
	CorrelationPercentage string   `json:"correlation_percentage"`
	Hypothesis            string   `json:"hypothesis"`
	ExampleCode           []string `json:"example_code"`
	StudentIDs            []int    `json:"student_ids"`
}

type PotentialMisconception struct {
	Misconception                   string   `json:"misconception"`
	RelatedErrorCategories          []string `json:"related_error_categories"`
	OccurrenceCount                 int      `json:"occurrence_count"`
	OccurrencePercentage            string   `json:"occurrence_percentage"`
	ExplanationDiagnostic           string   `json:"explanation_diagnostic"`
	ExampleCodeError                []string `json:"example_code_error"`
	StudentIDs                      []int    `json:"student_ids"`
	SuggestedExplanationForStudents string   `json:"suggested_explanation_for_students"`
	CorrectCodeExample              []string `json:"correct_code_example"`
	FollowUpQuestion                string   `json:"follow_up_question"`
}

type AggregateAnalysis struct {
	TopErrors               []TopError               `json:"top_errors"`
	ErrorCorrelations       []ErrorCorrelation       `json:"error_correlations"`
	PotentialMisconceptions []PotentialMisconception `json:"potential_misconceptions"`
}

type ProblemSummary struct {
	Title string `json:"title"`
}

type AnalysisData struct {
	ProblemSummary       ProblemSummary         `json:"problem_summary"`
	IsEnabled            bool                   `json:"isEnable"`
	OverallAssessment    OverallAssessment      `json:"overall_assessment"`
	IndividualAssessment []IndividualAssessment `json:"individual_assessment"`
	AggregateAnalysis    AggregateAnalysis      `json:"aggregate_analysis"`
}

func HandleMergedData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	query := r.URL.Query()
	problemIDStr := query.Get("problem_id")
	regenerateStr := query.Get("regenerate")

	// 🔁 Convert problemID to int
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		http.Error(w, "Invalid problem_id", http.StatusBadRequest)
		return
	}

	// 🔁 Convert regenerate to bool (optional; use if you need it)
	generateNew := regenerateStr == "true"

	uploadedAt, description, filename, err := repository.GetProblemUploadTimeAndDescription(problemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch problem description: %v", err), http.StatusInternalServerError)
		return
	}
	var problemDescription ProblemDescription
	if models.ActiveProblems[filename] == nil || models.ActiveProblems[filename].Info == nil {
		problemDescription = ProblemDescription{
			Timestamp:          uploadedAt.Format("2006-01-02 15:04:05"),
			ProblemDescription: description,
			IsMCQ:              false, // Default to false if no active problem found
		}
	} else {
		problemDescription = ProblemDescription{
			Timestamp:          uploadedAt.Format("2006-01-02 15:04:05"),
			ProblemDescription: description,
			IsMCQ:              models.ActiveProblems[filename].Info.Answer != "",
		}
	}

	codeSnapshotss, err := repository.GetCompleteCodeSnapshotsByProblemID(problemID)
	if err != nil {
		http.Error(w, "Error fetching code snapshots", http.StatusInternalServerError)
		return
	}

	allCodeSnapshots, err := repository.GetAllCodeSnapshotsByProblemID(problemID)
	if err != nil {
		http.Error(w, "Error fetching code snapshots", http.StatusInternalServerError)
		return
	}

	var grades []Grade.Grade
	if err := models.DB.Where("problem_id = ?", problemID).Find(&grades).Error; err != nil {
		http.Error(w, "Error fetching grades", http.StatusInternalServerError)
		return
	}

	// Step 1: Build a map for fast lookup: student_id -> grade
	gradeMap := make(map[int]string)
	for _, g := range grades {
		gradeMap[g.StudentID] = g.Grade
	}

	// Step 2: Merge into snapshots
	formattedSnapshots := []CodeSnapshot{}
	for _, snap := range allCodeSnapshots {
		formattedSnapshots = append(formattedSnapshots, CodeSnapshot{
			StudentID:  snap.StudentID,
			Timestamp:  snap.LastUpdatedAt.Format("2006-01-02 15:04:05"),
			Content:    snap.Code,
			Grade:      gradeMap[snap.StudentID],
			SnapshotId: snap.ID,
		})
	}

	formattedSnapshotsCL := []CodeSnapshot{}
	for _, snap := range codeSnapshotss {
		formattedSnapshotsCL = append(formattedSnapshotsCL, CodeSnapshot{
			StudentID:  snap.StudentID,
			Timestamp:  snap.LastUpdatedAt.Format("2006-01-02 15:04:05"),
			Content:    snap.Code,
			Grade:      gradeMap[snap.StudentID],
			SnapshotId: snap.ID,
		})
	}

	// Step 3: Initialize performance counters
	performanceCounts := map[string]int{
		"correct":      0,
		"incorrect":    0,
		"not_assessed": 0,
	}

	individualAssessment := []IndividualAssessment{}
	total := len(codeSnapshotss)

	for _, snap := range codeSnapshotss {
		grade := gradeMap[snap.StudentID]
		performanceLevel := ""

		switch grade {
		case "correct":
			performanceLevel = "Correct"
			performanceCounts["correct"]++
		case "incorrect":
			performanceLevel = "Incorrect"
			performanceCounts["incorrect"]++
		default:
			performanceLevel = "NotAssessed"
			performanceCounts["not_assessed"]++
		}

		individualAssessment = append(individualAssessment, IndividualAssessment{
			StudentID:        snap.StudentID,
			PerformanceLevel: performanceLevel,
		})
	}

	// Helper to format percentage
	formatPercent := func(count, total int) string {
		if total == 0 {
			return "0.00%"
		}
		return fmt.Sprintf("%.2f%%", float64(count)*100/float64(total))
	}
	codeSnapshots := CodeSnapshots{
		Entries: formattedSnapshots,
	}
	submissions, err := repository.GetSubmissionsByProblemID(problemID)
	if err != nil {
		http.Error(w, "Error fetching submission time", http.StatusInternalServerError)
		return
	}

	formattedSubmissionTimes := []SubmissionTimeEntry{}
	for _, sub := range submissions {
		formattedSubmissionTimes = append(formattedSubmissionTimes, SubmissionTimeEntry{
			StudentID: sub.StudentID,
			Timestamp: sub.CodeSubmittedAt.Format("2006-01-02 15:04:05"),
		})
	}

	submissionTimes := SubmissionTimes{
		SubmissionTimes: formattedSubmissionTimes,
	}

	studentStatuses, err := repository.GetStudentStatusesByProblemID(problemID)
	if err != nil {
		http.Error(w, "Error fetching student statuses", http.StatusInternalServerError)
		return
	}

	interventions := []TAIntervention{}
	students := repository.GetAllStudents() // Get students map for FetchMessages

	for _, status := range studentStatuses {
		if status.HelpStat != "" {
			intervention := TAIntervention{
				StudentID: status.StudentID,
				Timestamp: status.LastUpdatedAt.Format("2006-01-02 15:04:05"),
				HelpStat:  status.HelpStat,
			}

			// If student asked for help, fetch their messages
			if status.HelpStat == "Asked for help" {
				messages, err := restHandlers.FetchMessages(students, problemID, status.StudentID, "teacher")
				if err == nil && len(messages) > 0 {
					// sort messages by the GivenAt field
					sort.Slice(messages, func(i, j int) bool {
						return messages[i].GivenAt.After(messages[j].GivenAt)
					})
					// Get the most recent help request message
					for _, msg := range messages {
						if msg.Role == "student" {
							intervention.Message = msg.Message
							break
						}
					}
				}
			}

			interventions = append(interventions, intervention)
		}
	}

	taInterventions := TAInterventions{
		Interventions: interventions,
	}

	var analysisData AnalysisData

	if generateNew {
		aggregateAnalysis, err := makeRequest(description, formattedSnapshotsCL, len(codeSnapshotss), problemID, gradeMap, w, r)
		if err != nil {
			http.Error(w, "Error generating analysis data", http.StatusInternalServerError)
			return
		}
		analysisData.IsEnabled = true
		analysisData.AggregateAnalysis = aggregateAnalysis
	} else {
		jsonText, _ := repository.GetLatestCodeInsightByProblemID(problemID)
		if jsonText != nil {
			if err := json.Unmarshal([]byte(jsonText.Response), &analysisData); err != nil {
				log.Printf("Unmarshal error: %v", err)
			}
			analysisData.IsEnabled = true
		} else {
			analysisData = AnalysisData{
				ProblemSummary: ProblemSummary{
					Title: filename,
				},
				IsEnabled: false,
				OverallAssessment: OverallAssessment{
					TotalEntries: total,
					PerformanceDistribution: PerformanceDistribution{
						Correct: PerformanceCategory{
							Count:      performanceCounts["correct"],
							Percentage: formatPercent(performanceCounts["correct"], total),
						},
						Incorrect: PerformanceCategory{
							Count:      performanceCounts["incorrect"],
							Percentage: formatPercent(performanceCounts["incorrect"], total),
						},
						NotAssessed: PerformanceCategory{
							Count:      performanceCounts["not_assessed"],
							Percentage: formatPercent(performanceCounts["not_assessed"], total),
						},
					},
				},
				IndividualAssessment: individualAssessment,
			}

		}
	}
	analysisData.OverallAssessment.TotalEntries = total
	analysisData.OverallAssessment.PerformanceDistribution = PerformanceDistribution{
		Correct: PerformanceCategory{
			Count:      performanceCounts["correct"],
			Percentage: formatPercent(performanceCounts["correct"], total),
		},
		Incorrect: PerformanceCategory{
			Count:      performanceCounts["incorrect"],
			Percentage: formatPercent(performanceCounts["incorrect"], total),
		},
		NotAssessed: PerformanceCategory{
			Count:      performanceCounts["not_assessed"],
			Percentage: formatPercent(performanceCounts["not_assessed"], total),
		},
	}

	analysisData.IndividualAssessment = individualAssessment
	merged := map[string]interface{}{
		"analysisData":       analysisData,
		"problemDescription": problemDescription,
		"codeSnapshots":      codeSnapshots,
		"submissionTimes":    submissionTimes,
		"taInterventions":    taInterventions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(merged)
}

func makeRequest(problemDescription string, formattedSnapshots []CodeSnapshot, NumberOfStudents, problemID int, gradeMap map[int]string, w http.ResponseWriter, r *http.Request) (AggregateAnalysis, error) {

	NewInstructions := `# LLM Prompt: Analyze, Assess, and Generate Remediation Ideas for CS1 Student Code Submissions

## Role and Goal

You are an expert AI assistant specializing in analyzing and assessing student code submissions for introductory computer science exercises (CS1/CS2 level). Your goal is to process a batch of student code submissions for a specific programming problem, **evaluate each submission's likely correctness and adherence to requirements based *solely* on static code analysis**, and generate a structured JSON summary. This summary should provide instructors with actionable insights for **Monitoring** student progress (performance levels), **Analyzing** common issues (errors, correlations, misconceptions), and **Responding** with targeted instructional support (explanations, examples, follow-up questions).

## Input Data Provided to You

You will be provided with the following inputs:

1.  **Problem Description:** The full text of the programming problem, including examples and specific requirements/constraints (e.g., functions not to use).
2.  **Student Submissions:** A JSON list where each object represents a single student's submission, containing:
3. **GradeMap:** A map of student Id and Assigned grade. If the map does not contain a students value, assign it to "NotAssessed"
* 'student_id': A unique identifier for the student.
* 'timestamp': The time of submission.
* 'content': A string containing the student's source code.
* **(Note: Assessment results like test pass/fail counts are NOT provided. You must infer correctness and errors from the code.)**

## Analysis Process and Output Structure

Analyze the provided data using the following staged process. You must perform the detailed error/misconception analysis described internally, as it is required for the aggregate stages. Format your final output **exclusively** as a single JSON object adhering to the structure specified at the end.

**Crucial Task:** Your primary challenge is to **simulate an assessment** through static analysis of each student's 'content' against the 'Problem Description'.

### Stage 1: Individual Code Assessment and Classification

For **each** student submission:

1.  **Analyze Code Logic & Requirements:**
* Does the code attempt to solve the correct problem as described?
* Does the core algorithm appear logically sound for typical cases?
* Does the code adhere to all specific requirements mentioned in the 'Problem Description' (e.g., not using forbidden functions like 'min()'/'max()', specific output formats)?
* **Internally identify** potential logical errors, incorrect initializations, mishandled edge cases (e.g., empty lists, single items, negatives, zeros), inefficiencies, requirement violations, or other flaws based on the code structure. This internal analysis is crucial for later stages.
* Estimate the likelihood of runtime errors ('IndexError', 'TypeError', etc.) based on the logic.
2.  **Classify Performance Level (Inferred): Use Grade map to do this
3.  **Output (Simplified):** Contribute an object containing only the 'student_id' and your inferred 'performance_level' to the 'individual_assessment' array in the final JSON for *each* student.

### Stage 2: Error Identification and Categorization (Aggregate)

Focusing primarily on errors *you identified internally* in submissions classified as "correct" or "incorrect":

1.  **Consolidate & Categorize Inferred Errors:** Group the errors identified across failing submissions using these categories:
* **Requirement Violation:** Code ignores explicit problem constraints.
* **Misinterpretation of Problem:** Code solves the wrong problem.
* **Logic Error:** Flawed reasoning in the algorithm.
* **Initialization Error:** Incorrect starting values for variables.
* **Control Flow Error:** Incorrect loops, branching, recursion.
* **Off-by-One Error:** Loop boundaries, indexing issues.
* **Edge Case Handling Error:** Failure on non-standard valid inputs.
* **Data Type Error:** Incorrect use or conversion of types.
* **Data Structure Error:** Incorrect use of lists, dictionaries, etc.
* **Function/Method Error:** Issues with definition, calls, parameters, returns.
* **Variable Scope Error:** Misunderstanding local vs. global scope.
* **Inefficiency/Suboptimal Algorithm:** Correct but slow/resource-intensive solution.
* **Potential Runtime Error:** High likelihood of crash (IndexError, TypeError, etc.).
* *(You may identify additional specific error patterns if frequent and distinct).*
2.  **Frequency Analysis:** Count occurrences for each category among "correct"/"incorrect".
3.  **Select Top Errors:** Identify the top 5 most frequent *inferred* error categories.
4.  **Output:** For each top error, populate an object in the 'top_errors' array containing 'category', 'occurrence_count', 'occurrence_percentage' (of failing students, format "XX.XX%"), 'description', 'example_code' (concise snippet illustrating the error), and 'student_ids'.

### Stage 3: Correlation and Pattern Analysis (Aggregate)

Analyze which *inferred* error categories (from Stage 2) frequently co-occur within the same "correct" or "incorrect" submissions.

1.  **Identify Strong Correlations:** Find the 3-5 strongest co-occurrence pairs among "correct"/"incorrect".
2.  **Output:** For each pair, populate an object in the 'error_correlations' array containing 'correlated_errors' (list of 2 categories), 'correlation_count', 'correlation_percentage' (of failing students, format "XX.XX%"), 'hypothesis' (why they might be linked), 'example_code' (concise snippet showing both errors), and 'student_ids'.

### Stage 4: Potential Misconception Inference and Remediation Content (Aggregate)

Based on the top *inferred* errors, correlations, and code patterns:

1.  **Infer Misconceptions:** Identify 1-3 high-level *potential* underlying conceptual misunderstandings likely explaining prevalent error patterns among "correct"/"incorrect" students.
2.  **Generate Remediation Content:** For each inferred misconception, *also* generate content suitable for instructor intervention (for the "Respond" dashboard).
3.  **Output:** For each inferred misconception, populate an object in the 'potential_misconceptions' array containing:
* 'misconception': Concise description of the potential misunderstanding.
* 'related_error_categories': List of inferred error categories strongly associated.
* 'occurrence_count': Approximate number of failing students whose inferred errors align.
* 'occurrence_percentage': Approximate percentage of failing students potentially affected (format "XX.XX%").
* 'explanation_diagnostic': Clear explanation of the likely misunderstanding (for the instructor's analysis).
* 'example_code_error': Concise (max 5-7 lines) code snippet vividly illustrating the *result* of this misconception (can reuse from errors/correlations if appropriate).
* 'student_ids': Array of 'student_id's of failing students whose code strongly suggests this misconception.
* **'suggested_explanation_for_students':** **(New)** A brief, clear, student-friendly explanation of the correct concept or why the misconception leads to errors. Suitable for direct use or adaptation by the instructor.
* **'correct_code_example':** **(New)** A minimal, correct code snippet (max 5-7 lines) demonstrating the *proper* way to handle the specific concept related to the misconception.
* **'follow_up_question':** **(New)** A short question (e.g., conceptual, code prediction, fill-in-the-blank) designed to check student understanding after the explanation.

## Important Considerations & Limitations

* **Static Analysis Only:** Your assessment is based purely on reading the code. You cannot execute it. Inferences about errors, performance levels, and the generated remediation content require instructor validation.
* **Focus on Clarity:** Provide clear, concise descriptions, examples, explanations, and questions. Ensure generated code examples are minimal and directly relevant.
* **Adhere to JSON:** Ensure the entire output is a single, valid JSON object matching the structure below.

## Final Output Format (Single JSON Object)

'''json
  "aggregate_analysis": {
    "top_errors": [
      {
        "category": "Initialization Error",
        "occurrence_count": 18,
        "occurrence_percentage": "90.00%",
        "description": "Students incorrectly initialized min/max variables (e.g., to 0), likely causing failure with negative numbers.",
        "example_code": [ /* ... */ ],
        "student_ids": [ /* ... */ ]
      }
      // ... other top errors
    ],
    "error_correlations": [
      {
        "correlated_errors": ["Initialization Error", "Edge Case Handling Error"],
        "correlation_count": 15,
        "correlation_percentage": "75.00%",
        "hypothesis": "Incorrect initialization likely leads directly to failure on edge cases involving values outside the initial assumption.",
        "example_code": [ /* ... */ ],
        "student_ids": [ /* ... */ ]
      }
      // ... other strong correlations
    ],
    "potential_misconceptions": [
      {
        "misconception": "Assuming default/zero initialization works universally for finding extremes.",
        "related_error_categories": ["Initialization Error", "Edge Case Handling Error", "Logic Error"],
        "occurrence_count": 18,
        "occurrence_percentage": "90.00%",
        "explanation_diagnostic": "Students often default min/max to 0, suggesting a lack of consideration for input ranges like all negatives.",
        "example_code_error": [
            "def find_max(nums):",
            "  max_v = 0 # Misconception here",
            "  for n in nums: # Fails if nums = [-10, -5, -2]",
            "    if n > max_v:",
            "      max_v = n",
            "  return max_v"
         ],
        "student_ids": [ /* ... */ ],
        // --- Fields for Respond Dashboard ---
        "suggested_explanation_for_students": "When finding a minimum or maximum, initializing your tracking variable to a value like 0 can fail if all numbers are outside that range (e.g., all negative). A safer way is to initialize using the *first element* of the list, then loop through the *rest*.",
        "correct_code_example": [
            "def find_max_safe(nums):",
            "  if not nums: return None # Handle empty list",
            "  max_v = nums[0] # Initialize with first element",
            "  for i in range(1, len(nums)):",
            "    if nums[i] > max_v:",
            "      max_v = nums[i]",
            "  return max_v"
        ],
        "follow_up_question": "Consider 'find_max_safe([-5, -12, -3])'. What value will 'max_v' be initialized to, and what will the function return?"
      }
      // ... other potential misconceptions with remediation content
    ]
  }
just give me the whole json object as response ...nothing else.
dont even write Here is the analysis of the student code submissions for the even number counter problem:.
JUST THE JSON OBJECT! so that the output is  in the format that I just have to unstructured it to the below go object and it works
dont even give '''json


type AggregateAnalysis struct {
TopErrors               []TopError               'json:"top_errors"'
ErrorCorrelations       []ErrorCorrelation       'json:"error_correlations"'
PotentialMisconceptions []PotentialMisconception 'json:"potential_misconceptions"'
}
type TopError struct {
Category             string   'json:"category"'
OccurrenceCount      int      'json:"occurrence_count"'
OccurrencePercentage string   'json:"occurrence_percentage"'
Description          string   'json:"description"'
ExampleCode          []string 'json:"example_code"'
StudentIDs           []int    'json:"student_ids"'
}
type ErrorCorrelation struct {
CorrelatedErrors      []string 'json:"correlated_errors"'
CorrelationCount      int      'json:"correlation_count"'
CorrelationPercentage string   'json:"correlation_percentage"'
Hypothesis            string   'json:"hypothesis"'
ExampleCode           []string 'json:"example_code"'
StudentIDs            []int    'json:"student_ids"'
}

type PotentialMisconception struct {
Misconception                   string   'json:"misconception"'
RelatedErrorCategories          []string 'json:"related_error_categories"'
OccurrenceCount                 int      'json:"occurrence_count"'
OccurrencePercentage            string   'json:"occurrence_percentage"'
ExplanationDiagnostic           string   'json:"explanation_diagnostic"'
ExampleCodeError                []string 'json:"example_code_error"'
StudentIDs                      []int    'json:"student_ids"'
SuggestedExplanationForStudents string   'json:"suggested_explanation_for_students"'
CorrectCodeExample              []string 'json:"correct_code_example"'
FollowUpQuestion                string   'json:"follow_up_question"'
}
`

	prompt := fmt.Sprintf(
		"Problem Description:\n%s\n\nInstructions: %s\nTotal Students:\n%sStudent Submissions:\n%s",
		problemDescription, NewInstructions, NumberOfStudents, formattedSnapshots,
	)

	// Create a Gin context from http.ResponseWriter and http.Request
	c, _ := gin.CreateTestContext(w)
	c.Request = r

	jsonText := ""

	provider, err := openAI.GetMostRecentlyUpdatedProvider()
	if err != nil {
		log.Println("Error determining selected AI provider:", err)
	} else {
		fmt.Println("Currently selected provider:", provider.Name)
	}
	if provider.ID == 1 {
		messages := []map[string]string{
			{"role": "user", "content": prompt},
		}
		jsonText = openAI.MakeRequestClaudeAnalyze(c, messages)
	}
	if provider.ID == 2 {
		var rawMessages = []map[string]string{
			{"role": "user", "content": prompt},
		}

		var messages []goopenai.Message
		for _, m := range rawMessages {
			messages = append(messages, goopenai.Message{
				Role:    m["role"],
				Content: m["content"],
			})
		}
		jsonText = openAI.MakeRequestOpenAIAnalyze(c, messages)
	}
	if provider.ID == 3 {
		messages := []map[string]string{
			{"role": "user", "content": prompt},
		}
		jsonText = openAI.MakeRequestGeminiAnalyze(c, messages)
	}
	var aggregateAnalysis AggregateAnalysis

	if jsonText == "" {
		log.Printf("Unmarshal error: %v", err)
		return aggregateAnalysis, errors.New("could not generate response")
	}

	cleanedJSON, err := CleanAndUnmarshal(jsonText, &aggregateAnalysis)
	if err != nil {
		log.Printf("Failed to clean JSON: %v", err)
		return aggregateAnalysis, err
	}

	// Extract the "aggregate_analysis" field manually
	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(cleanedJSON), &raw); err != nil {
		log.Printf("Unmarshal into RawMessage map failed: %v", err)
		return aggregateAnalysis, err
	}

	inner, ok := raw["aggregate_analysis"]
	if !ok {
		log.Printf("Key 'aggregate_analysis' not found in JSON")
		return aggregateAnalysis, errors.New("missing 'aggregate_analysis' key in JSON")
	}

	if err := json.Unmarshal(inner, &aggregateAnalysis); err != nil {
		log.Printf("Failed to unmarshal inner aggregate_analysis: %v", err)
		return aggregateAnalysis, err
	}

	// Save to DB (if needed)
	_, err = repository.AddCodeInsight(problemID, cleanedJSON, time.Now())
	if err != nil {
		log.Printf("Failed to add code insight: %v", err)
	}

	return aggregateAnalysis, nil
}
func CleanAndUnmarshal(jsonText string, target interface{}) (string, error) {
	// Trim leading/trailing whitespace
	jsonText = strings.TrimSpace(jsonText)

	// Remove ```json and ``` or just ``` if present
	if strings.HasPrefix(jsonText, "```") {
		re := regexp.MustCompile("(?s)```(?:json)?\\s*(.*?)\\s*```")
		matches := re.FindStringSubmatch(jsonText)
		if len(matches) >= 2 {
			jsonText = matches[1]
		} else {
			return "", errors.New("failed to extract valid JSON from code block")
		}
	}

	// Now unmarshal cleaned JSON
	if err := json.Unmarshal([]byte(jsonText), target); err != nil {
		log.Printf("Unmarshal error: %v", err)
		return "", err
	}

	return jsonText, nil
}

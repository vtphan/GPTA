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
	FeedbackToStudents   string   `json:"feedback_to_students"`
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

	NewInstructions := restHandlers.AIPromptsMap["analyze-class-performance"]

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

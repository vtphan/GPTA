package openAI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/GPTA/src/models"
	"github.com/GPTA/src/repository"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//func makeRequest(c *gin.Context, messages []goopenai.Message) {
//	apiKey := OpenaiAPIKey
//	organization := ""
//
//	client := goopenai.NewClient(apiKey, organization)
//
//	r := goopenai.CreateChatCompletionsRequest{
//		Model:       ChatGPTModel,
//		Messages:    messages,
//		Temperature: 0,
//	}
//	fmt.Print(messages)
//	for i := 0; i < NumRetry; i++ {
//		completions, err := client.CreateChatCompletions(context.Background(), r)
//		if err != nil {
//			panic(err)
//		}
//		// fmt.Println("Printing Content:")
//		// fmt.Print(completions)
//		var feedbacksJSON map[string]interface{}
//		err = json.Unmarshal([]byte(completions.Choices[0].Message.Content), &feedbacksJSON)
//		if err == nil {
//			c.JSON(http.StatusOK, feedbacksJSON)
//			return
//		}
//		// fmt.Println(feedbacksJSON)
//	}
//	c.JSON(http.StatusInternalServerError, "Couldn't process the request!")
//}

//func makeRequest2(c *gin.Context, messages []goopenai.Message) {
//	apiKey := OpenaiAPIKey
//	organization := ""
//
//	client := goopenai.NewClient(apiKey, organization)
//
//	r := goopenai.CreateChatCompletionsRequest{
//		Model:       ChatGPTModel,
//		Messages:    messages,
//		Temperature: 0,
//	}
//	for i := 0; i < NumRetry; i++ {
//		completions, err := client.CreateChatCompletions(context.Background(), r)
//		if err != nil {
//			fmt.Println(err)
//			continue
//		}
//
//		// Return the response in a structured JSON format
//		c.JSON(http.StatusOK, gin.H{
//			"response": completions.Choices[0].Message.Content,
//		})
//		return
//	}
//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't process the request!"})
//}

const ClaudeAPIKey = ""
const ClaudeEndpoint = "https://api.anthropic.com/v1/messages"
const ClaudeModel = "claude-3-opus-20240229" // Adjust model as needed

func MakeRequestClaude(c *gin.Context, messages []map[string]string) {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"model":      ClaudeModel,
		"messages":   messages,
		"max_tokens": 4096,
	})

	req, err := http.NewRequest("POST", ClaudeEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// ✅ Use "x-api-key" instead of "Authorization"
	req.Header.Set("x-api-key", ClaudeAPIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Anthropic-Version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API request failed"})
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	// Parse JSON response
	var responseJSON map[string]interface{}
	json.Unmarshal(body, &responseJSON)

	// Extract text from the content array
	if contentArray, found := responseJSON["content"].([]interface{}); found && len(contentArray) > 0 {
		// Extract the text field from the first content object
		if firstContent, ok := contentArray[0].(map[string]interface{}); ok {
			if text, exists := firstContent["text"].(string); exists {
				c.JSON(http.StatusOK, gin.H{"response": text})
				return
			}
		}
	}

	// Handle unexpected response structure
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from Claude API"})
	// todo : save it to db
}

func makeRequestClaude2(c *gin.Context, messages []map[string]string, problemId int) string {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"model":      ClaudeModel,
		"messages":   messages,
		"max_tokens": 4096,
	})

	req, err := http.NewRequest("POST", ClaudeEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return ""
	}

	req.Header.Set("x-api-key", ClaudeAPIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Anthropic-Version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API request failed"})
		return ""
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	// Parse JSON response
	var responseJSON map[string]interface{}
	json.Unmarshal(body, &responseJSON)

	// Extract text response
	if contentArray, found := responseJSON["content"].([]interface{}); found && len(contentArray) > 0 {
		if firstContent, ok := contentArray[0].(map[string]interface{}); ok {
			if text, exists := firstContent["text"].(string); exists {
				err = repository.SaveClassFeedback(problemId, text)
				if err != nil {
					fmt.Println(err)
				}
				return text
			}
		}
	}

	return ""
}

func makeRequestClaudeSC(c *gin.Context, messages []map[string]string, studentId, problemId int) string {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"model":      ClaudeModel,
		"messages":   messages,
		"max_tokens": 4096,
	})

	req, err := http.NewRequest("POST", ClaudeEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return ""
	}

	req.Header.Set("x-api-key", ClaudeAPIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Anthropic-Version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API request failed"})
		return ""
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	// Parse JSON response
	var responseJSON map[string]interface{}
	json.Unmarshal(body, &responseJSON)

	// Extract text response
	if contentArray, found := responseJSON["content"].([]interface{}); found && len(contentArray) > 0 {
		if firstContent, ok := contentArray[0].(map[string]interface{}); ok {
			if text, exists := firstContent["text"].(string); exists {
				err = repository.SaveScaffoldingFeedback(studentId, problemId, text)
				if err != nil {
					fmt.Println(err)
				}
				return text
			}
		}
	}

	return ""
}

func makeRequestClaude3(c *gin.Context, messages []map[string]string, problemId int) {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"model":      ClaudeModel,
		"messages":   messages,
		"max_tokens": 4096, // Increased token limit for detailed feedback
	})

	req, err := http.NewRequest("POST", ClaudeEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("x-api-key", ClaudeAPIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Anthropic-Version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API request failed"})
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	// Parse top-level JSON response
	var responseJSON map[string]interface{}
	if err := json.Unmarshal(body, &responseJSON); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse API response"})
		return
	}

	// Extract the actual feedback JSON from the "text" field inside "content"
	var extractedJSON string
	if contentArray, found := responseJSON["content"].([]interface{}); found && len(contentArray) > 0 {
		if firstContent, ok := contentArray[0].(map[string]interface{}); ok {
			if text, exists := firstContent["text"].(string); exists {
				extractedJSON = text
			}
		}
	}

	if extractedJSON == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response structure from Claude API"})
		return
	}

	// Parse extracted JSON into feedbackList
	var feedbackList []struct {
		StudentID   int    `json:"student_id"`
		Explanation string `json:"explanation"`
		Percentage  int    `json:"percentage"`
	}

	if err := json.Unmarshal([]byte(extractedJSON), &feedbackList); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse extracted JSON response"})
		return
	}

	// Save each student's progress to the database
	for _, feedback := range feedbackList {
		studentProgress := models.StudentProgress{
			ProblemID:   problemId,
			StudentID:   feedback.StudentID,
			Explanation: feedback.Explanation,
			Percentage:  feedback.Percentage,
		}

		if err := repository.SaveStudentProgress(studentProgress); err != nil {
			fmt.Println("Error saving student progress:", err)
		}
		if err := repository.UpdateStudentPercentStat(studentProgress.Percentage, feedback.Explanation, time.Now(), feedback.StudentID, problemId); err != nil {
			fmt.Println("Error saving student progress:", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student progress saved successfully"})
}

func makeRequestClaude4(c *gin.Context, messages []map[string]string, problemId int, scaffoldingStrategy int) {
	var existingScaffolding []models.Scaffolding
	models.DB.Where("problem_id = ? AND scaffolding_strategy = ?", problemId, scaffoldingStrategy).Find(&existingScaffolding)

	existingLevels := make(map[int]bool)
	for _, record := range existingScaffolding {
		existingLevels[record.ScaffoldingLevel] = true
	}

	if len(existingLevels) >= 3 {
		c.JSON(http.StatusConflict, gin.H{"message": "Scaffolding already exists for the given problem ID and strategy"})
		return
	}

	requestBody, _ := json.Marshal(map[string]interface{}{
		"model":      ClaudeModel,
		"messages":   messages,
		"max_tokens": 2048, // Increased token limit for detailed feedback
	})

	req, err := http.NewRequest("POST", ClaudeEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("x-api-key", ClaudeAPIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Anthropic-Version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API request failed"})
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	// Parse JSON response
	var responseJSON map[string]interface{}
	json.Unmarshal(body, &responseJSON)

	var extractedJSON string
	if contentArray, found := responseJSON["content"].([]interface{}); found && len(contentArray) > 0 {
		if firstContent, ok := contentArray[0].(map[string]interface{}); ok {
			if text, exists := firstContent["text"].(string); exists {
				extractedJSON = text
			}
		}
	}

	if extractedJSON == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response structure from Claude API"})
		return
	}

	extractedJSON = strings.TrimSpace(extractedJSON)

	if !json.Valid([]byte(extractedJSON)) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Extracted JSON is not valid"})
		return
	}

	type ScaffoldingVariations struct {
		Struggling       string `json:"struggling"`
		Developing       string `json:"developing"`
		NearlyProficient string `json:"nearly_proficient"`
	}

	var variations ScaffoldingVariations
	if err := json.Unmarshal([]byte(extractedJSON), &variations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse scaffolding variations JSON"})
		return
	}

	levels := map[int]string{
		models.ScaffoldingLevelStruggling:       variations.Struggling,
		models.ScaffoldingLevelDeveloping:       variations.Developing,
		models.ScaffoldingLevelNearlyProficient: variations.NearlyProficient,
	}

	var newScaffoldings []models.Scaffolding
	for level, material := range levels {
		if !existingLevels[level] {
			newScaffoldings = append(newScaffoldings, models.Scaffolding{
				ProblemID:           problemId,
				ScaffoldingStrategy: scaffoldingStrategy,
				ScaffoldingLevel:    level,
				ScaffoldingMaterial: material,
				Time:                time.Now(),
			})
		}
	}

	if err := models.DB.Create(&newScaffoldings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert scaffolding records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Scaffolding inserted successfully"})
}

func MakeRequestClaudeAnalyze(c *gin.Context, messages []map[string]string) string {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"model":      ClaudeModel,
		"messages":   messages,
		"max_tokens": 4096,
	})

	req, err := http.NewRequest("POST", ClaudeEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return ""
	}

	// ✅ Use "x-api-key" instead of "Authorization"
	req.Header.Set("x-api-key", ClaudeAPIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Anthropic-Version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API request failed"})
		return ""
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	// Parse JSON response
	var responseJSON map[string]interface{}
	json.Unmarshal(body, &responseJSON)

	// Extract text from the content array
	if contentArray, found := responseJSON["content"].([]interface{}); found && len(contentArray) > 0 {
		if firstContent, ok := contentArray[0].(map[string]interface{}); ok {
			if text, exists := firstContent["text"].(string); exists {
				// ✅ Extract only JSON portion using regex
				re := regexp.MustCompile(`(?s)\{.*\}`)
				jsonPart := re.FindString(text)
				if jsonPart == "" {
					log.Println("Failed to extract JSON from Claude response")
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not extract JSON from Claude response"})
					return ""
				}

				// ✅ Return only the JSON string for unmarshalling later
				return jsonPart
			}
		}
	}

	// Handle unexpected response structure
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from Claude API"})
	return ""
	// todo : save it to db
}

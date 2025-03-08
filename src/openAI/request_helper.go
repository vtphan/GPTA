package openAI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/GPTA/src/repository"
	"io/ioutil"
	"net/http"

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

func makeRequestClaude(c *gin.Context, messages []map[string]string) {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"model":      ClaudeModel,
		"messages":   messages,
		"max_tokens": 1024,
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

func makeRequestClaude2(c *gin.Context, messages []map[string]string, problemId int) {
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

	// Extract text response
	if contentArray, found := responseJSON["content"].([]interface{}); found && len(contentArray) > 0 {
		if firstContent, ok := contentArray[0].(map[string]interface{}); ok {
			if text, exists := firstContent["text"].(string); exists {
				err = repository.SaveClassFeedback(problemId, text)
				if err != nil {
					fmt.Println(err)
				}
				c.JSON(http.StatusOK, gin.H{"response": text})
				return
			}
		}
	}

	// Handle unexpected response structure
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from Claude API"})
	// todo: generate class summary/ get latest class summary
}

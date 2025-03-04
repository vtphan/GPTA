package openAI

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/franciscoescher/goopenai"
	"github.com/gin-gonic/gin"
)

func makeRequest(c *gin.Context, messages []goopenai.Message) {
	apiKey := OpenaiAPIKey
	organization := ""

	client := goopenai.NewClient(apiKey, organization)

	r := goopenai.CreateChatCompletionsRequest{
		Model:       ChatGPTModel,
		Messages:    messages,
		Temperature: 0,
	}
	fmt.Print(messages)
	for i := 0; i < NumRetry; i++ {
		completions, err := client.CreateChatCompletions(context.Background(), r)
		if err != nil {
			panic(err)
		}
		// fmt.Println("Printing Content:")
		// fmt.Print(completions)
		var feedbacksJSON map[string]interface{}
		err = json.Unmarshal([]byte(completions.Choices[0].Message.Content), &feedbacksJSON)
		if err == nil {
			c.JSON(http.StatusOK, feedbacksJSON)
			return
		}
		// fmt.Println(feedbacksJSON)
	}
	c.JSON(http.StatusInternalServerError, "Couldn't process the request!")
}

func makeRequest2(c *gin.Context, messages []goopenai.Message) {
	apiKey := OpenaiAPIKey
	organization := ""

	client := goopenai.NewClient(apiKey, organization)

	r := goopenai.CreateChatCompletionsRequest{
		Model:       ChatGPTModel,
		Messages:    messages,
		Temperature: 0,
	}
	for i := 0; i < NumRetry; i++ {
		completions, err := client.CreateChatCompletions(context.Background(), r)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Return the response in a structured JSON format
		c.JSON(http.StatusOK, gin.H{
			"response": completions.Choices[0].Message.Content,
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't process the request!"})
}

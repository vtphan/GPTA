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

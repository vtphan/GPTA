package openAI

import (
	"encoding/json"
	"fmt"
	"github.com/franciscoescher/goopenai"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InstructionsWithExampleHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := gin.CreateTestContext(w)
	c.Request = r
	EffectiveFeedbackInstructionWithExample(c)
}

func ProcessCodeWithPromptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ProcessCodeWithPrompt(w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func EffectiveFeedbackInstructionWithExample(c *gin.Context) {
	processHandlers(c, getFeedbackWithInstructionsWithExample)
}

func processHandlers(c *gin.Context, f getPromptFunc) {
	var requestBody RequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	requestData := createRequestData(f(requestBody.Course, requestBody.Duration), requestBody)
	makeRequest(c, requestData)
}

func ProcessCodeWithPrompt(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Code         string `json:"code"`
		Instructions string `json:"custom_prompt"`
		UserID       int    `json:"user_id"`
	}

	// Parse the incoming JSON body
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//exampleFeedback := &ExpectedFeedback{
	//	Feedbacks: []Feedback{{
	//		SolutionID: 0,
	//		NumErrors:  2,
	//		HelpNeeded: true,
	//		Feedback:   "Response to Instructions",
	//	}},
	//}
	//jsonExample, _ := json.Marshal(exampleFeedback)
	//jsonExampleStr := string(jsonExample)

	// Construct the prompt for OpenAI
	responseFormat := " Please provide the response strictly in JSON format. The key should be feedback followed by a string value"
	NewInstructions := "With reference to the below code, " + requestData.Instructions + responseFormat

	prompt := fmt.Sprintf(
		"Instructions:  %s\nCode:\n%s\n",
		NewInstructions, requestData.Code,
	)

	// Create OpenAI request message
	messages := []goopenai.Message{
		{Role: "user", Content: prompt},
	}

	// Create a Gin context from http.ResponseWriter and http.Request
	c, _ := gin.CreateTestContext(w)
	c.Request = r

	// Call the existing makeRequest function
	makeRequest(c, messages)
}

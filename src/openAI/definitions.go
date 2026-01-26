package openAI

import (
	"encoding/json"
	"fmt"
	"github.com/franciscoescher/goopenai"
	"log"
)

// NumRetry is the number of retry a request should do.
const NumRetry = 3

// Feedback is one instance of expected feedback from ChatGpt.
type Feedback struct {
	SolutionID int    `json:"solution_id"`
	NumErrors  int    `json:"number_of_errors"`
	HelpNeeded bool   `json:"help_needed"`
	Feedback   string `json:"feedback"`
}

// ExpectedFeedback is an example feedback structure which will be pass
//
//	to the ChatGpt as an example.
type ExpectedFeedback struct {
	Feedbacks []Feedback `json:"feedbacks"`
}

// ChatGPTAPIEndpoint is the API endpoint for ChatGPT. Currently using prompt completion.
const ChatGPTAPIEndpoint = "https://api.openai.com/v1/chat/completions"

// ChatGPTModel is the model/engine using for the feedback.
const ChatGPTModel = "gpt-3.5-turbo"

// OpenaiAPIKey is the key for ChatGPT connection. Should hide before publishing the code.
var OpenaiAPIKey = ""

const generalFeedbackPrompt = "give feedback for each of the solutions"
const hintsPrompt = "give hints to the students to fix problems in the code (if there is any)"
const explanationPrompt = "give an explanation for each of the code mentioning what the student did wrong"
const examplePrompt = "give an example to help the student understand the concept"
const noPointErrorPrompt = "do not point to the error explicitly"
const notGiveawayPrompt = "do not give answer away; make them think"
const encouragePrompt = "encourage the student by pointing out what they did right"

const expectedFormatPrompt = "give the answser as a JSON text as the following format and fields, "

func getInitialPrompt(courseName string, duration int) string {
	exampleFeedback := &ExpectedFeedback{
		Feedbacks: []Feedback{{
			SolutionID: 0,
			NumErrors:  2,
			HelpNeeded: true,
			Feedback:   "This is an example",
		}},
	}
	jsonExample, _ := json.Marshal(exampleFeedback)
	prompt := fmt.Sprintf("Suppose you are the instructor in the course %s "+
		" and you are conducting a %d minutes in-class exercise and the "+
		"students have been provided with a problem / exercise mentioned "+
		"below. All the data is provided into a JSON text including course "+
		"name, exercise duration, problem statement and student solutions "+
		"where each of the solutions contains the solution id, student code "+
		"and number of minutes left when submitted. Remember to take to the "+
		"solution id from the student's solutions. %s\n%s\n", courseName,
		duration, expectedFormatPrompt, string(jsonExample))
	return prompt
}

func getDefaultPrompt(courseName string, duration int) string {
	prompt := getInitialPrompt(courseName, duration)
	prompt += generalFeedbackPrompt
	return prompt
}

func getHintsPrompt(courseName string, duration int) string {
	prompt := getInitialPrompt(courseName, duration)
	prompt += hintsPrompt + " and " + noPointErrorPrompt + " and " + notGiveawayPrompt
	return prompt
}

func getExplanationPrompt(courseName string, duration int) string {
	prompt := getInitialPrompt(courseName, duration)
	prompt += explanationPrompt + " and " + noPointErrorPrompt + " and " + notGiveawayPrompt
	return prompt
}

func getGiveExamplePrompt(courseName string, duration int) string {
	prompt := getInitialPrompt(courseName, duration)
	prompt += examplePrompt + " and " + noPointErrorPrompt + " and " + notGiveawayPrompt
	return prompt
}

func getEcouragementWithHintsPrompt(courseName string, duration int) string {
	prompt := getInitialPrompt(courseName, duration)
	prompt += encouragePrompt + " and " + hintsPrompt + " and " + noPointErrorPrompt + " and " + notGiveawayPrompt
	return prompt
}

func getEcouragementWithExamplePrompt(courseName string, duration int) string {
	prompt := getInitialPrompt(courseName, duration)
	prompt += encouragePrompt + " and " + examplePrompt + " and " + explanationPrompt + " and " + noPointErrorPrompt + " and " + notGiveawayPrompt
	return prompt
}

func getEcouragementWithExplanationPrompt(courseName string, duration int) string {
	prompt := getInitialPrompt(courseName, duration)
	prompt += encouragePrompt + " and " + explanationPrompt + " and " + noPointErrorPrompt + " and " + notGiveawayPrompt
	return prompt
}

func getAllNonDefaultPrompt(courseName string, duration int) string {
	prompt := getInitialPrompt(courseName, duration)
	prompt += encouragePrompt + " and " + hintsPrompt + " and " + explanationPrompt + " and " + examplePrompt + " and " + noPointErrorPrompt + " and " + notGiveawayPrompt
	return prompt
}

func getFeedbackWithInstructionsWithExample(courseName string, duration int) string {
	prompt := getInitialPrompt(courseName, duration)
	prompt += " Please follow the following instructions for your feedback.\n" + EffectiveFeedbackWithExample
	return prompt
}

func getFeedbackWithInstructionsWithoutExample(courseName string, duration int) string {
	prompt := getInitialPrompt(courseName, duration)
	prompt += " Please follow the following instructions for your feedback.\n" + EffectiveFeedbackWithoutExample
	return prompt
}

type getPromptFunc func(string, int) string

// Solution represents each of the code's detail information.
type Solution struct {
	ID         int    `json:"solution_id"`
	Code       string `json:"code"`
	MinuteLeft int    `json:"minute_left"`
}

// RequestBody represents the request body structure
type RequestBody struct {
	Problem   string     `json:"problem"`
	Solutions []Solution `json:"solutions"`
	Course    string     `json:"course"`
	Duration  int        `json:"duration"`
}

func createRequestData(prompt string, request RequestBody) []goopenai.Message {
	b, err := json.Marshal(request)
	if err != nil {
		log.Fatal("Unable to parse struct to json.")
	}
	newPrompt := prompt + "\n" + string(b)
	return []goopenai.Message{{Role: "user", Content: newPrompt}}
}

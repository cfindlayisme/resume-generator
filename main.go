package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/cfindlayisme/resume-generator/env"
	"github.com/sashabaranov/go-openai"
)

//go:embed resume.json jobdescription.txt
var content embed.FS

//go:embed example-format.json
var gptResponseFormat string

// Resume struct to match resume.json structure
type Resume struct {
	Name       string       `json:"name"`
	Email      string       `json:"email"`
	Summary    string       `json:"summary"`
	Skills     []string     `json:"skills"`
	Experience []Experience `json:"experience"`
}

type Experience struct {
	Company  string   `json:"company"`
	Role     string   `json:"role"`
	Duration string   `json:"duration"`
	Details  []string `json:"details"`
}

// Load the resume from the embedded JSON file
func loadResume() (*Resume, error) {
	data, err := content.ReadFile("resume.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read resume.json: %v", err)
	}

	var resume Resume
	err = json.Unmarshal(data, &resume)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal resume JSON: %v", err)
	}

	return &resume, nil
}

// Load the job description from the embedded text file
func loadJobDescription() (string, error) {
	data, err := content.ReadFile("jobdescription.txt")
	if err != nil {
		return "", fmt.Errorf("failed to read jobdescription.txt: %v", err)
	}

	// Convert the content to a string and remove any extra whitespace
	return strings.TrimSpace(string(data)), nil
}

func generateTailoredResume(apiKey, jobDescription string, resume *Resume) (string, error) {
	client := openai.NewClient(apiKey)

	// Build the prompt to send to ChatGPT
	prompt := fmt.Sprintf(`Here's a resume:

Name: %s
Email: %s
Summary: %s
Skills: %v
Experience: %v

The below is an example JSON file format for the resume. Please respond in this format:
%s

Based on the following job description, generate a tailored resume emphasizing relevant skills and experiences. Feel free to remove or rephrase any information as needed.

Job Description:
%s
`, resume.Name, resume.Email, resume.Summary, resume.Skills, resume.Experience, jobDescription, gptResponseFormat)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "system",
					Content: "You are a helpful assistant that rewrites resumes to match job descriptions.",
				},
				{
					Role:    "user",
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("ChatGPT request failed: %v", err)
	}

	return resp.Choices[0].Message.Content, nil
}

func main() {

	err := env.Init()
	if err != nil {
		log.Fatalf("Failed to initialize environment: %v", err)
	}
	apiKey := env.GetOpenAIKey()

	// Load the resume data from the embedded JSON file
	resume, err := loadResume()
	if err != nil {
		log.Fatalf("Failed to load resume: %v", err)
	}

	// Load the job description from the embedded text file
	jobDescription, err := loadJobDescription()
	if err != nil {
		log.Fatalf("Failed to load job description: %v", err)
	}

	// Generate the tailored resume using the OpenAI API
	tailoredResume, err := generateTailoredResume(apiKey, jobDescription, resume)
	if err != nil {
		log.Fatalf("Failed to generate tailored resume: %v", err)
	}

	// Output the tailored resume
	fmt.Println("Tailored Resume:\n")
	fmt.Println(tailoredResume)
}

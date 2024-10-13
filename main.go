package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "embed"

	"github.com/cfindlayisme/resume-generator/env"
	"github.com/cfindlayisme/resume-generator/model"
	"github.com/sashabaranov/go-openai"
)

//go:embed example-format.json
var gptResponseFormat string

func loadResume() (*model.Resume, error) {
	data, err := os.ReadFile("resume.json") // Reading from file system
	if err != nil {
		return nil, fmt.Errorf("failed to read resume.json: %v", err)
	}

	var resume model.Resume
	err = json.Unmarshal(data, &resume)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal resume JSON: %v", err)
	}

	return &resume, nil
}

// Load the job description from the filesystem
func loadJobDescription() (string, error) {
	data, err := os.ReadFile("jobdescription.txt") // Reading from file system
	if err != nil {
		return "", fmt.Errorf("failed to read jobdescription.txt: %v", err)
	}

	// Convert the content to a string and remove any extra whitespace
	return strings.TrimSpace(string(data)), nil
}

func generateTailoredResume(apiKey, jobDescription string, resume *model.Resume) (string, error) {
	client := openai.NewClient(apiKey)

	// Build the prompt to send to ChatGPT
	prompt := fmt.Sprintf(`Here's a resume:

Name: %s
Email: %s
Summary: %s
Skills: %v
Experience: %v

Based on the following job description, generate a tailored resume emphasizing relevant skills and experiences. Feel free to remove or rephrase any information as needed.

Job Description:
%s
`, resume.Name, resume.Email, resume.Summary, resume.Skills, resume.Experience, jobDescription)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "system",
					Content: "You are a helpful assistant that rewrites resumes to match job descriptions.",
				},
				{
					Role:    "system",
					Content: "Respond in this format: " + gptResponseFormat,
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
func generateTailoredCoverLetter(apiKey, jobDescription string, resume *model.Resume) (string, error) {
	client := openai.NewClient(apiKey)

	// Build the prompt to send to ChatGPT
	prompt := fmt.Sprintf(`Here's a resume:

Name: %s
Email: %s
Summary: %s
Skills: %v
Experience: %v

Based on the following job description, generate a tailored cover letter emphasizing relevant skills and experiences.

Job Description:
%s
`, resume.Name, resume.Email, resume.Summary, resume.Skills, resume.Experience, jobDescription)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "system",
					Content: "You are a helpful assistant that writes cover letters to match job descriptions.",
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

	// Load the resume data from the filesystem
	resume, err := loadResume()
	if err != nil {
		log.Fatalf("Failed to load resume: %v", err)
	}

	// Load the job description from the filesystem
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
	var tailoredResumeObj model.Resume
	err = json.Unmarshal([]byte(tailoredResume), &tailoredResumeObj)
	if err != nil {
		log.Fatalf("Failed to unmarshal tailored resume JSON: %v", err)
	}

	// Generate the tailored cover letter using the OpenAI API
	tailoredCoverLetter, err := generateTailoredCoverLetter(apiKey, jobDescription, &tailoredResumeObj)
	if err != nil {
		log.Fatalf("Failed to generate tailored cover letter: %v", err)
	}

	response := model.TailoredResponse{
		TailoredResume:      tailoredResumeObj,
		GeneratedTime:       time.Now(),
		TailoredCoverLetter: model.CoverLetter{Content: tailoredCoverLetter},
	}

	fmt.Println("Tailored response:")
	fmt.Printf("%+v\n", response)
}

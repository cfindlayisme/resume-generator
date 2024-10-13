package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cfindlayisme/resume-generator/env"
	"github.com/cfindlayisme/resume-generator/llm"
	"github.com/cfindlayisme/resume-generator/model"
)

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
	tailoredResume, err := llm.GenerateTailoredResume(apiKey, jobDescription, resume)
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
	tailoredCoverLetter, err := llm.GenerateTailoredCoverLetter(apiKey, jobDescription, &tailoredResumeObj)
	if err != nil {
		log.Fatalf("Failed to generate tailored cover letter: %v", err)
	}

	response := model.TailoredResponse{
		TailoredResume:      tailoredResumeObj,
		GeneratedTime:       time.Now(),
		TailoredCoverLetter: model.CoverLetter{Content: tailoredCoverLetter},
	}

	fmt.Println("Tailored response:")
	prettyResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal tailored response: %v", err)
	}
	fmt.Println(string(prettyResponse))
}

package llm

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/cfindlayisme/resume-generator/model"
	"github.com/sashabaranov/go-openai"
)

//go:embed example-format.json
var gptResponseFormat string

func GenerateTailoredResume(apiKey, jobDescription string, resume *model.Resume) (string, error) {
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
func GenerateTailoredCoverLetter(apiKey, jobDescription string, resume *model.Resume) (string, error) {
	client := openai.NewClient(apiKey)

	// Build the prompt to send to ChatGPT
	prompt := fmt.Sprintf(`Here's a resume:

Name: %s
Email: %s
Summary: %s
Skills: %v
Experience: %v

Based on the following job description, generate a tailored cover letter emphasizing relevant skills and experiences. Do not add information to replace - such as address, date, etc. I only want the Dear Hiring Manager, Introduction, Body, and Closing."

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

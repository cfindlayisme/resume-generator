package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
)

//go:embed resume.json
var resumeFile embed.FS

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

func main() {
	data, err := resumeFile.ReadFile("resume.json")
	if err != nil {
		log.Fatalf("Failed to read resume.json: %v", err)
	}

	var resume Resume
	err = json.Unmarshal(data, &resume)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	fmt.Printf("Resume of %s:\n", resume.Name)
	fmt.Printf("Email: %s\n", resume.Email)
	fmt.Printf("Summary: %s\n", resume.Summary)
	fmt.Println("Skills:")
	for _, skill := range resume.Skills {
		fmt.Printf(" - %s\n", skill)
	}

	fmt.Println("Experience:")
	for _, exp := range resume.Experience {
		fmt.Printf("Company: %s, Role: %s (%s)\n", exp.Company, exp.Role, exp.Duration)
		for _, detail := range exp.Details {
			fmt.Printf("  - %s\n", detail)
		}
	}
}

package model

import "time"

type TailoredResponse struct {
	TailoredResume      Resume      `json:"tailoredResume"`
	GeneratedTime       time.Time   `json:"generatedTime"`
	TailoredCoverLetter CoverLetter `json:"tailoredCoverLetter"`
}

// Resume struct to match resume.json structure
type Resume struct {
	Name       string       `json:"name"`
	Email      string       `json:"email"`
	Summary    string       `json:"summary"`
	Skills     []string     `json:"skills"`
	Experience []Experience `json:"experience"`
}

type CoverLetter struct {
	Content string `json:"content"`
}

type Experience struct {
	Company  string   `json:"company"`
	Role     string   `json:"role"`
	Duration string   `json:"duration"`
	Details  []string `json:"details"`
}

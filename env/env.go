package env

import (
	"fmt"
	"os"
)

func Init() error {
	if os.Getenv("OPENAI_API_KEY") == "" {
		return fmt.Errorf("environment variable OPENAI_API_KEY is not set")
	}
	return nil
}

func GetOpenAIKey() string {
	return os.Getenv("OPENAI_API_KEY")
}

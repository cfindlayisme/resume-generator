package env_test

import (
	"os"
	"testing"

	"github.com/cfindlayisme/resume-generator/env"
	"github.com/go-playground/assert"
)

func Test_GetOpenAIKeyReturnsEmptyIfUnset(t *testing.T) {
	os.Unsetenv("OPENAI_API_KEY")

	assert.Equal(t, env.GetOpenAIKey(), "")
}

func Test_GetOpenAIKeyGetsEnvironmentVariable(t *testing.T) {
	os.Unsetenv("OPENAI_API_KEY")
	os.Setenv("OPENAI_API_KEY", "mockapikey")

	assert.Equal(t, env.GetOpenAIKey(), "mockapikey")
}

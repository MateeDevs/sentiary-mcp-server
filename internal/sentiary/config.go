package sentiary

import (
	"os"
	"strings"
)

const defaultBaseURL = "https://api.sentiary.com"

type Config struct {
	BaseURL          string
	APIKey           string
	DefaultProjectID string
}

func ConfigFromEnv() Config {
	return Config{
		BaseURL:          defaultBaseURL,
		APIKey:           firstEnv("SENTIARY_USER_API_KEY", ""),
		DefaultProjectID: firstEnv("SENTIARY_PROJECT_ID", ""),
	}
}

func firstEnv(keys ...string) string {
	for _, key := range keys {
		value := strings.TrimSpace(os.Getenv(key))
		if value != "" {
			return value
		}
	}
	return ""
}

package config

import "os"
import "strings"

func loadEnvironment() envValues {
	ev := envValues{}

	if val := os.Getenv("AICLI_PROTOCOL"); val != "" {
		ev.protocol = val
	}
	if val := os.Getenv("AICLI_URL"); val != "" {
		ev.url = val
	}
	if val := os.Getenv("AICLI_API_KEY"); val != "" {
		ev.key = val
	} else if val := os.Getenv("AICLI_API_KEY_FILE"); val != "" {
		content, err := os.ReadFile(val)
		if err == nil {
			ev.key = strings.TrimSpace(string(content))
		}
	}
	if val := os.Getenv("AICLI_MODEL"); val != "" {
		ev.model = val
	}
	if val := os.Getenv("AICLI_FALLBACK"); val != "" {
		ev.fallback = val
	}
	if val := os.Getenv("AICLI_SYSTEM"); val != "" {
		ev.system = val
	}

	return ev
}

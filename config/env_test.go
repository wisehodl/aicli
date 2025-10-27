package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadEnvironment(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want envValues
	}{
		{
			name: "empty environment",
			env:  map[string]string{},
			want: envValues{},
		},
		{
			name: "protocol only",
			env:  map[string]string{"AICLI_PROTOCOL": "ollama"},
			want: envValues{protocol: "ollama"},
		},
		{
			name: "url only",
			env:  map[string]string{"AICLI_URL": "http://localhost:11434"},
			want: envValues{url: "http://localhost:11434"},
		},
		{
			name: "api key direct",
			env:  map[string]string{"AICLI_API_KEY": "sk-test123"},
			want: envValues{key: "sk-test123"},
		},
		{
			name: "model only",
			env:  map[string]string{"AICLI_MODEL": "llama3"},
			want: envValues{model: "llama3"},
		},
		{
			name: "fallback only",
			env:  map[string]string{"AICLI_FALLBACK": "gpt-3.5,gpt-4"},
			want: envValues{fallback: "gpt-3.5,gpt-4"},
		},
		{
			name: "system only",
			env:  map[string]string{"AICLI_SYSTEM": "You are helpful"},
			want: envValues{system: "You are helpful"},
		},
		{
			name: "all variables set",
			env: map[string]string{
				"AICLI_PROTOCOL": "openai",
				"AICLI_URL":      "https://api.openai.com/v1/chat/completions",
				"AICLI_API_KEY":  "sk-abc",
				"AICLI_MODEL":    "gpt-4",
				"AICLI_FALLBACK": "gpt-3.5",
				"AICLI_SYSTEM":   "system prompt",
			},
			want: envValues{
				protocol: "openai",
				url:      "https://api.openai.com/v1/chat/completions",
				key:      "sk-abc",
				model:    "gpt-4",
				fallback: "gpt-3.5",
				system:   "system prompt",
			},
		},
		{
			name: "empty string values preserved",
			env:  map[string]string{"AICLI_SYSTEM": ""},
			want: envValues{system: ""},
		},
		{
			name: "whitespace preserved",
			env:  map[string]string{"AICLI_SYSTEM": "  spaces  "},
			want: envValues{system: "  spaces  "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			got := loadEnvironment()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLoadEnvironmentKeyFile(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want envValues
	}{
		{
			name: "key file when no direct key",
			env:  map[string]string{"AICLI_API_KEY_FILE": "testdata/api.key"},
			want: envValues{key: "sk-test-key-123"},
		},
		{
			name: "direct key overrides key file",
			env: map[string]string{
				"AICLI_API_KEY":      "sk-direct",
				"AICLI_API_KEY_FILE": "testdata/api.key",
			},
			want: envValues{key: "sk-direct"},
		},
		{
			name: "key file not found",
			env:  map[string]string{"AICLI_API_KEY_FILE": "/nonexistent/key.txt"},
			want: envValues{},
		},
		{
			name: "key file with whitespace trimmed",
			env:  map[string]string{"AICLI_API_KEY_FILE": "testdata/api_whitespace.key"},
			want: envValues{key: "sk-whitespace-key"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			got := loadEnvironment()
			assert.Equal(t, tt.want, got)
		})
	}
}

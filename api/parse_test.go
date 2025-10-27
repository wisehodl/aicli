package api

import (
	"os"
	"testing"

	"git.wisehodl.dev/jay/aicli/config"
	"github.com/stretchr/testify/assert"
)

func TestParseResponse(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		protocol    config.APIProtocol
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name:     "openai success",
			body:     `{"choices":[{"message":{"content":"This is the response text."}}]}`,
			protocol: config.ProtocolOpenAI,
			want:     "This is the response text.",
		},
		{
			name:        "openai empty choices",
			body:        `{"choices":[]}`,
			protocol:    config.ProtocolOpenAI,
			wantErr:     true,
			errContains: "empty choices array",
		},
		{
			name:        "openai no choices field",
			body:        `{"result":"ok"}`,
			protocol:    config.ProtocolOpenAI,
			wantErr:     true,
			errContains: "no choices in response",
		},
		{
			name:        "openai invalid choice format",
			body:        `{"choices":["invalid"]}`,
			protocol:    config.ProtocolOpenAI,
			wantErr:     true,
			errContains: "invalid choice format",
		},
		{
			name:        "openai no message field",
			body:        `{"choices":[{"text":"wrong structure"}]}`,
			protocol:    config.ProtocolOpenAI,
			wantErr:     true,
			errContains: "no message in choice",
		},
		{
			name:        "openai no content field",
			body:        `{"choices":[{"message":{"role":"assistant"}}]}`,
			protocol:    config.ProtocolOpenAI,
			wantErr:     true,
			errContains: "no content in message",
		},
		{
			name:     "ollama success",
			body:     `{"response":"This is the Ollama response."}`,
			protocol: config.ProtocolOllama,
			want:     "This is the Ollama response.",
		},
		{
			name:        "ollama no response field",
			body:        `{"model":"llama3"}`,
			protocol:    config.ProtocolOllama,
			wantErr:     true,
			errContains: "no response field in ollama response",
		},
		{
			name:        "malformed json",
			body:        `{invalid json`,
			protocol:    config.ProtocolOpenAI,
			wantErr:     true,
			errContains: "parse response",
		},
		{
			name:     "openai multiline content",
			body:     `{"choices":[{"message":{"content":"Line 1\nLine 2\nLine 3"}}]}`,
			protocol: config.ProtocolOpenAI,
			want:     "Line 1\nLine 2\nLine 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseResponse([]byte(tt.body), tt.protocol)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseResponseWithTestdata(t *testing.T) {
	tests := []struct {
		name        string
		file        string
		protocol    config.APIProtocol
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name:     "openai success from file",
			file:     "testdata/openai_success.json",
			protocol: config.ProtocolOpenAI,
			want:     "This is the response text.",
		},
		{
			name:        "openai empty choices from file",
			file:        "testdata/openai_empty_choices.json",
			protocol:    config.ProtocolOpenAI,
			wantErr:     true,
			errContains: "empty choices array",
		},
		{
			name:     "ollama success from file",
			file:     "testdata/ollama_success.json",
			protocol: config.ProtocolOllama,
			want:     "This is the Ollama response.",
		},
		{
			name:        "ollama no response from file",
			file:        "testdata/ollama_no_response.json",
			protocol:    config.ProtocolOllama,
			wantErr:     true,
			errContains: "no response field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := os.ReadFile(tt.file)
			assert.NoError(t, err, "failed to read test file")

			got, err := parseResponse(body, tt.protocol)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

package api

import (
	"testing"

	"git.wisehodl.dev/jay/aicli/config"
	"github.com/stretchr/testify/assert"
)

func TestBuildPayload(t *testing.T) {
	tests := []struct {
		name  string
		cfg   config.ConfigData
		model string
		query string
		want  map[string]interface{}
	}{
		{
			name: "openai without system prompt",
			cfg: config.ConfigData{
				Protocol: config.ProtocolOpenAI,
			},
			model: "gpt-4",
			query: "analyze this",
			want: map[string]interface{}{
				"model": "gpt-4",
				"messages": []map[string]string{
					{"role": "user", "content": "analyze this"},
				},
			},
		},
		{
			name: "openai with system prompt",
			cfg: config.ConfigData{
				Protocol:     config.ProtocolOpenAI,
				SystemPrompt: "You are helpful",
			},
			model: "gpt-4",
			query: "analyze this",
			want: map[string]interface{}{
				"model": "gpt-4",
				"messages": []map[string]string{
					{"role": "system", "content": "You are helpful"},
					{"role": "user", "content": "analyze this"},
				},
			},
		},
		{
			name: "ollama without system prompt",
			cfg: config.ConfigData{
				Protocol: config.ProtocolOllama,
			},
			model: "llama3",
			query: "analyze this",
			want: map[string]interface{}{
				"model":  "llama3",
				"prompt": "analyze this",
				"stream": false,
			},
		},
		{
			name: "ollama with system prompt",
			cfg: config.ConfigData{
				Protocol:     config.ProtocolOllama,
				SystemPrompt: "You are helpful",
			},
			model: "llama3",
			query: "analyze this",
			want: map[string]interface{}{
				"model":  "llama3",
				"prompt": "analyze this",
				"system": "You are helpful",
				"stream": false,
			},
		},
		{
			name: "empty query",
			cfg: config.ConfigData{
				Protocol: config.ProtocolOpenAI,
			},
			model: "gpt-4",
			query: "",
			want: map[string]interface{}{
				"model": "gpt-4",
				"messages": []map[string]string{
					{"role": "user", "content": ""},
				},
			},
		},
		{
			name: "multiline query",
			cfg: config.ConfigData{
				Protocol: config.ProtocolOpenAI,
			},
			model: "gpt-4",
			query: "line1\nline2\nline3",
			want: map[string]interface{}{
				"model": "gpt-4",
				"messages": []map[string]string{
					{"role": "user", "content": "line1\nline2\nline3"},
				},
			},
		},
		{
			name: "model name injection",
			cfg: config.ConfigData{
				Protocol: config.ProtocolOpenAI,
			},
			model: "custom-model-name",
			query: "test",
			want: map[string]interface{}{
				"model": "custom-model-name",
				"messages": []map[string]string{
					{"role": "user", "content": "test"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildPayload(tt.cfg, tt.model, tt.query)
			assert.Equal(t, tt.want, got)
		})
	}
}

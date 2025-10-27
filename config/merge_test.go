package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMergeSources(t *testing.T) {
	tests := []struct {
		name  string
		flags flagValues
		env   envValues
		file  fileValues
		want  ConfigData
	}{
		{
			name:  "all empty uses defaults",
			flags: flagValues{},
			env:   envValues{},
			file:  fileValues{},
			want:  defaultConfig,
		},
		{
			name:  "file overrides defaults",
			flags: flagValues{},
			env:   envValues{},
			file: fileValues{
				protocol: "ollama",
				model:    "llama3",
			},
			want: ConfigData{
				Protocol:       ProtocolOllama,
				URL:            "https://api.ppq.ai/chat/completions",
				Model:          "llama3",
				FallbackModels: []string{"gpt-4.1-mini"},
			},
		},
		{
			name:  "env overrides file",
			flags: flagValues{},
			env: envValues{
				model: "gpt-4",
			},
			file: fileValues{
				model: "llama3",
			},
			want: ConfigData{
				Protocol:       ProtocolOpenAI,
				URL:            "https://api.ppq.ai/chat/completions",
				Model:          "gpt-4",
				FallbackModels: []string{"gpt-4.1-mini"},
			},
		},
		{
			name: "flags override env",
			flags: flagValues{
				model: "claude-3",
			},
			env: envValues{
				model: "gpt-4",
			},
			file: fileValues{},
			want: ConfigData{
				Protocol:       ProtocolOpenAI,
				URL:            "https://api.ppq.ai/chat/completions",
				Model:          "claude-3",
				FallbackModels: []string{"gpt-4.1-mini"},
			},
		},
		{
			name: "full precedence chain",
			flags: flagValues{
				protocol: "ollama",
				quiet:    true,
			},
			env: envValues{
				protocol: "openai",
				model:    "gpt-4",
				url:      "http://custom.api",
			},
			file: fileValues{
				protocol: "openai",
				model:    "llama3",
				url:      "http://file.api",
				fallback: "mistral",
			},
			want: ConfigData{
				Protocol:       ProtocolOllama,
				URL:            "http://custom.api",
				Model:          "gpt-4",
				FallbackModels: []string{"mistral"},
				Quiet:          true,
			},
		},
		{
			name: "fallback string split",
			flags: flagValues{
				fallback: "model1,model2,model3",
			},
			env:  envValues{},
			file: fileValues{},
			want: ConfigData{
				Protocol:       ProtocolOpenAI,
				URL:            "https://api.ppq.ai/chat/completions",
				Model:          "gpt-4o-mini",
				FallbackModels: []string{"model1", "model2", "model3"},
			},
		},
		{
			name: "direct key flag",
			flags: flagValues{
				key: "sk-direct",
			},
			env:  envValues{},
			file: fileValues{},
			want: ConfigData{
				Protocol:       ProtocolOpenAI,
				URL:            "https://api.ppq.ai/chat/completions",
				Model:          "gpt-4o-mini",
				FallbackModels: []string{"gpt-4.1-mini"},
				APIKey:         "sk-direct",
			},
		},
		{
			name: "direct system flag",
			flags: flagValues{
				system: "You are helpful",
			},
			env:  envValues{},
			file: fileValues{},
			want: ConfigData{
				Protocol:       ProtocolOpenAI,
				URL:            "https://api.ppq.ai/chat/completions",
				Model:          "gpt-4o-mini",
				FallbackModels: []string{"gpt-4.1-mini"},
				SystemPrompt:   "You are helpful",
			},
		},
		{
			name: "file paths collected",
			flags: flagValues{
				files:      []string{"a.go", "b.go"},
				prompts:    []string{"prompt1", "prompt2"},
				promptFile: "prompt.txt",
			},
			env:  envValues{},
			file: fileValues{},
			want: ConfigData{
				Protocol:       ProtocolOpenAI,
				URL:            "https://api.ppq.ai/chat/completions",
				Model:          "gpt-4o-mini",
				FallbackModels: []string{"gpt-4.1-mini"},
				FilePaths:      []string{"a.go", "b.go"},
				PromptFlags:    []string{"prompt1", "prompt2"},
				PromptPaths:    []string{"prompt.txt"},
			},
		},
		{
			name: "stdin file flag",
			flags: flagValues{
				stdinFile: true,
			},
			env:  envValues{},
			file: fileValues{},
			want: ConfigData{
				Protocol:       ProtocolOpenAI,
				URL:            "https://api.ppq.ai/chat/completions",
				Model:          "gpt-4o-mini",
				FallbackModels: []string{"gpt-4.1-mini"},
				StdinAsFile:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeSources(tt.flags, tt.env, tt.file)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMergeSourcesKeyFile(t *testing.T) {
	tests := []struct {
		name  string
		flags flagValues
		env   envValues
		file  fileValues
		want  ConfigData
	}{
		{
			name: "key file from flags",
			flags: flagValues{
				keyFile: "testdata/api.key",
			},
			env:  envValues{},
			file: fileValues{},
			want: ConfigData{
				Protocol:       ProtocolOpenAI,
				URL:            "https://api.ppq.ai/chat/completions",
				Model:          "gpt-4o-mini",
				FallbackModels: []string{"gpt-4.1-mini"},
				APIKey:         "sk-test-key-123",
			},
		},
		{
			name:  "key file from file config",
			flags: flagValues{},
			env:   envValues{},
			file: fileValues{
				keyFile: "testdata/api.key",
			},
			want: ConfigData{
				Protocol:       ProtocolOpenAI,
				URL:            "https://api.ppq.ai/chat/completions",
				Model:          "gpt-4o-mini",
				FallbackModels: []string{"gpt-4.1-mini"},
				APIKey:         "sk-test-key-123",
			},
		},
		{
			name: "direct key overrides key file",
			flags: flagValues{
				key:     "sk-direct",
				keyFile: "testdata/api.key",
			},
			env:  envValues{},
			file: fileValues{},
			want: ConfigData{
				Protocol:       ProtocolOpenAI,
				URL:            "https://api.ppq.ai/chat/completions",
				Model:          "gpt-4o-mini",
				FallbackModels: []string{"gpt-4.1-mini"},
				APIKey:         "sk-direct",
			},
		},
		{
			name:  "env key overrides file key file",
			flags: flagValues{},
			env: envValues{
				key: "sk-env",
			},
			file: fileValues{
				keyFile: "testdata/api.key",
			},
			want: ConfigData{
				Protocol:       ProtocolOpenAI,
				URL:            "https://api.ppq.ai/chat/completions",
				Model:          "gpt-4o-mini",
				FallbackModels: []string{"gpt-4.1-mini"},
				APIKey:         "sk-env",
			},
		},
		{
			name: "key file with whitespace trimmed",
			flags: flagValues{
				keyFile: "testdata/api_whitespace.key",
			},
			env:  envValues{},
			file: fileValues{},
			want: ConfigData{
				Protocol:       ProtocolOpenAI,
				URL:            "https://api.ppq.ai/chat/completions",
				Model:          "gpt-4o-mini",
				FallbackModels: []string{"gpt-4.1-mini"},
				APIKey:         "sk-whitespace-key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeSources(tt.flags, tt.env, tt.file)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMergeSourcesSystemFile(t *testing.T) {
	tests := []struct {
		name  string
		flags flagValues
		env   envValues
		file  fileValues
		want  ConfigData
	}{
		{
			name: "system file from flags",
			flags: flagValues{
				systemFile: "testdata/system.txt",
			},
			env:  envValues{},
			file: fileValues{},
			want: ConfigData{
				Protocol:       ProtocolOpenAI,
				URL:            "https://api.ppq.ai/chat/completions",
				Model:          "gpt-4o-mini",
				FallbackModels: []string{"gpt-4.1-mini"},
				SystemPrompt:   "You are a helpful assistant.",
			},
		},
		{
			name:  "system file from file config",
			flags: flagValues{},
			env:   envValues{},
			file: fileValues{
				systemFile: "testdata/system.txt",
			},
			want: ConfigData{
				Protocol:       ProtocolOpenAI,
				URL:            "https://api.ppq.ai/chat/completions",
				Model:          "gpt-4o-mini",
				FallbackModels: []string{"gpt-4.1-mini"},
				SystemPrompt:   "You are a helpful assistant.",
			},
		},
		{
			name: "direct system overrides system file",
			flags: flagValues{
				system:     "Direct system",
				systemFile: "testdata/system.txt",
			},
			env:  envValues{},
			file: fileValues{},
			want: ConfigData{
				Protocol:       ProtocolOpenAI,
				URL:            "https://api.ppq.ai/chat/completions",
				Model:          "gpt-4o-mini",
				FallbackModels: []string{"gpt-4.1-mini"},
				SystemPrompt:   "Direct system",
			},
		},
		{
			name:  "env system overrides file system file",
			flags: flagValues{},
			env: envValues{
				system: "System from env",
			},
			file: fileValues{
				systemFile: "testdata/system.txt",
			},
			want: ConfigData{
				Protocol:       ProtocolOpenAI,
				URL:            "https://api.ppq.ai/chat/completions",
				Model:          "gpt-4o-mini",
				FallbackModels: []string{"gpt-4.1-mini"},
				SystemPrompt:   "System from env",
			},
		},
		{
			name: "empty system file",
			flags: flagValues{
				systemFile: "testdata/system_empty.txt",
			},
			env:  envValues{},
			file: fileValues{},
			want: ConfigData{
				Protocol:       ProtocolOpenAI,
				URL:            "https://api.ppq.ai/chat/completions",
				Model:          "gpt-4o-mini",
				FallbackModels: []string{"gpt-4.1-mini"},
				SystemPrompt:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeSources(tt.flags, tt.env, tt.file)
			assert.Equal(t, tt.want, got)
		})
	}
}

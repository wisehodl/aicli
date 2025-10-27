package config

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want flagValues
	}{
		{
			name: "empty args",
			args: []string{},
			want: flagValues{},
		},
		{
			name: "single file short flag",
			args: []string{"-f", "main.go"},
			want: flagValues{files: []string{"main.go"}},
		},
		{
			name: "single file long flag",
			args: []string{"--file", "main.go"},
			want: flagValues{files: []string{"main.go"}},
		},
		{
			name: "multiple files",
			args: []string{"-f", "a.go", "-f", "b.go", "--file", "c.go"},
			want: flagValues{files: []string{"a.go", "b.go", "c.go"}},
		},
		{
			name: "single prompt short flag",
			args: []string{"-p", "analyze this"},
			want: flagValues{prompts: []string{"analyze this"}},
		},
		{
			name: "single prompt long flag",
			args: []string{"--prompt", "analyze this"},
			want: flagValues{prompts: []string{"analyze this"}},
		},
		{
			name: "multiple prompts",
			args: []string{"-p", "first", "-p", "second", "--prompt", "third"},
			want: flagValues{prompts: []string{"first", "second", "third"}},
		},
		{
			name: "prompt file",
			args: []string{"-pf", "prompt.txt"},
			want: flagValues{promptFile: "prompt.txt"},
		},
		{
			name: "prompt file long",
			args: []string{"--prompt-file", "prompt.txt"},
			want: flagValues{promptFile: "prompt.txt"},
		},
		{
			name: "system short",
			args: []string{"-s", "You are helpful"},
			want: flagValues{system: "You are helpful"},
		},
		{
			name: "system long",
			args: []string{"--system", "You are helpful"},
			want: flagValues{system: "You are helpful"},
		},
		{
			name: "system file short",
			args: []string{"-sf", "system.txt"},
			want: flagValues{systemFile: "system.txt"},
		},
		{
			name: "system file long",
			args: []string{"--system-file", "system.txt"},
			want: flagValues{systemFile: "system.txt"},
		},
		{
			name: "key short",
			args: []string{"-k", "sk-abc123"},
			want: flagValues{key: "sk-abc123"},
		},
		{
			name: "key long",
			args: []string{"--key", "sk-abc123"},
			want: flagValues{key: "sk-abc123"},
		},
		{
			name: "key file short",
			args: []string{"-kf", "api.key"},
			want: flagValues{keyFile: "api.key"},
		},
		{
			name: "key file long",
			args: []string{"--key-file", "api.key"},
			want: flagValues{keyFile: "api.key"},
		},
		{
			name: "protocol short",
			args: []string{"-l", "ollama"},
			want: flagValues{protocol: "ollama"},
		},
		{
			name: "protocol long",
			args: []string{"--protocol", "ollama"},
			want: flagValues{protocol: "ollama"},
		},
		{
			name: "url short",
			args: []string{"-u", "http://localhost:11434"},
			want: flagValues{url: "http://localhost:11434"},
		},
		{
			name: "url long",
			args: []string{"--url", "http://localhost:11434"},
			want: flagValues{url: "http://localhost:11434"},
		},
		{
			name: "model short",
			args: []string{"-m", "gpt-4"},
			want: flagValues{model: "gpt-4"},
		},
		{
			name: "model long",
			args: []string{"--model", "gpt-4"},
			want: flagValues{model: "gpt-4"},
		},
		{
			name: "fallback short",
			args: []string{"-b", "gpt-3.5-turbo"},
			want: flagValues{fallback: "gpt-3.5-turbo"},
		},
		{
			name: "fallback long",
			args: []string{"--fallback", "gpt-3.5-turbo"},
			want: flagValues{fallback: "gpt-3.5-turbo"},
		},
		{
			name: "output short",
			args: []string{"-o", "result.txt"},
			want: flagValues{output: "result.txt"},
		},
		{
			name: "output long",
			args: []string{"--output", "result.txt"},
			want: flagValues{output: "result.txt"},
		},
		{
			name: "config short",
			args: []string{"-c", "config.yaml"},
			want: flagValues{config: "config.yaml"},
		},
		{
			name: "config long",
			args: []string{"--config", "config.yaml"},
			want: flagValues{config: "config.yaml"},
		},
		{
			name: "stdin file short",
			args: []string{"-F"},
			want: flagValues{stdinFile: true},
		},
		{
			name: "stdin file long",
			args: []string{"--stdin-file"},
			want: flagValues{stdinFile: true},
		},
		{
			name: "quiet short",
			args: []string{"-q"},
			want: flagValues{quiet: true},
		},
		{
			name: "quiet long",
			args: []string{"--quiet"},
			want: flagValues{quiet: true},
		},
		{
			name: "verbose short",
			args: []string{"-v"},
			want: flagValues{verbose: true},
		},
		{
			name: "verbose long",
			args: []string{"--verbose"},
			want: flagValues{verbose: true},
		},
		{
			name: "version flag",
			args: []string{"--version"},
			want: flagValues{version: true},
		},
		{
			name: "complex combination",
			args: []string{
				"-f", "a.go",
				"-f", "b.go",
				"-p", "first prompt",
				"-pf", "prompt.txt",
				"-s", "system prompt",
				"-k", "key123",
				"-m", "gpt-4",
				"-b", "gpt-3.5",
				"-o", "out.txt",
				"-q",
				"-v",
			},
			want: flagValues{
				files:      []string{"a.go", "b.go"},
				prompts:    []string{"first prompt"},
				promptFile: "prompt.txt",
				system:     "system prompt",
				key:        "key123",
				model:      "gpt-4",
				fallback:   "gpt-3.5",
				output:     "out.txt",
				quiet:      true,
				verbose:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFlags()

			got, err := parseFlags(tt.args)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseFlagsErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "unknown flag",
			args: []string{"--unknown"},
		},
		{
			name: "flag without value",
			args: []string{"-f"},
		},
		{
			name: "model without value",
			args: []string{"-m"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFlags()

			_, err := parseFlags(tt.args)
			assert.Error(t, err, "parseFlags() should return error for %s", tt.name)
		})
	}
}

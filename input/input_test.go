package input

import (
	"testing"

	"git.wisehodl.dev/jay/aicli/config"
	"github.com/stretchr/testify/assert"
)

func TestResolveInputs(t *testing.T) {
	tests := []struct {
		name         string
		cfg          config.ConfigData
		stdinContent string
		hasStdin     bool
		want         InputData
		wantErr      bool
		errContains  string
	}{
		{
			name:         "no input returns error",
			cfg:          config.ConfigData{},
			stdinContent: "",
			hasStdin:     false,
			wantErr:      true,
			errContains:  "no input provided",
		},
		{
			name:         "stdin only as prompt",
			cfg:          config.ConfigData{},
			stdinContent: "analyze this",
			hasStdin:     true,
			want: InputData{
				Prompts: []string{"analyze this"},
				Files:   []FileData{},
			},
		},
		{
			name: "prompt flag only",
			cfg: config.ConfigData{
				PromptFlags: []string{"test prompt"},
			},
			stdinContent: "",
			hasStdin:     false,
			want: InputData{
				Prompts: []string{"test prompt"},
				Files:   []FileData{},
			},
		},
		{
			name: "file flag only",
			cfg: config.ConfigData{
				FilePaths: []string{"testdata/code.go"},
			},
			stdinContent: "",
			hasStdin:     false,
			want: InputData{
				Prompts: []string{},
				Files: []FileData{
					{Path: "testdata/code.go", Content: "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}\n"},
				},
			},
		},
		{
			name: "stdin as file with -F flag",
			cfg: config.ConfigData{
				StdinAsFile: true,
			},
			stdinContent: "stdin content",
			hasStdin:     true,
			want: InputData{
				Prompts: []string{},
				Files: []FileData{
					{Path: "input", Content: "stdin content"},
				},
			},
		},
		{
			name: "stdin as file with -F and explicit files",
			cfg: config.ConfigData{
				StdinAsFile: true,
				FilePaths:   []string{"testdata/code.go"},
			},
			stdinContent: "stdin content",
			hasStdin:     true,
			want: InputData{
				Prompts: []string{},
				Files: []FileData{
					{Path: "input", Content: "stdin content"},
					{Path: "testdata/code.go", Content: "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}\n"},
				},
			},
		},
		{
			name: "stdin prefixed with explicit prompt",
			cfg: config.ConfigData{
				PromptFlags: []string{"analyze"},
			},
			stdinContent: "code to analyze",
			hasStdin:     true,
			want: InputData{
				Prompts: []string{"analyze", "code to analyze"},
				Files:   []FileData{},
			},
		},
		{
			name: "prompt from file",
			cfg: config.ConfigData{
				PromptPaths: []string{"testdata/prompt1.txt"},
			},
			stdinContent: "",
			hasStdin:     false,
			want: InputData{
				Prompts: []string{"Analyze the following code.\n"},
				Files:   []FileData{},
			},
		},
		{
			name: "complete scenario: prompts, files, stdin",
			cfg: config.ConfigData{
				PromptFlags: []string{"review this"},
				FilePaths:   []string{"testdata/code.go"},
			},
			stdinContent: "additional context",
			hasStdin:     true,
			want: InputData{
				Prompts: []string{"review this", "additional context"},
				Files: []FileData{
					{Path: "testdata/code.go", Content: "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}\n"},
				},
			},
		},
		{
			name: "file read error propagates",
			cfg: config.ConfigData{
				FilePaths: []string{"testdata/nonexistent.go"},
			},
			stdinContent: "",
			hasStdin:     false,
			wantErr:      true,
			errContains:  "read file",
		},
		{
			name: "prompt file read error propagates",
			cfg: config.ConfigData{
				PromptPaths: []string{"testdata/missing.txt"},
			},
			stdinContent: "",
			hasStdin:     false,
			wantErr:      true,
			errContains:  "read prompt file",
		},
		{
			name: "empty file path error propagates",
			cfg: config.ConfigData{
				FilePaths: []string{""},
			},
			stdinContent: "",
			hasStdin:     false,
			wantErr:      true,
			errContains:  "empty file path",
		},
		{
			name: "stdin replaces prompts when no explicit flags",
			cfg: config.ConfigData{
				PromptFlags: []string{},
			},
			stdinContent: "stdin prompt",
			hasStdin:     true,
			want: InputData{
				Prompts: []string{"stdin prompt"},
				Files:   []FileData{},
			},
		},
		{
			name: "multiple files in order",
			cfg: config.ConfigData{
				FilePaths: []string{"testdata/code.go", "testdata/data.json"},
			},
			stdinContent: "",
			hasStdin:     false,
			want: InputData{
				Prompts: []string{},
				Files: []FileData{
					{Path: "testdata/code.go", Content: "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}\n"},
					{Path: "testdata/data.json", Content: "{\n  \"name\": \"test\",\n  \"value\": 42\n}\n"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveInputs(tt.cfg, tt.stdinContent, tt.hasStdin)
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

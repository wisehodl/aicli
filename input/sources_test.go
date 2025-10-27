package input

import (
	"testing"

	"git.wisehodl.dev/jay/aicli/config"
	"github.com/stretchr/testify/assert"
)

func TestReadPromptSources(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.ConfigData
		want    []string
		wantErr bool
	}{
		{
			name: "no prompts returns empty",
			cfg:  config.ConfigData{},
			want: []string{},
		},
		{
			name: "single flag prompt",
			cfg: config.ConfigData{
				PromptFlags: []string{"analyze this"},
			},
			want: []string{"analyze this"},
		},
		{
			name: "multiple flag prompts",
			cfg: config.ConfigData{
				PromptFlags: []string{"first", "second", "third"},
			},
			want: []string{"first", "second", "third"},
		},
		{
			name: "single prompt file",
			cfg: config.ConfigData{
				PromptPaths: []string{"testdata/prompt1.txt"},
			},
			want: []string{"Analyze the following code.\n"},
		},
		{
			name: "multiple prompt files",
			cfg: config.ConfigData{
				PromptPaths: []string{"testdata/prompt1.txt", "testdata/prompt2.txt"},
			},
			want: []string{
				"Analyze the following code.\n",
				"Focus on:\n- Performance\n- Security\n- Readability\n",
			},
		},
		{
			name: "empty prompt file",
			cfg: config.ConfigData{
				PromptPaths: []string{"testdata/prompt_empty.txt"},
			},
			want: []string{""},
		},
		{
			name: "flags and files combined",
			cfg: config.ConfigData{
				PromptFlags: []string{"first flag", "second flag"},
				PromptPaths: []string{"testdata/prompt1.txt"},
			},
			want: []string{
				"first flag",
				"second flag",
				"Analyze the following code.\n",
			},
		},
		{
			name: "file not found",
			cfg: config.ConfigData{
				PromptPaths: []string{"testdata/nonexistent.txt"},
			},
			wantErr: true,
		},
		{
			name: "mixed valid and invalid",
			cfg: config.ConfigData{
				PromptFlags: []string{"valid flag"},
				PromptPaths: []string{"testdata/nonexistent.txt"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadPromptSources(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReadFileSources(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.ConfigData
		want    []FileData
		wantErr bool
	}{
		{
			name: "no files returns empty",
			cfg:  config.ConfigData{},
			want: []FileData{},
		},
		{
			name: "single file",
			cfg: config.ConfigData{
				FilePaths: []string{"testdata/code.go"},
			},
			want: []FileData{
				{
					Path:    "testdata/code.go",
					Content: "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}\n",
				},
			},
		},
		{
			name: "multiple files",
			cfg: config.ConfigData{
				FilePaths: []string{"testdata/code.go", "testdata/data.json"},
			},
			want: []FileData{
				{
					Path:    "testdata/code.go",
					Content: "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}\n",
				},
				{
					Path:    "testdata/data.json",
					Content: "{\n  \"name\": \"test\",\n  \"value\": 42\n}\n",
				},
			},
		},
		{
			name: "empty file path",
			cfg: config.ConfigData{
				FilePaths: []string{""},
			},
			wantErr: true,
		},
		{
			name: "file not found",
			cfg: config.ConfigData{
				FilePaths: []string{"testdata/nonexistent.go"},
			},
			wantErr: true,
		},
		{
			name: "permission denied",
			cfg: config.ConfigData{
				FilePaths: []string{"/root/secret.txt"},
			},
			wantErr: true,
		},
		{
			name: "mixed valid and invalid",
			cfg: config.ConfigData{
				FilePaths: []string{"testdata/code.go", "testdata/nonexistent.go"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadFileSources(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

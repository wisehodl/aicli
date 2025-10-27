package input

import (
	"testing"

	"git.wisehodl.dev/jay/aicli/config"
	"github.com/stretchr/testify/assert"
)

func TestDetermineRole(t *testing.T) {
	tests := []struct {
		name     string
		cfg      config.ConfigData
		hasStdin bool
		want     StdinRole
	}{
		{
			name:     "no stdin returns StdinAsPrompt",
			cfg:      config.ConfigData{},
			hasStdin: false,
			want:     StdinAsPrompt,
		},
		{
			name:     "stdin with no flags returns StdinAsPrompt",
			cfg:      config.ConfigData{},
			hasStdin: true,
			want:     StdinAsPrompt,
		},
		{
			name: "stdin with -p flag returns StdinAsPrefixedContent",
			cfg: config.ConfigData{
				PromptFlags: []string{"analyze this"},
			},
			hasStdin: true,
			want:     StdinAsPrefixedContent,
		},
		{
			name: "stdin with -pf flag returns StdinAsPrefixedContent",
			cfg: config.ConfigData{
				PromptPaths: []string{"prompt.txt"},
			},
			hasStdin: true,
			want:     StdinAsPrefixedContent,
		},
		{
			name: "stdin with -F flag returns StdinAsFile",
			cfg: config.ConfigData{
				StdinAsFile: true,
			},
			hasStdin: true,
			want:     StdinAsFile,
		},
		{
			name: "stdin with -F and -p returns StdinAsFile (explicit wins)",
			cfg: config.ConfigData{
				StdinAsFile: true,
				PromptFlags: []string{"analyze"},
			},
			hasStdin: true,
			want:     StdinAsFile,
		},
		{
			name: "stdin with file flags returns StdinAsPrompt",
			cfg: config.ConfigData{
				FilePaths: []string{"main.go"},
			},
			hasStdin: true,
			want:     StdinAsPrompt,
		},
		{
			name: "no stdin with -F returns StdinAsFile (role set but unused)",
			cfg: config.ConfigData{
				StdinAsFile: true,
			},
			hasStdin: false,
			want:     StdinAsPrompt,
		},
		{
			name: "stdin with both -p and -pf returns StdinAsPrefixedContent",
			cfg: config.ConfigData{
				PromptFlags: []string{"first"},
				PromptPaths: []string{"prompt.txt"},
			},
			hasStdin: true,
			want:     StdinAsPrefixedContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetermineRole(tt.cfg, tt.hasStdin)
			assert.Equal(t, tt.want, got)
		})
	}
}

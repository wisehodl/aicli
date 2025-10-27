package input

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregatePrompts(t *testing.T) {
	tests := []struct {
		name    string
		prompts []string
		stdin   string
		role    StdinRole
		want    []string
	}{
		{
			name:    "empty inputs returns empty",
			prompts: []string{},
			stdin:   "",
			role:    StdinAsPrompt,
			want:    []string{},
		},
		{
			name:    "stdin as prompt with no other prompts",
			prompts: []string{},
			stdin:   "stdin content",
			role:    StdinAsPrompt,
			want:    []string{"stdin content"},
		},
		{
			name:    "stdin as prompt replaces existing prompts",
			prompts: []string{"prompt1", "prompt2"},
			stdin:   "stdin content",
			role:    StdinAsPrompt,
			want:    []string{"stdin content"},
		},
		{
			name:    "no stdin with role prompt returns prompts unchanged",
			prompts: []string{"prompt1", "prompt2"},
			stdin:   "",
			role:    StdinAsPrompt,
			want:    []string{"prompt1", "prompt2"},
		},
		{
			name:    "stdin as prefixed appends to prompts",
			prompts: []string{"prompt1", "prompt2"},
			stdin:   "stdin content",
			role:    StdinAsPrefixedContent,
			want:    []string{"prompt1", "prompt2", "stdin content"},
		},
		{
			name:    "stdin as prefixed with no prompts",
			prompts: []string{},
			stdin:   "stdin content",
			role:    StdinAsPrefixedContent,
			want:    []string{"stdin content"},
		},
		{
			name:    "no stdin with role prefixed returns prompts unchanged",
			prompts: []string{"prompt1"},
			stdin:   "",
			role:    StdinAsPrefixedContent,
			want:    []string{"prompt1"},
		},
		{
			name:    "stdin as file excludes stdin from prompts",
			prompts: []string{"prompt1"},
			stdin:   "stdin content",
			role:    StdinAsFile,
			want:    []string{"prompt1"},
		},
		{
			name:    "no stdin with role file returns prompts unchanged",
			prompts: []string{"prompt1"},
			stdin:   "",
			role:    StdinAsFile,
			want:    []string{"prompt1"},
		},
		{
			name:    "empty string stdin with role prompt",
			prompts: []string{"prompt1"},
			stdin:   "",
			role:    StdinAsPrompt,
			want:    []string{"prompt1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AggregatePrompts(tt.prompts, tt.stdin, tt.role)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAggregateFiles(t *testing.T) {
	tests := []struct {
		name  string
		files []FileData
		stdin string
		role  StdinRole
		want  []FileData
	}{
		{
			name:  "empty inputs returns empty",
			files: []FileData{},
			stdin: "",
			role:  StdinAsFile,
			want:  []FileData{},
		},
		{
			name:  "stdin as file prepends to files",
			files: []FileData{{Path: "a.go", Content: "code"}},
			stdin: "stdin content",
			role:  StdinAsFile,
			want: []FileData{
				{Path: "input", Content: "stdin content"},
				{Path: "a.go", Content: "code"},
			},
		},
		{
			name:  "stdin as file with no other files",
			files: []FileData{},
			stdin: "stdin content",
			role:  StdinAsFile,
			want: []FileData{
				{Path: "input", Content: "stdin content"},
			},
		},
		{
			name:  "no stdin with role file returns files unchanged",
			files: []FileData{{Path: "a.go", Content: "code"}},
			stdin: "",
			role:  StdinAsFile,
			want:  []FileData{{Path: "a.go", Content: "code"}},
		},
		{
			name:  "stdin as prompt excludes stdin from files",
			files: []FileData{{Path: "a.go", Content: "code"}},
			stdin: "stdin content",
			role:  StdinAsPrompt,
			want:  []FileData{{Path: "a.go", Content: "code"}},
		},
		{
			name:  "stdin as prefixed excludes stdin from files",
			files: []FileData{{Path: "a.go", Content: "code"}},
			stdin: "stdin content",
			role:  StdinAsPrefixedContent,
			want:  []FileData{{Path: "a.go", Content: "code"}},
		},
		{
			name: "stdin as file with multiple files",
			files: []FileData{
				{Path: "a.go", Content: "code a"},
				{Path: "b.go", Content: "code b"},
			},
			stdin: "stdin content",
			role:  StdinAsFile,
			want: []FileData{
				{Path: "input", Content: "stdin content"},
				{Path: "a.go", Content: "code a"},
				{Path: "b.go", Content: "code b"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AggregateFiles(tt.files, tt.stdin, tt.role)
			assert.Equal(t, tt.want, got)
		})
	}
}

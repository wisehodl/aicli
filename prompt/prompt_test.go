package prompt

import (
	"testing"

	"git.wisehodl.dev/jay/aicli/input"
	"github.com/stretchr/testify/assert"
)

func TestFormatPrompts(t *testing.T) {
	tests := []struct {
		name    string
		prompts []string
		want    string
	}{
		{
			name:    "empty array returns empty string",
			prompts: []string{},
			want:    "",
		},
		{
			name:    "single prompt unchanged",
			prompts: []string{"analyze this"},
			want:    "analyze this",
		},
		{
			name:    "multiple prompts joined with newline",
			prompts: []string{"first", "second", "third"},
			want:    "first\nsecond\nthird",
		},
		{
			name:    "prompts with trailing newlines preserved",
			prompts: []string{"line one\n", "line two\n"},
			want:    "line one\n\nline two\n",
		},
		{
			name:    "empty string in array produces empty line",
			prompts: []string{"first", "", "third"},
			want:    "first\n\nthird",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatPrompts(tt.prompts)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormatFiles(t *testing.T) {
	tests := []struct {
		name  string
		files []input.FileData
		want  string
	}{
		{
			name:  "empty array returns empty string",
			files: []input.FileData{},
			want:  "",
		},
		{
			name: "single file formatted with template",
			files: []input.FileData{
				{Path: "main.go", Content: "package main"},
			},
			want: "File: main.go\n\n```\npackage main\n```",
		},
		{
			name: "multiple files separated by double newline",
			files: []input.FileData{
				{Path: "a.go", Content: "code a"},
				{Path: "b.go", Content: "code b"},
			},
			want: "File: a.go\n\n```\ncode a\n```\n\nFile: b.go\n\n```\ncode b\n```",
		},
		{
			name: "stdin path 'input' appears correctly",
			files: []input.FileData{
				{Path: "input", Content: "stdin content"},
			},
			want: "File: input\n\n```\nstdin content\n```",
		},
		{
			name: "file path with directory",
			files: []input.FileData{
				{Path: "src/main.go", Content: "package main"},
			},
			want: "File: src/main.go\n\n```\npackage main\n```",
		},
		{
			name: "content with backticks still wrapped",
			files: []input.FileData{
				{Path: "test.md", Content: "```go\nfunc main() {}\n```"},
			},
			want: "File: test.md\n\n```\n```go\nfunc main() {}\n```\n```",
		},
		{
			name: "empty content",
			files: []input.FileData{
				{Path: "empty.txt", Content: ""},
			},
			want: "File: empty.txt\n\n```\n\n```",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatFiles(tt.files)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCombineContent(t *testing.T) {
	tests := []struct {
		name      string
		promptStr string
		filesStr  string
		want      string
	}{
		{
			name:      "both empty returns empty",
			promptStr: "",
			filesStr:  "",
			want:      "",
		},
		{
			name:      "prompt only",
			promptStr: "analyze this",
			filesStr:  "",
			want:      "analyze this",
		},
		{
			name:      "files only uses default prompt",
			promptStr: "",
			filesStr:  "File: a.go\n\n```\ncode\n```",
			want:      "Analyze the following:\n\nFile: a.go\n\n```\ncode\n```",
		},
		{
			name:      "prompt and files combined with separator",
			promptStr: "review this code",
			filesStr:  "File: a.go\n\n```\ncode\n```",
			want:      "review this code\n\nFile: a.go\n\n```\ncode\n```",
		},
		{
			name:      "multiline prompt preserved",
			promptStr: "first line\nsecond line",
			filesStr:  "File: a.go\n\n```\ncode\n```",
			want:      "first line\nsecond line\n\nFile: a.go\n\n```\ncode\n```",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := combineContent(tt.promptStr, tt.filesStr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConstructQuery(t *testing.T) {
	tests := []struct {
		name    string
		prompts []string
		files   []input.FileData
		want    string
	}{
		{
			name:    "empty inputs returns empty",
			prompts: []string{},
			files:   []input.FileData{},
			want:    "",
		},
		{
			name:    "prompt only",
			prompts: []string{"analyze this"},
			files:   []input.FileData{},
			want:    "analyze this",
		},
		{
			name:    "file only with default prompt",
			prompts: []string{},
			files: []input.FileData{
				{Path: "main.go", Content: "package main"},
			},
			want: "Analyze the following:\n\nFile: main.go\n\n```\npackage main\n```",
		},
		{
			name:    "multiple prompts and files",
			prompts: []string{"review", "focus on bugs"},
			files: []input.FileData{
				{Path: "a.go", Content: "code a"},
				{Path: "b.go", Content: "code b"},
			},
			want: "review\nfocus on bugs\n\nFile: a.go\n\n```\ncode a\n```\n\nFile: b.go\n\n```\ncode b\n```",
		},
		{
			name:    "stdin as file with explicit prompt",
			prompts: []string{"analyze"},
			files: []input.FileData{
				{Path: "input", Content: "stdin data"},
			},
			want: "analyze\n\nFile: input\n\n```\nstdin data\n```",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConstructQuery(tt.prompts, tt.files)
			assert.Equal(t, tt.want, got)
		})
	}
}

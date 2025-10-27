package prompt

import (
	"fmt"
	"strings"

	"git.wisehodl.dev/jay/aicli/input"
)

const defaultPrompt = "Analyze the following:"

// ConstructQuery formats prompts and files into a complete query string.
func ConstructQuery(prompts []string, files []input.FileData) string {
	promptStr := formatPrompts(prompts)
	filesStr := formatFiles(files)
	return combineContent(promptStr, filesStr)
}

// formatPrompts joins prompt strings with newlines.
func formatPrompts(prompts []string) string {
	if len(prompts) == 0 {
		return ""
	}
	return strings.Join(prompts, "\n")
}

// formatFiles wraps each file in a template with path and content.
func formatFiles(files []input.FileData) string {
	if len(files) == 0 {
		return ""
	}

	var parts []string
	for _, f := range files {
		parts = append(parts, fmt.Sprintf("File: %s\n\n```\n%s\n```", f.Path, f.Content))
	}
	return strings.Join(parts, "\n\n")
}

// combineContent merges formatted prompts and files with appropriate separators.
func combineContent(promptStr, filesStr string) string {
	if promptStr == "" && filesStr == "" {
		return ""
	}

	if promptStr == "" && filesStr != "" {
		return defaultPrompt + "\n\n" + filesStr
	}

	if promptStr != "" && filesStr == "" {
		return promptStr
	}

	return promptStr + "\n\n" + filesStr
}

package input

import (
	"fmt"
	"os"

	"git.wisehodl.dev/jay/aicli/config"
)

// ReadPromptSources reads all prompt content from flags and files.
// Returns arrays of prompt strings in source order.
func ReadPromptSources(cfg config.ConfigData) ([]string, error) {
	prompts := []string{}

	// Add flag prompts first
	prompts = append(prompts, cfg.PromptFlags...)

	// Add prompt file contents
	for _, path := range cfg.PromptPaths {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read prompt file %s: %w", path, err)
		}
		prompts = append(prompts, string(content))
	}

	return prompts, nil
}

// ReadFileSources reads all input files specified in config.
// Returns FileData array in source order.
func ReadFileSources(cfg config.ConfigData) ([]FileData, error) {
	files := []FileData{}

	for _, path := range cfg.FilePaths {
		if path == "" {
			return nil, fmt.Errorf("empty file path provided")
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read file %s: %w", path, err)
		}

		files = append(files, FileData{
			Path:    path,
			Content: string(content),
		})
	}

	return files, nil
}

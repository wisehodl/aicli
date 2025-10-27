package input

import (
	"fmt"

	"git.wisehodl.dev/jay/aicli/config"
)

// ResolveInputs orchestrates the complete input resolution pipeline.
// Returns aggregated prompts and files ready for query construction.
func ResolveInputs(cfg config.ConfigData, stdinContent string, hasStdin bool) (InputData, error) {
	// Determine stdin role (CA -> CB)
	role := DetermineRole(cfg, hasStdin)

	// Read all sources (CC, CD)
	prompts, err := ReadPromptSources(cfg)
	if err != nil {
		return InputData{}, err
	}

	files, err := ReadFileSources(cfg)
	if err != nil {
		return InputData{}, err
	}

	// Aggregate with stdin (CE, CF)
	finalPrompts := AggregatePrompts(prompts, stdinContent, role)
	finalFiles := AggregateFiles(files, stdinContent, role)

	// Validate at least one input exists
	if len(finalPrompts) == 0 && len(finalFiles) == 0 {
		return InputData{}, fmt.Errorf("no input provided: supply stdin, --file, or --prompt")
	}

	return InputData{
		Prompts: finalPrompts,
		Files:   finalFiles,
	}, nil
}

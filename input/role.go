package input

import "git.wisehodl.dev/jay/aicli/config"

// DetermineRole decides how stdin content participates in the query based on
// flags and stdin presence. Per spec §7 rules.
func DetermineRole(cfg config.ConfigData, hasStdin bool) StdinRole {
	if !hasStdin {
		return StdinAsPrompt // unused, but set for consistency
	}

	// Explicit -F flag forces stdin as file
	if cfg.StdinAsFile {
		return StdinAsFile
	}

	// Any explicit prompt flag (-p or -pf) makes stdin prefixed content
	hasExplicitPrompt := len(cfg.PromptFlags) > 0 || len(cfg.PromptPaths) > 0

	if hasExplicitPrompt {
		return StdinAsPrefixedContent
	}

	// Default: stdin replaces any default prompt
	return StdinAsPrompt
}

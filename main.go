package main

import (
	"fmt"
	"os"

	"git.wisehodl.dev/jay/aicli/api"
	"git.wisehodl.dev/jay/aicli/config"
	"git.wisehodl.dev/jay/aicli/input"
	"git.wisehodl.dev/jay/aicli/output"
	"git.wisehodl.dev/jay/aicli/prompt"
	"git.wisehodl.dev/jay/aicli/version"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Phase 1: Version check (early exit)
	if config.IsVersionRequest(os.Args[1:]) {
		fmt.Printf("aicli %s\n", version.GetVersion())
		return nil
	}

	if config.IsHelpRequest(os.Args[1:]) {
		fmt.Fprint(os.Stderr, config.UsageText)
		return nil
	}

	// Phase 2: Configuration resolution
	cfg, err := config.BuildConfig(os.Args[1:])
	if err != nil {
		return err
	}

	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "[verbose] Configuration loaded\n")
		fmt.Fprintf(os.Stderr, "  Protocol: %s\n", protocolString(cfg.Protocol))
		fmt.Fprintf(os.Stderr, "  URL: %s\n", cfg.URL)
		fmt.Fprintf(os.Stderr, "  Model: %s\n", cfg.Model)
		fmt.Fprintf(os.Stderr, "  Fallbacks: %v\n", cfg.FallbackModels)
	}

	// Phase 3: Input collection
	stdinContent, hasStdin := input.DetectStdin()

	inputData, err := input.ResolveInputs(cfg, stdinContent, hasStdin)
	if err != nil {
		return err
	}

	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "[verbose] Input resolved: %d prompts, %d files\n",
			len(inputData.Prompts), len(inputData.Files))
	}

	// Phase 4: Query construction
	query := prompt.ConstructQuery(inputData.Prompts, inputData.Files)

	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "[verbose] Query length: %d bytes\n", len(query))
	}

	// Phase 5: API communication
	response, model, duration, err := api.SendChatRequest(cfg, query)
	if err != nil {
		return err
	}

	// Phase 6: Output delivery
	return output.WriteOutput(response, model, duration, cfg)
}

func protocolString(p config.APIProtocol) string {
	if p == config.ProtocolOllama {
		return "ollama"
	}
	return "openai"
}

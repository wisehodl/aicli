package config

import (
	"os"
	"strings"
)

func mergeSources(flags flagValues, env envValues, file fileValues) ConfigData {
	cfg := defaultConfig

	// Apply file values
	if file.protocol != "" {
		cfg.Protocol = parseProtocol(file.protocol)
	}
	if file.url != "" {
		cfg.URL = file.url
	}
	if file.model != "" {
		cfg.Model = file.model
	}
	if file.fallback != "" {
		cfg.FallbackModels = strings.Split(file.fallback, ",")
	}

	// Apply env values
	if env.protocol != "" {
		cfg.Protocol = parseProtocol(env.protocol)
	}
	if env.url != "" {
		cfg.URL = env.url
	}
	if env.model != "" {
		cfg.Model = env.model
	}
	if env.fallback != "" {
		cfg.FallbackModels = strings.Split(env.fallback, ",")
	}
	if env.system != "" {
		cfg.SystemPrompt = env.system
	}
	if env.key != "" {
		cfg.APIKey = env.key
	}

	// Apply flag values
	if flags.protocol != "" {
		cfg.Protocol = parseProtocol(flags.protocol)
	}
	if flags.url != "" {
		cfg.URL = flags.url
	}
	if flags.model != "" {
		cfg.Model = flags.model
	}
	if flags.fallback != "" {
		cfg.FallbackModels = strings.Split(flags.fallback, ",")
	}
	if flags.output != "" {
		cfg.Output = flags.output
	}
	cfg.Quiet = flags.quiet
	cfg.Verbose = flags.verbose
	cfg.StdinAsFile = flags.stdinFile

	// Collect input paths
	cfg.FilePaths = flags.files
	cfg.PromptFlags = flags.prompts
	if flags.promptFile != "" {
		cfg.PromptPaths = []string{flags.promptFile}
	}

	// Resolve system prompt (direct > file)
	if flags.system != "" {
		cfg.SystemPrompt = flags.system
	} else if flags.systemFile != "" {
		content, err := os.ReadFile(flags.systemFile)
		if err == nil {
			cfg.SystemPrompt = strings.TrimRight(string(content), "\n")
		}
	} else if file.systemFile != "" && cfg.SystemPrompt == "" {
		content, err := os.ReadFile(file.systemFile)
		if err == nil {
			cfg.SystemPrompt = strings.TrimRight(string(content), "\n")
		}
	}

	// Resolve API key (direct > file)
	if flags.key != "" {
		cfg.APIKey = flags.key
	} else if flags.keyFile != "" {
		content, err := os.ReadFile(flags.keyFile)
		if err == nil {
			cfg.APIKey = strings.TrimSpace(string(content))
		}
	} else if cfg.APIKey == "" && file.keyFile != "" {
		content, err := os.ReadFile(file.keyFile)
		if err == nil {
			cfg.APIKey = strings.TrimSpace(string(content))
		}
	}

	return cfg
}

func parseProtocol(s string) APIProtocol {
	switch s {
	case "ollama":
		return ProtocolOllama
	default:
		return ProtocolOpenAI
	}
}

package config

import (
	"fmt"
)

func validateConfig(cfg ConfigData) error {
	if cfg.APIKey == "" {
		return fmt.Errorf("API key required: use --key, --key-file, AICLI_API_KEY, AICLI_API_KEY_FILE, or key_file in config")
	}

	if cfg.Protocol != ProtocolOpenAI && cfg.Protocol != ProtocolOllama {
		return fmt.Errorf("invalid protocol: must be openai or ollama")
	}

	return nil
}

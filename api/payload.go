package api

import "git.wisehodl.dev/jay/aicli/config"

// buildPayload constructs the JSON payload for the API request based on protocol.
func buildPayload(cfg config.ConfigData, model string, query string) map[string]interface{} {
	if cfg.Protocol == config.ProtocolOllama {
		payload := map[string]interface{}{
			"model":  model,
			"prompt": query,
			"stream": false,
		}
		if cfg.SystemPrompt != "" {
			payload["system"] = cfg.SystemPrompt
		}
		return payload
	}

	// OpenAI protocol
	messages := []map[string]string{}
	if cfg.SystemPrompt != "" {
		messages = append(messages, map[string]string{
			"role":    "system",
			"content": cfg.SystemPrompt,
		})
	}
	messages = append(messages, map[string]string{
		"role":    "user",
		"content": query,
	})

	return map[string]interface{}{
		"model":    model,
		"messages": messages,
	}
}

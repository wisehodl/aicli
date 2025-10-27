package api

import (
	"encoding/json"
	"fmt"

	"git.wisehodl.dev/jay/aicli/config"
)

// parseResponse extracts the response content from the API response body.
func parseResponse(body []byte, protocol config.APIProtocol) (string, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	if protocol == config.ProtocolOllama {
		response, ok := result["response"].(string)
		if !ok {
			return "", fmt.Errorf("no response field in ollama response")
		}
		return response, nil
	}

	// OpenAI protocol
	choices, ok := result["choices"].([]interface{})
	if !ok {
		return "", fmt.Errorf("no choices in response")
	}

	if len(choices) == 0 {
		return "", fmt.Errorf("empty choices array")
	}

	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid choice format")
	}

	message, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("no message in choice")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("no content in message")
	}

	return content, nil
}

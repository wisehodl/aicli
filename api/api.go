package api

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"git.wisehodl.dev/jay/aicli/config"
)

// tryModel attempts a single model request through the complete pipeline:
// payload construction, HTTP execution, and response parsing.
func tryModel(cfg config.ConfigData, model string, query string) (string, error) {
	payload := buildPayload(cfg, model, query)

	if cfg.Verbose {
		payloadJSON, _ := json.Marshal(payload)
		fmt.Fprintf(os.Stderr, "[verbose] Request payload: %s\n", string(payloadJSON))
	}

	body, err := executeHTTP(cfg, payload)
	if err != nil {
		return "", err
	}

	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "[verbose] Response: %s\n", string(body))
	}

	response, err := parseResponse(body, cfg.Protocol)
	if err != nil {
		return "", err
	}

	return response, nil
}

// SendChatRequest sends a query to the configured model with automatic fallback.
// Returns the response content, the model name that succeeded, total duration, and any error.
// On failure, attempts each fallback model in sequence until one succeeds or all fail.
func SendChatRequest(cfg config.ConfigData, query string) (string, string, time.Duration, error) {
	models := append([]string{cfg.Model}, cfg.FallbackModels...)
	start := time.Now()

	for i, model := range models {
		if !cfg.Quiet && i > 0 {
			fmt.Fprintf(os.Stderr, "Model %s failed, trying %s...\n", models[i-1], model)
		}

		response, err := tryModel(cfg, model, query)
		if err == nil {
			return response, model, time.Since(start), nil
		}

		if !cfg.Quiet {
			fmt.Fprintf(os.Stderr, "Model %s failed: %v\n", model, err)
		}
	}

	return "", "", time.Since(start), fmt.Errorf("all models failed")
}

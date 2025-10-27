package output

import (
	"fmt"
	"os"
	"time"

	"git.wisehodl.dev/jay/aicli/config"
)

// WriteOutput orchestrates complete output delivery based on configuration.
func WriteOutput(response, model string, duration time.Duration, cfg config.ConfigData) error {
	if cfg.Output == "" {
		// Write to stdout with optional metadata
		formatted := formatOutput(response, model, duration, cfg.Quiet)
		return writeStdout(formatted)
	}

	// Write raw response to file
	if err := writeFile(response, cfg.Output); err != nil {
		return err
	}

	// Write metadata to stderr unless quiet
	if !cfg.Quiet {
		metadata := fmt.Sprintf("Used model: %s\nQuery duration: %.1fs\nWrote response to: %s\n",
			model, duration.Seconds(), cfg.Output)
		return writeStderr(metadata)
	}

	return nil
}

// formatOutput constructs the final output string with optional metadata header.
func formatOutput(response, model string, duration time.Duration, quiet bool) string {
	if quiet {
		return response
	}

	return fmt.Sprintf(`--- aicli ---

Used model: %s
Query duration: %.1fs

--- response ---

%s`, model, duration.Seconds(), response)
}

// writeStdout writes content to stdout.
func writeStdout(content string) error {
	_, err := fmt.Println(content)
	if err != nil {
		return fmt.Errorf("write stdout: %w", err)
	}
	return nil
}

// writeStderr writes logs to stderr.
func writeStderr(content string) error {
	_, err := fmt.Fprint(os.Stderr, content)
	if err != nil {
		return fmt.Errorf("write stderr: %w", err)
	}
	return nil
}

// writeFile writes content to the specified path with permissions 0644.
func writeFile(content, path string) error {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}
	return nil
}

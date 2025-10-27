package input

import (
	"io"
	"os"
)

// DetectStdin checks if stdin contains piped data and reads it.
// Returns content and true if stdin is a pipe/file, empty string and false if terminal.
func DetectStdin() (string, bool) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", false
	}

	// Terminal (character device) = no stdin data
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", false
	}

	// Pipe or file redirection detected
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", false
	}

	return string(content), true
}

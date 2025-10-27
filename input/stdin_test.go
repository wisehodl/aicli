package input

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Integration test helper - run manually with:
// echo "test" | STDIN_TEST=1 go test ./input/ -run TestDetectStdinIntegration
func TestDetectStdinIntegration(t *testing.T) {
	if os.Getenv("STDIN_TEST") != "1" {
		t.Skip("Set STDIN_TEST=1 and pipe data to run this test")
	}

	content, hasStdin := DetectStdin()

	t.Logf("hasStdin: %v", hasStdin)
	t.Logf("content length: %d", len(content))
	t.Logf("content: %q", content)

	// When piped: should detect stdin
	if hasStdin {
		assert.NotEmpty(t, content, "Expected content when stdin detected")
	}
}

package output

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"git.wisehodl.dev/jay/aicli/config"
	"github.com/stretchr/testify/assert"
)

func TestFormatOutput(t *testing.T) {
	tests := []struct {
		name     string
		response string
		model    string
		duration time.Duration
		quiet    bool
		want     string
	}{
		{
			name:     "normal mode with metadata",
			response: "This is the response.",
			model:    "gpt-4",
			duration: 2500 * time.Millisecond,
			quiet:    false,
			want: `--- aicli ---

Used model: gpt-4
Query duration: 2.5s

--- response ---

This is the response.`,
		},
		{
			name:     "quiet mode response only",
			response: "This is the response.",
			model:    "gpt-4",
			duration: 2500 * time.Millisecond,
			quiet:    true,
			want:     "This is the response.",
		},
		{
			name:     "duration formatting subsecond",
			response: "response",
			model:    "gpt-3.5",
			duration: 123 * time.Millisecond,
			quiet:    false,
			want: `--- aicli ---

Used model: gpt-3.5
Query duration: 0.1s

--- response ---

response`,
		},
		{
			name:     "duration formatting multi-second",
			response: "response",
			model:    "claude-3",
			duration: 12345 * time.Millisecond,
			quiet:    false,
			want: `--- aicli ---

Used model: claude-3
Query duration: 12.3s

--- response ---

response`,
		},
		{
			name:     "multiline response preserved",
			response: "Line 1\nLine 2\nLine 3",
			model:    "gpt-4",
			duration: 1 * time.Second,
			quiet:    false,
			want: `--- aicli ---

Used model: gpt-4
Query duration: 1.0s

--- response ---

Line 1
Line 2
Line 3`,
		},
		{
			name:     "empty response",
			response: "",
			model:    "gpt-4",
			duration: 1 * time.Second,
			quiet:    false,
			want: `--- aicli ---

Used model: gpt-4
Query duration: 1.0s

--- response ---

`,
		},
		{
			name:     "model name with special chars",
			response: "response",
			model:    "gpt-4-1106-preview",
			duration: 5 * time.Second,
			quiet:    false,
			want: `--- aicli ---

Used model: gpt-4-1106-preview
Query duration: 5.0s

--- response ---

response`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatOutput(tt.response, tt.model, tt.duration, tt.quiet)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWriteStdout(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "normal content",
			content: "test output",
		},
		{
			name:    "empty string",
			content: "",
		},
		{
			name:    "multiline content",
			content: "line 1\nline 2\nline 3",
		},
		{
			name:    "large content",
			content: string(make([]byte, 10000)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := writeStdout(tt.content)
			assert.NoError(t, err)

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)

			// writeStdout uses fmt.Println which adds newline
			expected := tt.content + "\n"
			assert.Equal(t, expected, buf.String())
		})
	}
}

func TestWriteStderr(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "normal content",
			content: "error message",
		},
		{
			name:    "empty string",
			content: "",
		},
		{
			name:    "multiline content",
			content: "line 1\nline 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			err := writeStderr(tt.content)
			assert.NoError(t, err)

			w.Close()
			os.Stderr = old

			var buf bytes.Buffer
			io.Copy(&buf, r)

			assert.Equal(t, tt.content, buf.String())
		})
	}
}

func TestWriteFile(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantErr     bool
		errContains string
	}{
		{
			name:    "normal write",
			content: "test content",
		},
		{
			name:    "empty content",
			content: "",
		},
		{
			name:    "multiline content",
			content: "line 1\nline 2\nline 3",
		},
		{
			name:    "large content",
			content: string(make([]byte, 100000)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			path := filepath.Join(tmpDir, "output.txt")

			err := writeFile(tt.content, path)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			assert.NoError(t, err)

			// Verify file exists and has correct content
			got, err := os.ReadFile(path)
			assert.NoError(t, err)
			assert.Equal(t, tt.content, string(got))

			// Verify permissions
			info, err := os.Stat(path)
			assert.NoError(t, err)
			assert.Equal(t, os.FileMode(0644), info.Mode().Perm())
		})
	}
}

func TestWriteFileErrors(t *testing.T) {
	tests := []struct {
		name        string
		setupPath   func() string
		errContains string
	}{
		{
			name: "directory does not exist",
			setupPath: func() string {
				return "/nonexistent/dir/output.txt"
			},
			errContains: "write output file",
		},
		{
			name: "permission denied",
			setupPath: func() string {
				tmpDir := t.TempDir()
				dir := filepath.Join(tmpDir, "readonly")
				os.Mkdir(dir, 0444)
				return filepath.Join(dir, "output.txt")
			},
			errContains: "write output file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setupPath()
			err := writeFile("content", path)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errContains)
		})
	}
}

func TestWriteFileOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "output.txt")

	// Write initial content
	err := writeFile("initial", path)
	assert.NoError(t, err)

	got, _ := os.ReadFile(path)
	assert.Equal(t, "initial", string(got))

	// Overwrite with new content
	err = writeFile("overwritten", path)
	assert.NoError(t, err)

	got, _ = os.ReadFile(path)
	assert.Equal(t, "overwritten", string(got))
}

func TestWriteOutput(t *testing.T) {
	tests := []struct {
		name        string
		response    string
		model       string
		duration    time.Duration
		cfg         config.ConfigData
		checkStdout bool
		checkStderr bool
		checkFile   bool
		wantStdout  string
		wantStderr  string
		wantErr     bool
		errContains string
	}{
		{
			name:     "stdout with metadata",
			response: "response text",
			model:    "gpt-4",
			duration: 2 * time.Second,
			cfg: config.ConfigData{
				Quiet: false,
			},
			checkStdout: true,
			wantStdout: `--- aicli ---

Used model: gpt-4
Query duration: 2.0s

--- response ---

response text
`,
		},
		{
			name:     "stdout quiet mode",
			response: "response text",
			model:    "gpt-4",
			duration: 2 * time.Second,
			cfg: config.ConfigData{
				Quiet: true,
			},
			checkStdout: true,
			wantStdout:  "response text\n",
		},
		{
			name:     "file output with stderr metadata",
			response: "response text",
			model:    "gpt-4",
			duration: 3 * time.Second,
			cfg: config.ConfigData{
				Output: "output.txt",
				Quiet:  false,
			},
			checkFile:   true,
			checkStderr: true,
			wantStderr:  "Used model: gpt-4\nQuery duration: 3.0s\nWrote response to: .*output.txt\n",
		},
		{
			name:     "file output quiet mode",
			response: "response text",
			model:    "gpt-4",
			duration: 3 * time.Second,
			cfg: config.ConfigData{
				Output: "output.txt",
				Quiet:  true,
			},
			checkFile: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Capture stdout if needed
			oldStdout := os.Stdout
			var stdoutR *os.File
			if tt.checkStdout {
				r, w, _ := os.Pipe()
				os.Stdout = w
				stdoutR = r
			}

			// Capture stderr if needed
			oldStderr := os.Stderr
			var stderrR *os.File
			if tt.checkStderr {
				r, w, _ := os.Pipe()
				os.Stderr = w
				stderrR = r
			}

			// Set output path if needed
			if tt.cfg.Output != "" {
				tt.cfg.Output = filepath.Join(tmpDir, tt.cfg.Output)
			}

			err := WriteOutput(tt.response, tt.model, tt.duration, tt.cfg)

			// Close write ends and restore originals
			if tt.checkStdout {
				os.Stdout.Close()
				os.Stdout = oldStdout
			}
			if tt.checkStderr {
				os.Stderr.Close()
				os.Stderr = oldStderr
			}

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			assert.NoError(t, err)

			// Read stdout
			if tt.checkStdout {
				var stdoutBuf bytes.Buffer
				io.Copy(&stdoutBuf, stdoutR)
				stdoutR.Close()
				assert.Equal(t, tt.wantStdout, stdoutBuf.String())
			}

			// Read stderr
			if tt.checkStderr {
				var stderrBuf bytes.Buffer
				io.Copy(&stderrBuf, stderrR)
				stderrR.Close()

				got := stderrBuf.String()
				assert.Contains(t, got, "Used model: gpt-4")
				assert.Contains(t, got, "Query duration: 3.0s")
				assert.Contains(t, got, "output.txt")
			}

			// Check file
			if tt.checkFile {
				content, err := os.ReadFile(tt.cfg.Output)
				assert.NoError(t, err)
				assert.Equal(t, tt.response, string(content))
			}
		})
	}
}

func TestWriteOutputFileError(t *testing.T) {
	cfg := config.ConfigData{
		Output: "/nonexistent/dir/output.txt",
		Quiet:  false,
	}

	err := WriteOutput("response", "gpt-4", 1*time.Second, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "write output file")
}

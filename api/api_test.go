package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"git.wisehodl.dev/jay/aicli/config"
	"github.com/stretchr/testify/assert"
)

type sequenceTransport struct {
	responses []*http.Response
	index     int
}

func (t *sequenceTransport) RoundTrip(*http.Request) (*http.Response, error) {
	if t.index >= len(t.responses) {
		return nil, fmt.Errorf("no more responses in sequence")
	}
	resp := t.responses[t.index]
	t.index++
	return resp, nil
}

func TestTryModel(t *testing.T) {
	tests := []struct {
		name        string
		cfg         config.ConfigData
		model       string
		query       string
		mockResp    *http.Response
		mockErr     error
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name: "successful request",
			cfg: config.ConfigData{
				Protocol: config.ProtocolOpenAI,
				URL:      "https://api.example.com",
				APIKey:   "sk-test",
			},
			model:    "gpt-4",
			query:    "test query",
			mockResp: makeResponse(200, `{"choices":[{"message":{"content":"response text"}}]}`),
			want:     "response text",
		},
		{
			name: "http error",
			cfg: config.ConfigData{
				Protocol: config.ProtocolOpenAI,
				URL:      "https://api.example.com",
				APIKey:   "sk-test",
			},
			model:       "gpt-4",
			query:       "test query",
			mockResp:    makeResponse(500, `{"error":"server error"}`),
			wantErr:     true,
			errContains: "HTTP 500",
		},
		{
			name: "parse error",
			cfg: config.ConfigData{
				Protocol: config.ProtocolOpenAI,
				URL:      "https://api.example.com",
				APIKey:   "sk-test",
			},
			model:       "gpt-4",
			query:       "test query",
			mockResp:    makeResponse(200, `{"choices":[]}`),
			wantErr:     true,
			errContains: "empty choices array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport := &mockRoundTripper{
				response: tt.mockResp,
				err:      tt.mockErr,
			}
			oldClient := httpClient
			httpClient = &http.Client{
				Timeout:   5 * time.Minute,
				Transport: transport,
			}
			defer func() { httpClient = oldClient }()

			got, err := tryModel(tt.cfg, tt.model, tt.query)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func captureStderr(f func()) string {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f()

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestSendChatRequest(t *testing.T) {
	tests := []struct {
		name         string
		cfg          config.ConfigData
		query        string
		mockResp     []*http.Response
		wantResponse string
		wantModel    string
		wantErr      bool
		errContains  string
		checkStderr  func(*testing.T, string)
	}{
		{
			name: "primary model succeeds",
			cfg: config.ConfigData{
				Protocol:       config.ProtocolOpenAI,
				URL:            "https://api.example.com",
				APIKey:         "sk-test",
				Model:          "gpt-4",
				FallbackModels: []string{"gpt-3.5"},
			},
			query: "test",
			mockResp: []*http.Response{
				makeResponse(200, `{"choices":[{"message":{"content":"primary response"}}]}`),
			},
			wantResponse: "primary response",
			wantModel:    "gpt-4",
		},
		{
			name: "primary fails, fallback succeeds",
			cfg: config.ConfigData{
				Protocol:       config.ProtocolOpenAI,
				URL:            "https://api.example.com",
				APIKey:         "sk-test",
				Model:          "gpt-4",
				FallbackModels: []string{"gpt-3.5"},
			},
			query: "test",
			mockResp: []*http.Response{
				makeResponse(500, `{"error":"server error"}`),
				makeResponse(200, `{"choices":[{"message":{"content":"fallback response"}}]}`),
			},
			wantResponse: "fallback response",
			wantModel:    "gpt-3.5",
			checkStderr: func(t *testing.T, stderr string) {
				assert.Contains(t, stderr, "Model gpt-4 failed")
				assert.Contains(t, stderr, "trying gpt-3.5")
			},
		},
		{
			name: "all models fail",
			cfg: config.ConfigData{
				Protocol:       config.ProtocolOpenAI,
				URL:            "https://api.example.com",
				APIKey:         "sk-test",
				Model:          "gpt-4",
				FallbackModels: []string{"gpt-3.5"},
			},
			query: "test",
			mockResp: []*http.Response{
				makeResponse(500, `{"error":"error1"}`),
				makeResponse(500, `{"error":"error2"}`),
			},
			wantErr:     true,
			errContains: "all models failed",
			checkStderr: func(t *testing.T, stderr string) {
				assert.Contains(t, stderr, "Model gpt-4 failed")
				assert.Contains(t, stderr, "Model gpt-3.5 failed")
			},
		},
		{
			name: "quiet mode suppresses progress",
			cfg: config.ConfigData{
				Protocol:       config.ProtocolOpenAI,
				URL:            "https://api.example.com",
				APIKey:         "sk-test",
				Model:          "gpt-4",
				FallbackModels: []string{"gpt-3.5"},
				Quiet:          true,
			},
			query: "test",
			mockResp: []*http.Response{
				makeResponse(500, `{"error":"error1"}`),
				makeResponse(200, `{"choices":[{"message":{"content":"response"}}]}`),
			},
			wantResponse: "response",
			wantModel:    "gpt-3.5",
			checkStderr: func(t *testing.T, stderr string) {
				assert.Empty(t, stderr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport := &sequenceTransport{responses: tt.mockResp}

			oldClient := httpClient
			httpClient = &http.Client{
				Timeout:   5 * time.Minute,
				Transport: transport,
			}
			defer func() { httpClient = oldClient }()

			var stderr string
			var response string
			var model string
			var duration time.Duration
			var err error

			if tt.checkStderr != nil {
				stderr = captureStderr(func() {
					response, model, duration, err = SendChatRequest(tt.cfg, tt.query)
				})
			} else {
				response, model, duration, err = SendChatRequest(tt.cfg, tt.query)
			}

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				if tt.checkStderr != nil {
					tt.checkStderr(t, stderr)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantResponse, response)
			assert.Equal(t, tt.wantModel, model)
			assert.Greater(t, duration, time.Duration(0))

			if tt.checkStderr != nil {
				tt.checkStderr(t, stderr)
			}
		})
	}
}

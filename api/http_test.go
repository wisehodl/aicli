package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"git.wisehodl.dev/jay/aicli/config"
	"github.com/stretchr/testify/assert"
)

type mockRoundTripper struct {
	response *http.Response
	err      error
	request  *http.Request
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.request = req
	return m.response, m.err
}

func makeResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func TestExecuteHTTP(t *testing.T) {
	tests := []struct {
		name        string
		cfg         config.ConfigData
		payload     map[string]interface{}
		mockResp    *http.Response
		mockErr     error
		wantBody    string
		wantErr     bool
		errContains string
	}{
		{
			name: "successful request",
			cfg: config.ConfigData{
				URL:    "https://api.example.com/chat",
				APIKey: "sk-test123",
			},
			payload: map[string]interface{}{
				"model": "gpt-4",
			},
			mockResp: makeResponse(200, `{"choices":[{"message":{"content":"response"}}]}`),
			wantBody: `{"choices":[{"message":{"content":"response"}}]}`,
			wantErr:  false,
		},
		{
			name: "HTTP 400 error",
			cfg: config.ConfigData{
				URL:    "https://api.example.com/chat",
				APIKey: "sk-test123",
			},
			payload:     map[string]interface{}{"model": "gpt-4"},
			mockResp:    makeResponse(400, `{"error":"bad request"}`),
			wantErr:     true,
			errContains: "HTTP 400",
		},
		{
			name: "HTTP 401 unauthorized",
			cfg: config.ConfigData{
				URL:    "https://api.example.com/chat",
				APIKey: "invalid-key",
			},
			payload:     map[string]interface{}{"model": "gpt-4"},
			mockResp:    makeResponse(401, `{"error":"unauthorized"}`),
			wantErr:     true,
			errContains: "HTTP 401",
		},
		{
			name: "HTTP 429 rate limit",
			cfg: config.ConfigData{
				URL:    "https://api.example.com/chat",
				APIKey: "sk-test123",
			},
			payload:     map[string]interface{}{"model": "gpt-4"},
			mockResp:    makeResponse(429, `{"error":"rate limit exceeded"}`),
			wantErr:     true,
			errContains: "HTTP 429",
		},
		{
			name: "HTTP 500 server error",
			cfg: config.ConfigData{
				URL:    "https://api.example.com/chat",
				APIKey: "sk-test123",
			},
			payload:     map[string]interface{}{"model": "gpt-4"},
			mockResp:    makeResponse(500, `{"error":"internal server error"}`),
			wantErr:     true,
			errContains: "HTTP 500",
		},
		{
			name: "network error",
			cfg: config.ConfigData{
				URL:    "https://api.example.com/chat",
				APIKey: "sk-test123",
			},
			payload:     map[string]interface{}{"model": "gpt-4"},
			mockErr:     http.ErrHandlerTimeout,
			wantErr:     true,
			errContains: "execute request",
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

			got, err := executeHTTP(tt.cfg, tt.payload)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantBody, string(got))
		})
	}
}

func TestExecuteHTTPHeaders(t *testing.T) {
	cfg := config.ConfigData{
		URL:    "https://api.example.com/chat",
		APIKey: "sk-test-key",
	}
	payload := map[string]interface{}{"model": "gpt-4"}

	transport := &mockRoundTripper{
		response: makeResponse(200, `{"result":"ok"}`),
	}
	oldClient := httpClient
	httpClient = &http.Client{
		Timeout:   5 * time.Minute,
		Transport: transport,
	}
	defer func() { httpClient = oldClient }()

	_, err := executeHTTP(cfg, payload)
	assert.NoError(t, err)

	assert.Equal(t, "application/json", transport.request.Header.Get("Content-Type"))
	assert.Equal(t, "Bearer sk-test-key", transport.request.Header.Get("Authorization"))
}

func TestExecuteHTTPTimeout(t *testing.T) {
	cfg := config.ConfigData{
		URL:    "https://api.example.com/chat",
		APIKey: "sk-test123",
	}
	payload := map[string]interface{}{"model": "gpt-4"}

	transport := &mockRoundTripper{
		response: makeResponse(200, `{"ok":true}`),
	}

	client := &http.Client{
		Timeout:   5 * time.Minute,
		Transport: transport,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", cfg.URL, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     ConfigData
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			cfg: ConfigData{
				Protocol:       ProtocolOpenAI,
				URL:            "https://api.openai.com",
				Model:          "gpt-4",
				FallbackModels: []string{"gpt-3.5"},
				APIKey:         "sk-test123",
			},
			wantErr: false,
		},
		{
			name: "missing api key",
			cfg: ConfigData{
				Protocol:       ProtocolOpenAI,
				URL:            "https://api.openai.com",
				Model:          "gpt-4",
				FallbackModels: []string{"gpt-3.5"},
				APIKey:         "",
			},
			wantErr: true,
			errMsg:  "API key required",
		},
		{
			name: "invalid protocol",
			cfg: ConfigData{
				Protocol:       APIProtocol(99),
				URL:            "https://api.openai.com",
				Model:          "gpt-4",
				FallbackModels: []string{"gpt-3.5"},
				APIKey:         "sk-test123",
			},
			wantErr: true,
			errMsg:  "invalid protocol",
		},
		{
			name: "ollama protocol valid",
			cfg: ConfigData{
				Protocol:       ProtocolOllama,
				URL:            "http://localhost:11434",
				Model:          "llama3",
				FallbackModels: []string{},
				APIKey:         "not-used-but-required",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

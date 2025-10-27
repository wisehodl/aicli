package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildConfig(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		env     map[string]string
		wantErr bool
		check   func(*testing.T, ConfigData)
	}{
		{
			name: "valid full config",
			args: []string{"-k", "sk-test", "-m", "gpt-4"},
			check: func(t *testing.T, cfg ConfigData) {
				assert.Equal(t, "sk-test", cfg.APIKey)
				assert.Equal(t, "gpt-4", cfg.Model)
			},
		},
		{
			name: "config file from env",
			args: []string{"-k", "sk-test"},
			env:  map[string]string{"AICLI_CONFIG_FILE": "testdata/partial.yaml"},
			check: func(t *testing.T, cfg ConfigData) {
				assert.Equal(t, "gpt-4", cfg.Model)
			},
		},
		{
			name:    "missing api key",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "invalid config file",
			args:    []string{"-c", "testdata/invalid.yaml", "-k", "test"},
			wantErr: true,
		},
		{
			name:    "invalid protocol in flags",
			args:    []string{"-k", "sk-test", "-l", "invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all AICLI_* env vars
			t.Setenv("AICLI_API_KEY", "")
			t.Setenv("AICLI_API_KEY_FILE", "")
			t.Setenv("AICLI_PROTOCOL", "")
			t.Setenv("AICLI_URL", "")
			t.Setenv("AICLI_MODEL", "")
			t.Setenv("AICLI_FALLBACK", "")
			t.Setenv("AICLI_SYSTEM", "")
			t.Setenv("AICLI_CONFIG_FILE", "")

			// Apply test-specific env
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			cfg, err := BuildConfig(tt.args)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			if tt.check != nil {
				tt.check(t, cfg)
			}
		})
	}
}

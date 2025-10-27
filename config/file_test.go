package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfigFile(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    fileValues
		wantErr bool
	}{
		{
			name: "empty path returns nil",
			path: "",
			want: fileValues{},
		},
		{
			name: "valid config",
			path: "testdata/valid.yaml",
			want: fileValues{
				protocol:   "ollama",
				url:        "http://localhost:11434/api/chat",
				keyFile:    "~/.aicli_key",
				model:      "llama3",
				fallback:   "llama2,mistral",
				systemFile: "~/system.txt",
			},
		},
		{
			name: "partial config",
			path: "testdata/partial.yaml",
			want: fileValues{
				model:    "gpt-4",
				fallback: "gpt-3.5-turbo",
			},
		},
		{
			name: "empty file",
			path: "testdata/empty.yaml",
			want: fileValues{},
		},
		{
			name:    "file not found",
			path:    "testdata/nonexistent.yaml",
			wantErr: true,
		},
		{
			name:    "invalid yaml syntax",
			path:    "testdata/invalid.yaml",
			wantErr: true,
		},
		{
			name: "unknown keys ignored",
			path: "testdata/unknown_keys.yaml",
			want: fileValues{
				protocol: "openai",
				model:    "gpt-4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadConfigFile(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

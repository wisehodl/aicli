package config

import (
	"fmt"
	"os"
)

const UsageText = `Usage: aicli [OPTION]...
Send prompts and files to LLM chat endpoints.

Global:
  --version                display version and exit

Input:
  -f, --file PATH          input file (repeatable)
  -F, --stdin-file         treat stdin as file content
  -p, --prompt TEXT        prompt text (repeatable)
  -pf, --prompt-file PATH  read prompt from file

System:
  -s, --system TEXT        system prompt text
  -sf, --system-file PATH  read system prompt from file
                           (error if both -s and -sf provided)

API:
  -l, --protocol PROTO     openai or ollama (default: openai)
  -u, --url URL            endpoint (default: https://api.ppq.ai/chat/completions)
  -k, --key KEY            API key
  -kf, --key-file PATH     read API key from file

Models:
  -m, --model NAME         primary model (default: gpt-4o-mini)
  -b, --fallback NAMES     comma-separated fallback list (default: gpt-4.1-mini)

Output:
  -o, --output PATH        write to file (mode 0644) instead of stdout
  -q, --quiet              suppress progress messages
  -v, --verbose            log debug information to stderr

Config:
  -c, --config PATH        YAML config file

Environment Variables:
  AICLI_API_KEY            API key
  AICLI_API_KEY_FILE       path to API key file
  AICLI_PROTOCOL           API protocol
  AICLI_URL                endpoint URL
  AICLI_MODEL              primary model name
  AICLI_FALLBACK           comma-separated fallback models
  AICLI_SYSTEM             system prompt text
  AICLI_SYSTEM_FILE        path to system prompt file
  AICLI_CONFIG_FILE        path to config file
  AICLI_PROMPT_FILE        path to prompt file
  AICLI_DEFAULT_PROMPT     override default prompt

Precedence Rules:
  API key:      --key > --key-file > AICLI_API_KEY > AICLI_API_KEY_FILE > config key_file
  System:       --system > --system-file > AICLI_SYSTEM > AICLI_SYSTEM_FILE > config system_file
  Config file:  --config > AICLI_CONFIG_FILE
  All others:   flags > environment > config file > defaults

Stdin Behavior:
  No flags:     stdin becomes the prompt
  With -p/-pf:  stdin appends after explicit prompts
  With -F:      stdin becomes first file (path: "input")

Examples:
  echo "What is Rust?" | aicli
  cat log.txt | aicli -F -p "Find errors in this log"
  aicli -f main.go -p "Review this code"
  aicli -c ~/.aicli.yaml -f src/main.go -f src/util.go -o analysis.md
  aicli -p "Context:" -pf template.txt -p "Apply to finance sector"
`

func printUsage() {
	fmt.Fprint(os.Stderr, UsageText)
}

// BuildConfig resolves configuration from all sources with precedence:
// flags > env > file > defaults
func BuildConfig(args []string) (ConfigData, error) {
	flags, err := parseFlags(args)
	if err != nil {
		return ConfigData{}, fmt.Errorf("parse flags: %w", err)
	}

	// Validate protocol strings before merge
	if flags.protocol != "" && flags.protocol != "openai" && flags.protocol != "ollama" {
		return ConfigData{}, fmt.Errorf("invalid protocol: must be openai or ollama, got: %s", flags.protocol)
	}

	configPath := flags.config
	if configPath == "" {
		configPath = os.Getenv("AICLI_CONFIG_FILE")
	}

	env := loadEnvironment()

	// Validate env protocol
	if env.protocol != "" && env.protocol != "openai" && env.protocol != "ollama" {
		return ConfigData{}, fmt.Errorf("invalid protocol: must be openai or ollama, got: %s", env.protocol)
	}

	file, err := loadConfigFile(configPath)
	if err != nil {
		return ConfigData{}, fmt.Errorf("load config file: %w", err)
	}

	// Validate file protocol
	if file.protocol != "" && file.protocol != "openai" && file.protocol != "ollama" {
		return ConfigData{}, fmt.Errorf("invalid protocol: must be openai or ollama, got: %s", file.protocol)
	}

	cfg := mergeSources(flags, env, file)

	if err := validateConfig(cfg); err != nil {
		return ConfigData{}, err
	}

	return cfg, nil
}

// IsVersionRequest checks if --version flag was passed
func IsVersionRequest(args []string) bool {
	for _, arg := range args {
		if arg == "--version" {
			return true
		}
	}
	return false
}

// IsHelpRequest checks if -h or --help flag was passed
func IsHelpRequest(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}

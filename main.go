package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"git.wisehodl.dev/jay/aicli/version"
	"gopkg.in/yaml.v3"
)

const defaultPrompt = "Analyze the following:"

type stdinRole int

const (
	stdinAsPrompt stdinRole = iota
	stdinAsPrefixedContent
	stdinAsFile
)

type Config struct {
	Protocol   string
	URL        string
	Key        string
	Model      string
	Fallbacks  []string
	SystemText string
	PromptText string
	Files      []FileData
	OutputPath string
	Quiet      bool
	Verbose    bool
}

type FileData struct {
	Path    string
	Content string
}

type flagValues struct {
	files       []string
	prompts     []string
	promptFile  string
	system      string
	systemFile  string
	key         string
	keyFile     string
	protocol    string
	url         string
	model       string
	fallback    string
	output      string
	config      string
	stdinFile   bool
	quiet       bool
	verbose     bool
	showVersion bool
}

const usageText = `Usage: aicli [OPTION]... [FILE]...    
Send files and prompts to LLM chat endpoints.    
    
With no FILE, or when FILE is -, read standard input.    

Global:
  --version            	   display version information and exit
    
Input:    
  -f, --file PATH          input file (repeatable)    
  -F, --stdin-file         treat stdin as file contents
  -p, --prompt TEXT        prompt text (repeatable, can be combined with --prompt-file)    
  -pf, --prompt-file PATH  prompt from file (combined with any --prompt flags)    
    
System:    
  -s, --system TEXT        system prompt text    
  -sf, --system-file PATH  system prompt from file    
    
API:    
  -l, --protocol PROTO     API protocol: openai, ollama (default: openai)    
  -u, --url URL            API endpoint (default: https://api.ppq.ai/chat/completions)    
  -k, --key KEY            API key (if present, --key-file is ignored)    
  -kf, --key-file PATH     API key from file (used only if --key is not provided)    
    
Models:    
  -m, --model NAME         primary model (default: gpt-4o-mini)    
  -b, --fallback NAMES     comma-separated fallback models (default: gpt-4.1-mini)    
    
Output:    
  -o, --output PATH        write to file instead of stdout    
  -q, --quiet              suppress progress output    
  -v, --verbose            enable debug logging    
    
Config:    
  -c, --config PATH        YAML config file    
    
Environment variables:    
  AICLI_API_KEY        API key    
  AICLI_API_KEY_FILE   Path to file containing API key (used only if AICLI_API_KEY is not set)
  AICLI_PROTOCOL       API protocol    
  AICLI_URL            API endpoint    
  AICLI_MODEL          primary model    
  AICLI_FALLBACK       fallback models    
  AICLI_SYSTEM         system prompt    
  AICLI_DEFAULT_PROMPT default prompt override  
  AICLI_CONFIG_FILE    Path to config file  
  AICLI_PROMPT_FILE    Path to prompt file  
  AICLI_SYSTEM_FILE    Path to system file  
    
API Key precedence: --key flag > --key-file flag > AICLI_API_KEY > AICLI_API_KEY_FILE > config file
    
Examples:    
  echo "What is Rust?" | aicli    
  cat file.txt | aicli -F -p "Analyze this file"
  aicli -f main.go -p "Review this code"    
  aicli -c ~/.aicli.yaml -f src/*.go -o analysis.md    
  aicli -p "First prompt" -pf prompt.txt -p "Last prompt"    
`

func printUsage() {
	fmt.Fprint(os.Stderr, usageText)
}

type fileList []string

func (f *fileList) String() string {
	return strings.Join(*f, ", ")
}

func (f *fileList) Set(value string) error {
	*f = append(*f, value)
	return nil
}

type promptList []string

func (p *promptList) String() string {
	return strings.Join(*p, "\n")
}

func (p *promptList) Set(value string) error {
	*p = append(*p, value)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Check for verbose flag early
	verbose := false
	for _, arg := range os.Args {
		if arg == "-v" || arg == "--verbose" {
			verbose = true
			break
		}
	}

	// Check for config file in environment variable before parsing flags
	configFilePath := os.Getenv("AICLI_CONFIG_FILE")

	flags := parseFlags()

	if flags.showVersion {
		fmt.Printf("aicli %s\n", version.GetVersion())
		return nil
	}

	if flags.config == "" && configFilePath != "" {
		flags.config = configFilePath
	}

	envVals := loadEnvVars(verbose)
	fileVals, err := loadConfigFile(flags.config)
	if err != nil {
		return err
	}

	merged := mergeConfigSources(verbose, flags, envVals, fileVals)
	if err := validateConfig(merged); err != nil {
		return err
	}

	if promptFilePath := os.Getenv("AICLI_PROMPT_FILE"); promptFilePath != "" && flags.promptFile == "" {
		content, err := os.ReadFile(promptFilePath)
		if err != nil {
			if verbose {
				fmt.Fprintf(os.Stderr, "[verbose] Failed to read AICLI_PROMPT_FILE at %s: %v\n", promptFilePath, err)
			}
		} else {
			merged.PromptText = string(content)
		}
	}

	if systemFilePath := os.Getenv("AICLI_SYSTEM_FILE"); systemFilePath != "" && flags.systemFile == "" && flags.system == "" {
		content, err := os.ReadFile(systemFilePath)
		if err != nil {
			if verbose {
				fmt.Fprintf(os.Stderr, "[verbose] Failed to read AICLI_SYSTEM_FILE at %s: %v\n", systemFilePath, err)
			}
		} else {
			merged.SystemText = string(content)
		}
	}

	stdinContent, hasStdin := detectStdin()
	role := determineStdinRole(flags, hasStdin)

	inputData, err := resolveInputStreams(merged, stdinContent, hasStdin, role, flags)
	if err != nil {
		return err
	}

	config := buildCompletePrompt(inputData)

	if config.Verbose {
		logVerbose("Configuration resolved", config)
	}

	startTime := time.Now()
	response, usedModel, err := sendChatRequest(config)
	duration := time.Since(startTime)

	if err != nil {
		return err
	}

	return writeOutput(response, usedModel, duration, config)
}

func parseFlags() flagValues {
	fv := flagValues{}
	var files fileList
	var prompts promptList

	flag.Usage = printUsage

	flag.Var(&files, "f", "")
	flag.Var(&files, "file", "")
	flag.Var(&prompts, "p", "")
	flag.Var(&prompts, "prompt", "")
	flag.StringVar(&fv.promptFile, "pf", "", "")
	flag.StringVar(&fv.promptFile, "prompt-file", "", "")
	flag.StringVar(&fv.system, "s", "", "")
	flag.StringVar(&fv.system, "system", "", "")
	flag.StringVar(&fv.systemFile, "sf", "", "")
	flag.StringVar(&fv.systemFile, "system-file", "", "")
	flag.StringVar(&fv.key, "k", "", "")
	flag.StringVar(&fv.key, "key", "", "")
	flag.StringVar(&fv.keyFile, "kf", "", "")
	flag.StringVar(&fv.keyFile, "key-file", "", "")
	flag.StringVar(&fv.protocol, "l", "", "")
	flag.StringVar(&fv.protocol, "protocol", "", "")
	flag.StringVar(&fv.url, "u", "", "")
	flag.StringVar(&fv.url, "url", "", "")
	flag.StringVar(&fv.model, "m", "", "")
	flag.StringVar(&fv.model, "model", "", "")
	flag.StringVar(&fv.fallback, "b", "", "")
	flag.StringVar(&fv.fallback, "fallback", "", "")
	flag.StringVar(&fv.output, "o", "", "")
	flag.StringVar(&fv.output, "output", "", "")
	flag.StringVar(&fv.config, "c", "", "")
	flag.StringVar(&fv.config, "config", "", "")
	flag.BoolVar(&fv.stdinFile, "F", false, "")
	flag.BoolVar(&fv.stdinFile, "stdin-file", false, "")
	flag.BoolVar(&fv.quiet, "q", false, "")
	flag.BoolVar(&fv.quiet, "quiet", false, "")
	flag.BoolVar(&fv.verbose, "v", false, "")
	flag.BoolVar(&fv.verbose, "verbose", false, "")
	flag.BoolVar(&fv.showVersion, "version", false, "")

	flag.Parse()

	fv.files = files
	fv.prompts = prompts

	return fv
}

func loadEnvVars(verbose bool) map[string]string {
	env := make(map[string]string)
	if val := os.Getenv("AICLI_PROTOCOL"); val != "" {
		env["protocol"] = val
	}
	if val := os.Getenv("AICLI_URL"); val != "" {
		env["url"] = val
	}
	if val := os.Getenv("AICLI_API_KEY"); val != "" {
		env["key"] = val
	}
	if env["key"] == "" {
		if val := os.Getenv("AICLI_API_KEY_FILE"); val != "" {
			content, err := os.ReadFile(val)
			if err != nil && verbose {
				fmt.Fprintf(os.Stderr, "[verbose] Failed to read AICLI_API_KEY_FILE at %s: %v\n", val, err)
			} else {
				env["key"] = strings.TrimSpace(string(content))
			}
		}
	}
	if val := os.Getenv("AICLI_MODEL"); val != "" {
		env["model"] = val
	}
	if val := os.Getenv("AICLI_FALLBACK"); val != "" {
		env["fallback"] = val
	}
	if val := os.Getenv("AICLI_SYSTEM"); val != "" {
		env["system"] = val
	}
	if val := os.Getenv("AICLI_DEFAULT_PROMPT"); val != "" {
		env["prompt"] = val
	}
	return env
}

func loadConfigFile(path string) (map[string]interface{}, error) {
	if path == "" {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	return config, nil
}

func mergeConfigSources(verbose bool, flags flagValues, env map[string]string, file map[string]interface{}) Config {
	cfg := Config{
		Protocol:  "openai",
		URL:       "https://api.ppq.ai/chat/completions",
		Model:     "gpt-4o-mini",
		Fallbacks: []string{"gpt-4.1-mini"},
		Quiet:     flags.quiet,
		Verbose:   flags.verbose,
	}

	if env["protocol"] != "" {
		cfg.Protocol = env["protocol"]
	}
	if env["url"] != "" {
		cfg.URL = env["url"]
	}
	if env["key"] != "" {
		cfg.Key = env["key"]
	}
	if env["model"] != "" {
		cfg.Model = env["model"]
	}
	if env["fallback"] != "" {
		cfg.Fallbacks = strings.Split(env["fallback"], ",")
	}
	if env["system"] != "" {
		cfg.SystemText = env["system"]
	}

	if file != nil {
		if v, ok := file["protocol"].(string); ok {
			cfg.Protocol = v
		}
		if v, ok := file["url"].(string); ok {
			cfg.URL = v
		}
		if v, ok := file["model"].(string); ok {
			cfg.Model = v
		}
		if v, ok := file["fallback"].(string); ok {
			cfg.Fallbacks = strings.Split(v, ",")
		}
		if v, ok := file["system_file"].(string); ok {
			content, err := os.ReadFile(v)
			if err != nil {
				if verbose {
					fmt.Fprintf(os.Stderr, "[verbose] Failed to read system_file at %s: %v\n", v, err)
				}
			} else {
				cfg.SystemText = string(content)
			}
		}
		if v, ok := file["key_file"].(string); ok && cfg.Key == "" {
			content, err := os.ReadFile(v)
			if err != nil {
				if cfg.Verbose {
					fmt.Fprintf(os.Stderr, "[verbose] Failed to read key_file at %s: %v\n", v, err)
				}
			} else {
				cfg.Key = strings.TrimSpace(string(content))
			}
		}
	}

	if flags.protocol != "" {
		cfg.Protocol = flags.protocol
	}
	if flags.url != "" {
		cfg.URL = flags.url
	}
	if flags.model != "" {
		cfg.Model = flags.model
	}
	if flags.fallback != "" {
		cfg.Fallbacks = strings.Split(flags.fallback, ",")
	}
	if flags.system != "" {
		cfg.SystemText = flags.system
	}
	if flags.systemFile != "" {
		content, err := os.ReadFile(flags.systemFile)
		if err != nil {
			if cfg.Verbose {
				fmt.Fprintf(os.Stderr, "[verbose] Failed to read system file at %s: %v\n", flags.systemFile, err)
			}
		} else {
			cfg.SystemText = string(content)
		}
	}
	if flags.key != "" {
		cfg.Key = flags.key
	} else if flags.keyFile != "" {
		content, err := os.ReadFile(flags.keyFile)
		if err != nil {
			if cfg.Verbose {
				fmt.Fprintf(os.Stderr, "[verbose] Failed to read key file at %s: %v\n", flags.keyFile, err)
			}
		} else {
			cfg.Key = strings.TrimSpace(string(content))
		}
	}
	if flags.output != "" {
		cfg.OutputPath = flags.output
	}

	return cfg
}

func validateConfig(cfg Config) error {
	if cfg.Key == "" {
		return fmt.Errorf("API key required: use --key, --key-file, AICLI_API_KEY, AICLI_API_KEY_FILE, or key_file in config")
	}
	if cfg.Protocol != "openai" && cfg.Protocol != "ollama" {
		return fmt.Errorf("protocol must be 'openai' or 'ollama', got: %s", cfg.Protocol)
	}
	return nil
}

func detectStdin() (string, bool) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", false
	}

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", false
	}

	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", false
	}

	return string(content), true
}

func determineStdinRole(flags flagValues, hasStdin bool) stdinRole {
	if !hasStdin {
		return stdinAsPrompt
	}

	if flags.stdinFile {
		return stdinAsFile
	}

	hasExplicitPrompt := len(flags.prompts) > 0 || flags.promptFile != ""

	if hasExplicitPrompt {
		return stdinAsPrefixedContent
	}

	return stdinAsPrompt
}

func resolveInputStreams(cfg Config, stdinContent string, hasStdin bool, role stdinRole, flags flagValues) (Config, error) {
	hasPromptFlag := len(flags.prompts) > 0 || flags.promptFile != ""
	hasFileFlag := len(flags.files) > 0

	// Handle case where only stdin as file is provided
	if !hasPromptFlag && !hasFileFlag && hasStdin && flags.stdinFile {
		cfg.Files = append(cfg.Files, FileData{Path: "input", Content: stdinContent})
		return cfg, nil
	}

	if !hasStdin && !hasFileFlag && !hasPromptFlag {
		return cfg, fmt.Errorf("no input provided: supply stdin, --file, or --prompt")
	}

	for _, path := range flags.files {
		if path == "" {
			return cfg, fmt.Errorf("empty file path provided")
		}
	}

	if flags.system != "" && flags.systemFile != "" {
		return cfg, fmt.Errorf("cannot use both --system and --system-file")
	}

	if len(flags.prompts) > 0 {
		cfg.PromptText = strings.Join(flags.prompts, "\n")
	}

	if flags.promptFile != "" {
		content, err := os.ReadFile(flags.promptFile)
		if err != nil {
			return cfg, fmt.Errorf("read prompt file: %w", err)
		}
		if cfg.PromptText != "" {
			cfg.PromptText += "\n\n" + string(content)
		} else {
			cfg.PromptText = string(content)
		}
	}

	if hasStdin {
		switch role {
		case stdinAsPrompt:
			cfg.PromptText = stdinContent
		case stdinAsPrefixedContent:
			if cfg.PromptText != "" {
				cfg.PromptText += "\n\n" + stdinContent
			} else {
				cfg.PromptText = stdinContent
			}
		case stdinAsFile:
			cfg.Files = append(cfg.Files, FileData{Path: "input", Content: stdinContent})
		}
	}

	for _, path := range flags.files {
		content, err := os.ReadFile(path)
		if err != nil {
			return cfg, fmt.Errorf("read file %s: %w", path, err)
		}
		cfg.Files = append(cfg.Files, FileData{Path: path, Content: string(content)})
	}

	return cfg, nil
}

func buildCompletePrompt(inputData Config) Config {
	result := inputData
	promptParts := []string{}

	// Use inputData's prompt if set, otherwise check for overrides
	if inputData.PromptText != "" {
		promptParts = append(promptParts, inputData.PromptText)
	} else if override := os.Getenv("AICLI_DEFAULT_PROMPT"); override != "" {
		promptParts = append(promptParts, override)
	} else if len(inputData.Files) > 0 {
		promptParts = append(promptParts, defaultPrompt)
	}

	// Format files if present
	if len(inputData.Files) > 0 {
		fileSection := formatFiles(inputData.Files)
		if len(promptParts) > 0 {
			promptParts = append(promptParts, "", fileSection)
		} else {
			promptParts = append(promptParts, fileSection)
		}
	}

	result.PromptText = strings.Join(promptParts, "\n")
	return result
}

func formatFiles(files []FileData) string {
	var buf strings.Builder
	for i, f := range files {
		if i > 0 {
			buf.WriteString("\n\n")
		}
		buf.WriteString(fmt.Sprintf("File: %s\n\n```\n%s\n```", f.Path, f.Content))
	}
	return buf.String()
}

func sendChatRequest(cfg Config) (string, string, error) {
	models := append([]string{cfg.Model}, cfg.Fallbacks...)

	for i, model := range models {
		if !cfg.Quiet && i > 0 {
			fmt.Fprintf(os.Stderr, "Model %s failed, trying %s...\n", models[i-1], model)
		}

		response, err := tryModel(cfg, model)
		if err == nil {
			return response, model, nil
		}

		if !cfg.Quiet {
			fmt.Fprintf(os.Stderr, "Model %s failed: %v\n", model, err)
		}
	}

	return "", "", fmt.Errorf("all models failed")
}

func tryModel(cfg Config, model string) (string, error) {
	payload := buildPayload(cfg, model)
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}

	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "Request payload: %s\n", string(body))
	}

	req, err := http.NewRequest("POST", cfg.URL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.Key))

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if cfg.Verbose {
		fmt.Fprintf(os.Stderr, "Response: %s\n", string(respBody))
	}

	return parseResponse(respBody, cfg.Protocol)
}

func buildPayload(cfg Config, model string) map[string]interface{} {
	if cfg.Protocol == "ollama" {
		payload := map[string]interface{}{
			"model":  model,
			"prompt": cfg.PromptText,
			"stream": false,
		}
		if cfg.SystemText != "" {
			payload["system"] = cfg.SystemText
		}
		return payload
	}

	messages := []map[string]string{}
	if cfg.SystemText != "" {
		messages = append(messages, map[string]string{
			"role":    "system",
			"content": cfg.SystemText,
		})
	}
	messages = append(messages, map[string]string{
		"role":    "user",
		"content": cfg.PromptText,
	})

	return map[string]interface{}{
		"model":    model,
		"messages": messages,
	}
}

func parseResponse(body []byte, protocol string) (string, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	if protocol == "ollama" {
		if response, ok := result["response"].(string); ok {
			return response, nil
		}
		return "", fmt.Errorf("no response field in ollama response")
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid choice format")
	}

	message, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("no message in choice")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("no content in message")
	}

	return content, nil
}

func writeOutput(response, model string, duration time.Duration, cfg Config) error {
	if cfg.OutputPath == "" {
		if !cfg.Quiet {
			fmt.Println("--- aicli ---")
			fmt.Println()
			fmt.Printf("Used model: %s\n", model)
			fmt.Printf("Query duration: %.1fs\n", duration.Seconds())
			fmt.Println()
			fmt.Println("--- response ---")
			fmt.Println()
		}
		fmt.Println(response)
		return nil
	}

	if err := os.WriteFile(cfg.OutputPath, []byte(response), 0644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}

	if !cfg.Quiet {
		fmt.Printf("Used model: %s\n", model)
		fmt.Printf("Query duration: %.1fs\n", duration.Seconds())
		fmt.Printf("Wrote response to: %s\n", cfg.OutputPath)
	}

	return nil
}

func logVerbose(msg string, cfg Config) {
	fmt.Fprintf(os.Stderr, "[verbose] %s\n", msg)
	fmt.Fprintf(os.Stderr, "  Protocol: %s\n", cfg.Protocol)
	fmt.Fprintf(os.Stderr, "  URL: %s\n", cfg.URL)
	fmt.Fprintf(os.Stderr, "  Model: %s\n", cfg.Model)
	fmt.Fprintf(os.Stderr, "  Fallbacks: %v\n", cfg.Fallbacks)
	fmt.Fprintf(os.Stderr, "  Files: %d\n", len(cfg.Files))
}

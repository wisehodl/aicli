# AICLI

Source: https://git.wisehodl.dev/jay/aicli

Mirror: https://github.com/wisehodl/aicli

A flexible command-line interface for interacting with LLM chat APIs.

AICLI provides a streamlined way to interact with language models from your terminal. Send prompts and files to chat models like OpenAI's GPT or local Ollama models, customize system prompts, and receive responses right in your terminal or save them to files.

## Features

- Query OpenAI-compatible APIs or Ollama models directly
- Send files as context with your prompts
- Customize system prompts
- Configure via environment variables, config files, or CLI flags
- Save responses to files
- Automatic model fallbacks if primary models fail
- Flexible input handling (stdin, files, direct prompts)

## Installation

### Pre-built Binaries

Download the latest binary for your platform from the [Releases](https://git.wisehodl.dev/jay/aicli/releases) page:

Make the file executable (Linux/macOS):

```bash
chmod +x aicli-linux-amd64
mv aicli-linux-amd64 /usr/local/bin/aicli  # or any directory in your PATH
```

### Building from Source

Requires Go 1.16+:

```bash
git clone https://git.wisehodl.dev/jay/aicli.git
cd aicli
go build -o aicli
```

## Configuration

AICLI can be configured in multiple ways, with the following precedence (highest to lowest):

1. Command-line flags
2. Environment variables
3. Config file
4. Default values

### API Key Setup

Set up your API key using one of these methods:

```bash
# Command line
aicli --key "your-api-key" ...

# Environment variable
export AICLI_API_KEY="your-api-key"

# Key file
echo "your-api-key" > ~/.aicli_key
export AICLI_API_KEY_FILE=~/.aicli_key
```

### Environment Variables

```bash
# API Configuration
export AICLI_API_KEY="your-api-key"
export AICLI_API_KEY_FILE="~/.aicli_key"
export AICLI_PROTOCOL="openai"  # or "ollama"
export AICLI_URL="https://api.ppq.ai/chat/completions"  # custom endpoint

# Model Selection
export AICLI_MODEL="gpt-4o-mini"
export AICLI_FALLBACK="gpt-4.1-mini,gpt-3.5-turbo"

# Prompts
export AICLI_SYSTEM="You are a helpful AI assistant."
export AICLI_DEFAULT_PROMPT="Analyze the following:"

# File Paths
export AICLI_CONFIG_FILE="~/.aicli.yaml"
export AICLI_PROMPT_FILE="~/prompts/default.txt"
export AICLI_SYSTEM_FILE="~/prompts/system.txt"
```

### Config File (YAML)

Create a YAML config file (e.g., `~/.aicli.yaml`):

```yaml
protocol: openai
url: https://api.ppq.ai/chat/completions
model: gpt-4o-mini
fallback: gpt-4.1-mini,gpt-3.5-turbo
key_file: ~/.aicli_key
system_file: ~/prompts/system.txt
```

## Basic Usage

### Simple Queries

```bash
# Direct question
aicli -p "Explain quantum computing in simple terms"

# Using stdin
echo "What is the capital of France?" | aicli

# Save response to file
aicli -p "Write a short poem about coding" -o poem.txt
```

### Working with Files

```bash
# Analyze a code file
aicli -f main.go -p "Review this code for bugs and improvements"

# Analyze multiple files
aicli -f main.go -f utils.go -p "Explain how these files work together"

# Using stdin as a file
cat log.txt | aicli -F -p "Find problems in this log"

# Combining stdin file with regular files
cat log.txt | aicli -F -f config.json -p "Find problems in this log and config"
```

### Customizing Prompts

```bash
# Multiple prompt sections
aicli -p "Analyze this data:" -p "Focus on trends over time:" -f data.csv

# Combining prompts from files and flags
aicli -pf prompt_template.txt -p "Apply this to the finance sector" -f report.txt

# Using system prompt
aicli -s "You are a security expert" -p "Review this code for vulnerabilities" -f app.js
```

### API Configuration

```bash
# Using Ollama with local model
aicli -l ollama -u http://localhost:11434/api/chat -m llama3 -p "Explain Docker"

# Custom OpenAI-compatible endpoint
aicli -u https://api.company.ai/v1/chat/completions -p "Generate a marketing slogan"

# With fallback models
aicli -m claude-3-opus -b claude-3-sonnet,gpt-4o -p "Write a complex algorithm"
```

## Advanced Examples

### Code Review Workflow

```bash
# Review pull request changes
git diff main..feature-branch | aicli -p "Review these changes. Identify potential bugs and suggest improvements."
```

### Data Analysis

```bash
# Analyze CSV data
aicli -p "Analyze this CSV data and provide insights:" -f data.csv -o analysis.md
```

### Content Generation with Context

```bash
# Generate documentation with multiple context files
aicli -s "You are a technical writer creating clear documentation" \
     -p "Create a README for this project" \
     -f main.go -f api.go -f config.go -o README.md
```

### Translation Pipeline

```bash
# Translate text from a file
cat text.txt | aicli -p "Translate this text to Spanish" > text_es.txt
```

### Meeting Summarization

```bash
# Summarize meeting transcript
aicli -pf summarization_prompt.txt -f transcript.txt -o summary.md
```

## Tips and Tricks

### Creating Reusable Prompt Files

Store common prompt patterns in files:

```
# ~/prompts/code-review.txt
Review the following code:
1. Identify potential bugs or edge cases
2. Suggest performance improvements
3. Highlight security concerns
4. Recommend better patterns or practices
```

Then use them:

```bash
aicli -pf ~/prompts/code-review.txt -f main.go
```

### Environment Variable Shortcuts

Set up aliases with pre-configured environment variables:

```bash
# In your .bashrc or .zshrc
alias aicli-code="AICLI_SYSTEM_FILE=~/prompts/system-code.txt aicli"
alias aicli-creative="AICLI_SYSTEM='You are a creative writer' AICLI_MODEL=gpt-4o aicli"
```

### Combining with Other Tools

```bash
# Generate commit message from changes
git diff --staged | aicli -q -p "Generate a concise commit message for these changes:" | git commit -F -

# Analyze logs
grep ERROR /var/log/app.log | aicli -p "Identify patterns in these error logs"
```

## Full Command Reference

```
Usage: aicli [OPTION]... [FILE]...
Send files and prompts to LLM chat endpoints.

With no FILE, or when FILE is -, read standard input.

Global:
  --version                display version information and exit

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
```

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

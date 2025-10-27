# Installing AICLI

This document provides instructions for installing and configuring AICLI on various platforms.

## Pre-built Binaries

The easiest way to install AICLI is to download a pre-built binary from the [releases page](https://git.wisehodl.dev/jay/aicli/releases).

### Platform Selection Guide

Choose the appropriate binary for your platform:

- **Linux (64-bit x86)**: `aicli-linux-amd64`
- **Linux (32-bit x86)**: `aicli-linux-386`
- **Linux (64-bit ARM)**: `aicli-linux-arm64`
- **Linux (ARMv7)**: `aicli-linux-armv7`
- **Linux (ARMv6)**: `aicli-linux-armv6`
- **macOS (Intel)**: `aicli-darwin-amd64`
- **macOS (Apple Silicon)**: `aicli-darwin-arm64`
- **Windows (64-bit)**: `aicli-windows-amd64.exe`
- **Windows (32-bit)**: `aicli-windows-386.exe`
- **FreeBSD (64-bit)**: `aicli-freebsd-amd64`
- **OpenBSD (64-bit)**: `aicli-openbsd-amd64`
- **NetBSD (64-bit)**: `aicli-netbsd-amd64`
- **Solaris (64-bit)**: `aicli-solaris-amd64`

### Linux/macOS Installation

```bash
# Download the appropriate binary (replace with actual version and platform)
curl -LO https://git.wisehodl.dev/jay/aicli/releases/download/v1.0.0/aicli-linux-amd64

# Make executable
chmod +x aicli-linux-amd64

# Move to a directory in your PATH
sudo mv aicli-linux-amd64 /usr/local/bin/aicli

# Verify installation
aicli --version
````

### Windows Installation

1. Download the appropriate EXE file for your system
2. Rename the executable to `aicli.exe` if desired
3. Add the directory to your PATH or move the executable to a directory in your PATH
4. Open Command Prompt or PowerShell and verify the installation:

```
aicli --version
```

## Configuration

### API Key Setup

You'll need an API key to use AICLI. Set it up using one of these methods:

```bash
# Direct method (less secure)
export AICLI_API_KEY="your-api-key"

# File method (more secure)
echo "your-api-key" > ~/.aicli_key
export AICLI_API_KEY_FILE=~/.aicli_key

# Or specify in config file (see below)
```

### Configuration File

Create a configuration file at `~/.aicli.yaml` or use the sample config:

```bash
# Create config file
cat > ~/.aicli.yaml << 'EOF'
protocol: openai
url: https://api.ppq.ai/chat/completions
key_file: ~/.aicli_key
model: gpt-4o-mini
fallback: gpt-4.1-mini,o3
EOF

# Edit with your preferred editor
nano ~/.aicli.yaml
```

See the README for detailed configuration options.

## Verification

Verify the downloaded files against the provided checksums:

```bash
# Download the checksum file
curl -LO https://git.wisehodl.dev/jay/aicli/releases/download/v1.0.0/SHA256SUMS

# Verify your downloaded binary
sha256sum -c SHA256SUMS --ignore-missing
```

## Building from Source

If you prefer to build from source:

```bash
# Clone the repository
git clone https://git.wisehodl.dev/jay/aicli.git
cd aicli

# Build
go build -o aicli

# Install
sudo mv aicli /usr/local/bin/
```

## Next Steps

See the [README.md](https://claude.ai/chat/README.md) for usage instructions and examples.

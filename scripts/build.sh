#!/usr/bin/env bash
set -euo pipefail

# Build script for AICLI
# This script builds binaries for multiple platforms and architectures

# Set default values if not provided by environment
: "${VERSION:="dev"}"
: "${PACKAGE:="git.wisehodl.dev/jay/aicli"}"
: "${DATE:=$(date -u +"%Y-%m-%dT%H:%M:%SZ")}"
: "${COMMIT:=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")}"

# Output directories
DIST_DIR="$(pwd)/dist"
mkdir -p "$DIST_DIR"

# Build info
LDFLAGS="-s -w -X '${PACKAGE}/version.Version=${VERSION}' -X '${PACKAGE}/version.CommitHash=${COMMIT}' -X '${PACKAGE}/version.BuildDate=${DATE}'"

echo "Building AICLI version ${VERSION} (${COMMIT}) built at ${DATE}"

TARGETS=(
  "linux/amd64"
  "linux/arm64"
  "linux/386"
  "linux/arm/7" # ARMv7
  "linux/arm/6" # ARMv6
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
  "windows/386"
  "freebsd/amd64"
  "openbsd/amd64"
  "netbsd/amd64"
  "solaris/amd64"
)

# Build all targets
for target in "${TARGETS[@]}"; do
  os=$(echo "$target" | cut -d/ -f1)
  arch=$(echo "$target" | cut -d/ -f2)
  arm_version=""

  # Handle ARM version if specified
  if [[ "$target" == */* && "$target" != */amd64 && "$target" != */386 && "$target" != */arm64 ]]; then
    arm_version=$(echo "$target" | cut -d/ -f3)
    echo "Building for ${os}/${arch} (ARM version ${arm_version})"
  else
    echo "Building for ${os}/${arch}"
  fi

  # Set output filename
  if [[ "$os" == "windows" ]]; then
    output="${DIST_DIR}/aicli-${os}-${arch}.exe"
  elif [[ -n "$arm_version" ]]; then
    output="${DIST_DIR}/aicli-${os}-armv${arm_version}"
  else
    output="${DIST_DIR}/aicli-${os}-${arch}"
  fi

  # Set GOOS, GOARCH, and GOARM
  export GOOS=$os
  export GOARCH=$arch
  if [[ -n "$arm_version" ]]; then
    export GOARM=$arm_version
  else
    unset GOARM
  fi

  # Build the binary
  echo "Building ${output}..."
  go build -ldflags "${LDFLAGS}" -o "$output" .

  # Make binary executable
  if [[ "$os" != "windows" ]]; then
    chmod +x "$output"
  fi
done

echo "All binaries built to ${DIST_DIR}"

# Generate checksums
echo "Generating checksums..."
(cd "$DIST_DIR" && sha256sum ./* >SHA256SUMS)

echo "Build complete! Binaries and checksums available in ${DIST_DIR}"

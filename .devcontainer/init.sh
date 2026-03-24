#!/usr/bin/env bash
set -euo pipefail

if [ -f go.mod ]; then
  echo "Downloading Go modules..."
  go mod download
fi

echo "Tool versions:"
go version
golangci-lint --version | head -n 1
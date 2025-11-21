#!/usr/bin/env bash
set -euo pipefail

MODE="${1:-build}" # "build" (default) or "prepare"

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
BUILD_DIR="${ROOT_DIR}/_build"
COREDNS_DIR="${BUILD_DIR}/coredns"

mkdir -p "${BUILD_DIR}"

if [[ ! -d "${COREDNS_DIR}" ]]; then
  git clone --depth=1 https://github.com/coredns/coredns "${COREDNS_DIR}"
fi

cd "${COREDNS_DIR}"

# Wire the external plugin into plugin.cfg if not present
if ! grep -q '^llm:' plugin.cfg; then
  awk '1; $0 ~ /^template:template$/ { print "llm:github.com/ville/coredns-llm/plugin/llm" }' plugin.cfg > plugin.cfg.new
  mv plugin.cfg.new plugin.cfg
fi

# Regenerate plugin directives and imports so CoreDNS knows about the llm directive
make gen

# Ensure local module replace
if ! grep -q 'github.com/ville/coredns-llm' go.mod; then
  go mod edit -require=github.com/ville/coredns-llm@v0.0.0-00010101000000-000000000000
fi
go mod edit -replace=github.com/ville/coredns-llm="${ROOT_DIR}"

if [[ "${MODE}" == "build" ]]; then
  echo "Building CoreDNS with llm plugin..."
  go build -o coredns ./
  echo "Built: ${COREDNS_DIR}/coredns"
fi

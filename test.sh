#!/usr/bin/env bash
set -euo pipefail

go test -timeout 10s -cover -v ./...

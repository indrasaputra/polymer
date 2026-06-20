#!/usr/bin/env bash

set -eo pipefail

# sqlfluff fix -d postgres -e LT05,AM04,AM09,RF04
config=$(pwd)/.golangci.yml
golangci-lint run  --config=${config} ./...

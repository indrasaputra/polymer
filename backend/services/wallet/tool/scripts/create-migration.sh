#!/usr/bin/env bash

set -eo pipefail

atlas migrate diff "$1" --env local
sqlfluff fix --show-lint-violations -d postgres -e LT05,AM04,AM09,PG01
atlas migrate hash --env local

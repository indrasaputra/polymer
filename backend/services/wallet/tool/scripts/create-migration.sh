#!/usr/bin/env bash

set -eo pipefail

atlas migrate diff --dir file://db/migrations --dev-url docker://postgres/17 --to file://db/schema.sql
sqlfluff fix -d postgres -e LT05,AM04
atlas migrate hash --dir file://db/migrations

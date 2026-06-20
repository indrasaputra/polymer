#!/usr/bin/env bash

set -eo pipefail

run_fieldalignment() {
    local dir=$1
    local has_changes=1

    while [[ $has_changes -ne 0 ]]; do
        echo "running fieldalignment..."
        if output=$(cd "${dir}" && fieldalignment -fix ./... 2>&1); then
            has_changes=0
        else
            echo "fieldalignment detected issues, retrying..."
            echo "$output"
        fi
    done
}

run_fieldalignment

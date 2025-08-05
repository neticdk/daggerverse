#!/usr/bin/env bash

set -eu

test() {
    pwd
    find ./ -type f -name "test.sh" -print0 | while IFS= read -r -d '' test_script; do
        (
            dir=$(dirname "$test_script")
            echo "--> Running test in $dir"
            cd "$dir" && bash ./"$(basename "$test_script")"
        )
    done
}

"$@"
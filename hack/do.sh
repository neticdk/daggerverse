#!/usr/bin/env bash

set -eu

test() {
    find ./ -type f -name "test.sh" -print0 | while IFS= read -r -d '' test_script; do
        (
            dir=$(dirname "$test_script")
            echo "--> Running test in $dir"
            cd "$dir" && bash ./"$(basename "$test_script")"
        )
    done
}

dagger_version() {
    local version="0.0.0"
    while IFS= read -r -d '' test_script; do
        local dir
        dir=$(dirname "$test_script")
        if [ -f "$dir/dagger.json" ]; then
            v=$(yq '.engineVersion' "$dir/dagger.json" | sed 's/"//g' | sed 's/v//')
            if [[ -z "$v" || "$v" == "null" ]]; then
                continue
            fi
            set +e
            vercomp "$version" "$v"
            ret=$?
            set -e
            if [ $ret -eq 2 ]; then
                version="$v"
            fi
        fi
    done < <(find ./ -type f -name "test.sh" -print0)
    echo "$version"
}

vercomp () {
    if [[ $1 == $2 ]]
    then
        return 0
    fi
    local IFS=.
    local i ver1=($1) ver2=($2)
    # fill empty fields in ver1 with zeros
    for ((i=${#ver1[@]}; i<${#ver2[@]}; i++))
    do
        ver1[i]=0
    done
    for ((i=0; i<${#ver1[@]}; i++))
    do
        if ((10#${ver1[i]:=0} > 10#${ver2[i]:=0}))
        then
            return 1
        fi
        if ((10#${ver1[i]} < 10#${ver2[i]}))
        then
            return 2
        fi
    done
    return 0
}

"$@"
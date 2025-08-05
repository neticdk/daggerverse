#!/usr/bin/env bash

set -eu

test() {
    for dir in $(ls ./modules); do
        (
            if [ ! -f "./modules/$dir/hack/test.sh" ]; then
                echo "Missing test script for $dir"
                return 1
            fi
            echo "--> Running test in $dir"
            cd "./modules/$dir/hack" && bash test.sh
        )
    done
}

dagger_version() {
    local version="0.0.0"
    for dir in $(ls ./modules); do
        dagger_file="./modules/$dir/dagger.json"
        if [ ! -f "$dagger_file" ]; then
            echo "No dagger.json file in module $dir"
            return 1
        fi
        v=$(yq '.engineVersion' "$dagger_file" | sed 's/"//g' | sed 's/v//')
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
    done
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
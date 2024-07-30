#!/bin/sh

cmds=(jq tr cut)
for cmd in $cmds; do
    if ! command -v $cmd &>/dev/null; then
        echo "[ERROR] $cmd command not found"
        exit 1
    fi
done


JSON_DATA_FILE="$1"
OUTPUT_FILE="$2"

jq 'map(.uri)' $JSON_DATA_FILE | jq '.[]' | tr -d '"' > $OUTPUT_FILE

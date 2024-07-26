#!/usr/bin/env bash

cmds=(jq tr cut)
for cmd in $cmds; do
    if ! command -v $cmd &>/dev/null; then
        echo "[ERROR] $cmd command not found"
        exit 1
    fi
done


JSON_DATA_FILE="./doh_resolvers_data_20240119.json"
OUTPUT_FILE="./doh_servers.txt"

jq 'map(.uri)' $JSON_DATA_FILE | jq '.[]' | tr -d '"' > $OUTPUT_FILE

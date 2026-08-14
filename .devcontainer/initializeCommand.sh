#!/bin/sh

for vol in opencode_setting claude_setting codex_setting vscocde-golang-devcontainer; do
    docker volume inspect "$vol" > /dev/null 2>&1 || docker volume create "$vol"
done
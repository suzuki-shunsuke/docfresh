#!/usr/bin/env bash

set -euo pipefail

ls *.md | grep -v README.md | sort -u | while read -r file; do
    title=$(head -n 1 "$file" | sed -E "s/^# //")
    echo "- [${title}](${file})"
done

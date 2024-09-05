#!/bin/bash
cd "$(dirname "$1")"
set -e
FOLDER="builds"
for file in cmd/*; do
    if [ -d "$file" ]; then
        go build -o $PWD/$FOLDER/$(basename $file) $PWD/$file
    fi 
done

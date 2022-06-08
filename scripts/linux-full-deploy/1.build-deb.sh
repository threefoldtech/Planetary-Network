#!/bin/sh
FILE=LICENSE

if [ ! -f "$FILE" ]; then
    echo "Please run script from main directory"
    exit 1
fi

scripts/linux-build.sh
scripts/linux-installer.sh

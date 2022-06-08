#!/bin/bash
FILE=LICENSE

if [ ! -f "$FILE" ]; then
    echo "Please run script from main directory"
    exit 1
fi

./scripts/windows-build.sh
mkdir -p generated/builds/windows/
cp "src/deploy/windows/src.exe" "generated/builds/windows/ThreeFoldPlanetaryNetwork.exe"
cp libs/wintun.dll generated/builds/windows/wintun.dll
cp resources/icon.ico generated/builds/windows/icon.ico

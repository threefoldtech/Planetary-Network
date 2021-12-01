#!/bin/bash
cd ..
./windows-build.sh
mkdir -p generated/builds/windows/
cp "src/deploy/windows/src.exe" "generated/builds/windows/ThreeFoldNetworkConnector.exe"
cp libs/wintun.dll generated/builds/windows/wintun.dll
cp resources/icon.ico generated/builds/windows/icon.ico
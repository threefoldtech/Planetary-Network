#!/bin/bash
echo "ThreeFold Network Connector build script"

echo "Cleaning old build artifacts"
rm -rf src/deploy/darwin & rm -rf "ThreeFold Network Connector.dmg"

echo "Building the binary"
qtdeploy build darwin src/.

echo "Renaming and adding reference to /Application directory path."
mv "src/deploy/darwin/src.app" "src/deploy/darwin/ThreeFold Network Connector.app"
ln -s /Applications src/deploy/darwin/Applications

echo "Generation executable dmg file."
hdiutil create "ThreeFold Network Connector" -srcfolder src/deploy/darwin/

echo "Completed."
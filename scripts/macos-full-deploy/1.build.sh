#!/bin/bash
echo "ThreeFold Planetary Network build script"

echo "Cleaning old build artifacts"
rm -rf src/deploy/darwin & rm -rf "ThreeFold Planetary Network.dmg"

echo "Building the binary"
GO111MODULE=off 
qtmoc desktop src/.
qtdeploy build darwin src/.

echo "Renaming and adding reference to /Application directory path."
mv "src/deploy/darwin/src.app" "src/deploy/darwin/ThreeFold Planetary Network.app"

echo "Adding icns to .app"
cp "resources/src.icns" "src/deploy/darwin/ThreeFold Planetary Network.app/Contents/Resources/src.icns"
ln -s /Applications src/deploy/darwin/Applications

echo "Generation executable dmg file."
hdiutil create "ThreeFold Planetary Network" -srcfolder src/deploy/darwin/

echo "Completed."
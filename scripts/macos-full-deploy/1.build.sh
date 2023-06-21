#!/bin/bash
echo "ThreeFold Planetary Network build script"

echo "Cleaning old build artifacts"
rm -rf src/deploy/darwin & rm -rf "ThreeFold Planetary Network.dmg"

echo "Building the binary"
 
cd src

qtmoc desktop ./
qtdeploy build darwin ./

cd ..

echo overwriting the application bundle information
cp "resources/Info.plist" "src/deploy/darwin/src.app/Contents/Info.plist"

echo "Renaming"
mv "src/deploy/darwin/src.app" "src/deploy/darwin/ThreeFold Planetary Network.app"

echo "Adding icns to .app"
cp "resources/src.icns" "src/deploy/darwin/ThreeFold Planetary Network.app/Contents/Resources/src.icns"

echo "adding reference to /Application directory path."
ln -s /Applications src/deploy/darwin/Applications

echo "Generation executable dmg file."
hdiutil create "ThreeFold Planetary Network" -srcfolder src/deploy/darwin/

echo "Completed."
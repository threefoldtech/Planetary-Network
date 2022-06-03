# Deploy

## Building MSI for windows.
Due to the complex building process of this project it is currently manual. 

### Step 1 (From a Linux machine or using Windows WSL.)
- ./scripts/windows-full-deploy/1.build-and-prepare-windows-binary.sh

### Step 2: Add requird files
- Copy the .exe file to generated/builds/windows
- cp libs/wintun.dll generated/builds/windows/wintun.dll
- cp resources/icon.ico generated/builds/windows/icon.ico

### Step 3 (From cmd or powershell using Windows. Might work with wine, untested.) Note: ResourceHacker must be installed!
- Make sure the path inside the bat script is correct
- ./scripts/windows-full-deploy/2.build-and-prepare-windows-binary-icon.bat

### Step 4 other option. (From cmd or powershell using Windows. Might work with wine, untested.) Note: Can't be done inside WSL must be on a native windows PATH. UNC paths not supported.
- Extract the contents from `https://github.com/wixtoolset/wix3/releases/download/wix3112rtm/wix311-binaries.zip` into a folder `wixbin` and copy it into the root of this project.
- scripts\windows-full-deploy\3.build-and-create-windows-msi-installer.bat


### Step 4 (From cmd or powershell using Windows. Might work with wine, untested.) Note: Can't be done inside WSL must be on a native windows PATH. UNC paths not supported.
- Copy the `wix.xml` file from the msi directory into `generated/builds/windows/`.
- Extract the contents from `https://github.com/wixtoolset/wix3/releases/download/wix3112rtm/wix311-binaries.zip` into a folder `wixbin` and copy it into `generated/builds/windows/`.
- cd inside the `generated/builds/windows` directory and run the following commands:
- wixbin/candle -nologo -out ThreeFoldNetworkConnector-0.0.0.1-x64.wixobj -arch x64 wix.xml
- wixbin/light -nologo -spdb -sice:ICE71 -sice:ICE61 -ext WixUtilExtension.dll -out ThreeFoldNetworkConnector-0.0.0.1-x64.msi ThreeFoldNetworkConnector-0.0.0.1-x64.wixobj


# Dev

## TODO







https://jmeubank.github.io/tdm-gcc/articles/2021-05/10.3.0-release



Install QT 5.13.0
Install Go 1.18
Install MinGW
Install GCC

QT_DIR=C:\Qt\Qt5.13.0
QT_QMAKE_DIR=C:\Qt\Qt5.13.0\5.13.0\mingw73_64\bin


qtdeploy build desktop ./ (in src folder)
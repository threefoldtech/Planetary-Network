# TF-NetworkConnector
Threefold Network Connector is a desktop client (Windows, Mac OS X, Linux) to connect to the ThreeFold Planetary Network, or Yggdrasil Network. It's a GUI client to connect to the Planetary Network with one click.


## More information
See https://forum.threefold.io/t/how-our-planetary-network-works/1210

## Download 
See releases page: https://github.com/threefoldtech/TF-NetworkConnector/releases

## Building MSI for windows.
Due to the complex building process of this project it is currently manual. 

### Step 1 (From a Linux machine or using Windows WSL.)
- ./scripts/windows-full-deploy/1.build-and-prepare-windows-binary.sh

### Step 2 (From cmd or powershell using Windows. Might work with wine, untested.) Note: ResourceHacker must be installed!
- ./scripts/windows-full-deploy/2.build-and-prepare-windows-binary-icon.bat

### Step 3 (From cmd or powershell using Windows. Might work with wine, untested.) Note: Can't be done inside WSL must be on a native windows PATH. UNC paths not supported.
- Copy the `wix.xml` file from the msi directory into `generated/builds/windows/`.
- Extract the contents from `https://github.com/wixtoolset/wix3/releases/download/wix3112rtm/wix311-binaries.zip` into a folder `wixbin` and copy it into `generated/builds/windows/`.
- cd inside the `generated/builds/windows` directory and run the following commands:
- wixbin/candle -nologo -out ThreeFoldNetworkConnector-0.0.0.1-x64.wixobj -arch x64 wix.xml
- wixbin/light -nologo -spdb -sice:ICE71 -sice:ICE61 -ext WixUtilExtension.dll -out ThreeFoldNetworkConnector-0.0.0.1-x64.msi ThreeFoldNetworkConnector-0.0.0.1-x64.wixobj

### Step 3 other option. (From cmd or powershell using Windows. Might work with wine, untested.) Note: Can't be done inside WSL must be on a native windows PATH. UNC paths not supported.
- Extract the contents from `https://github.com/wixtoolset/wix3/releases/download/wix3112rtm/wix311-binaries.zip` into a folder `wixbin` and copy it into the root of this project.
- scripts\windows-full-deploy\3.build-and-create-windows-msi-installer.bat
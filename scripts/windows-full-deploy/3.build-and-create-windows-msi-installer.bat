@echo off
echo "Building the MSI ..."

copy msi\wix.xml generated\builds\windows\
xcopy /E /I wixbin generated\builds\windows\wixbin

cd generated/builds/windows/

"wixbin/candle.exe" -nologo -out ThreeFoldNetworkConnector-0.0.0.1-x64.wixobj -arch x64 wix.xml
"wixbin/light.exe" -nologo -spdb -sice:ICE71 -sice:ICE61 -ext WixUtilExtension.dll -out ThreeFoldNetworkConnector-0.0.0.1-x64.msi ThreeFoldNetworkConnector-0.0.0.1-x64.wixobj
echo "Finished"
pause
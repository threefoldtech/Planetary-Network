@echo off
echo "Inserting threefold ico in the executable binary"
"C:\Program Files (x86)\Resource Hacker\ResourceHacker.exe" -open "\\wsl.localhost\Ubuntu-18.04\home\singlecore\documents\jimber\new\TF-NetworkConnector\generated/builds/windows/ThreeFoldPlanetaryNetwork.exe" -save "\\wsl.localhost\Ubuntu-18.04\home\singlecore\documents\jimber\new\TF-NetworkConnector\generated/builds/windows/ThreeFoldPlanetaryNetwork.exe" -action addskip -res "\\wsl.localhost\Ubuntu-18.04\home\singlecore\documents\jimber\new\TF-NetworkConnector\generated/builds/windows/icon.ico" -mask ICONGROUP,MAIN

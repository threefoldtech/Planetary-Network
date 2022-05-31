@echo off
echo "Inserting threefold ico in the executable binary"
"C:\Program Files (x86)\Resource Hacker\ResourceHacker.exe" -open "C:\Users\jimber\Documents\Threefold\planetary_network\generated\builds\windows\ThreeFoldNetworkConnector.exe" -save "C:\Users\jimber\Documents\Threefold\planetary_network\generated\builds\windows\ThreeFoldNetworkConnector.exe" -action addskip -res "C:\Users\jimber\Documents\Threefold\planetary_network\generated\builds\windows\icon.ico" -mask ICONGROUP,MAIN

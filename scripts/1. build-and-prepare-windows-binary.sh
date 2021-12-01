#!/bin/bash
cd ..

docker build -t rcpwin - < scripts/Dockerfile.windows
docker rm -fv yggwin || true
docker run -d --name=yggwin -v $(pwd)/src:/src rcpwin tail -f /dev/null
docker exec yggwin /bin/buildwindows #run this to build again without restarting docker
mkdir -p generated/builds/windows/
cp "src/deploy/windows/src.exe" "generated/builds/windows/ThreeFoldNetworkConnector.exe"
cp libs/wintun.dll generated/builds/windows/wintun.dll
cp resources/icon.ico generated/builds/windows/icon.ico
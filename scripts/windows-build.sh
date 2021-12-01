#!/bin/bash
FILE=LICENSE

if [ ! -f "$FILE" ]; then
    echo "Please run script from main directory"
    exit 1
fi
sudo chown -R $USER src
cp src/go.mod.back src/go.mod


docker build -t rcpwin - < Dockerfile.windows
docker rm -fv yggwin || true
docker run -d --name=yggwin -v $(pwd)/src:/src rcpwin tail -f /dev/null
docker exec yggwin /bin/buildwindows #run this to build again without restarting docker

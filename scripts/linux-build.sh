#!/bin/bash
FILE=LICENSE

if [ ! -f "$FILE" ]; then
    echo "Please run script from main directory"
    exit 1
fi
sudo chown -R $USER src
cp src/go.mod.back src/go.mod 

docker build -t rcplinux - < Dockerfile.linux
docker rm -fv ygglinux || true
docker run -d --name=ygglinux -v $(pwd)/src:/src rcplinux tail -f /dev/null
docker exec ygglinux /bin/buildlinux  # Run this to build again without restarting docker

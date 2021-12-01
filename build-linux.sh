#!/bin/bash
docker build -t rcplinux - < Dockerfile.linux
docker rm -fv ygglinux || true
docker run -d --name=ygglinux -v $(pwd)/src:/src rcplinux tail -f /dev/null
docker exec ygglinux /bin/buildlinux #run this to build again without restarting docker

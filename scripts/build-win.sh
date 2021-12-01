#!/bin/bash
docker build -t rcpwin - < Dockerfile.windows
docker rm -fv yggwin || true
docker run -d --name=yggwin -v $(pwd)/src:/src rcpwin tail -f /dev/null
docker exec yggwin /bin/buildwindows #run this to build again without restarting docker

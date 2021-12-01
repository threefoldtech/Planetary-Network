#!/bin/bash
docker build --build-arg QT_CI_LOGIN=jonas.delrue@jimber.org --build-arg QT_CI_PASSWORD=RunLikeH3 -t qttest .
docker rm -fv ygg
docker run -d --name=ygg qttest tail -f /dev/null
docker cp ygg:/root/go/src/github.com/threefoldtech/yggdrasil-desktop-client/deploy/linux/yggdrasil-desktop-client .
docker rm -fv ygg

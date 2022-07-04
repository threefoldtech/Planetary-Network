## Dev

#### Build and install for Linux (Creating docker and make .deb file) 
`./scripts/linux-build.sh`

#### Run the UI
`./src/deploy/linux/src`

#### Run the Server
`sudo ./src/deploy/linux/src -server`

#### When you have code changes:
`docker exec ygglinux /bin/buildlinux`


## Deploy

`./scripts/linux-full-deploy/1.build-deb.sh`

** File is located in the root of the project **
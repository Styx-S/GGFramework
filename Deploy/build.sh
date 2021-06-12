#!/usr/bin/env bash

env_os=`go env | grep GOOS | sed 's/\"//g'`
env_arch=`go env | grep GOARCH | sed 's/\"//g'`

echo $env_os
echo $env_arch

# cross compile
go env -w GOOS="linux"
go env -w GOARCH="386"

go build ../*.go
mv ./main app/
# clear old
docker rmi -f ggframework:v1
# build
docker build -t ggframework:v1 .
docker images

if [ $1 ]; then
  echo "enable export docker and transport..."
  docker save -o ggf.tar ggframework:v1
  scp ./ggf.tar root@${1}:~/docker/
  echo "auto delete"
  rm -f ./ggf.tar
fi

# reset
go env -w ${env_arch}
go env -w ${env_os}

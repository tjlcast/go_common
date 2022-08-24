#!/usr/bin/env bash

osType=$1
appName="xxxx"
packageName="main"

./git-info.sh $packageName

if [ -z $osType ]; then
    osType="linux"
fi

if [ $osType = mac ]; then
    echo "Choose: mac"
    ./build-mac.sh $appName
elif [ $osType = win ]; then
    echo "Choose: win"
    ./build-win.sh $appName
else
    echo "Choose: linux"
    ./build-linux.sh $appName
fi

go env -w CGO_ENABLED=1 GOOS=darwin GOARCH=amd64
#go env -w CGO_ENABLED=1 GOOS=linux GOARCH=amd64
#go env -w CGO_ENABLED=1 GOOS=darwin GOARCH=amd64

echo "Finish build."

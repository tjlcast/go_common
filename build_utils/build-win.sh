#!/usr/bin/env bash

appName=$1

echo "build to win."
go env -w CGO_ENABLED=0 GOOS=windows GOARCH=amd64
go build -o $appName

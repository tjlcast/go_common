#!/usr/bin/env bash

appName=$(./app.sh)

nohup ./$appName &
echo $appName > .app

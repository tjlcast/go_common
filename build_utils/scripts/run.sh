#!/usr/bin/env bash

appName="awesome"

nohup ./$appName &
echo $appName > .app

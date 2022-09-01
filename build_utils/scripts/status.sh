#!/usr/bin/env bash

appName=$(./app.sh)

ps -ef | grep -E "./"$appName

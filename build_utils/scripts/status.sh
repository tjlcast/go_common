#!/usr/bin/env bash

appName=$(cat .app)

ps -ef | grep -E "./"$appName

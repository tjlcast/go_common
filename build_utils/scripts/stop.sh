#!/usr/bin/env bash

pid=$(echo $(.status.sh) | awk '{pid=$2; print pid}')
kill -9 $pid

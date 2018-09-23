#!/usr/bin/env bash

if [ $# != 1 ]; then
    echo "Usage: $0 [Code File Name]"
    exit 0
fi

_BIN_NAME=${1%.*}

rm ./pkg/*
GOOS=linux GOARCH=amd64 gom build -ldflags="-w" -o ./pkg/${_BIN_NAME}_linux_amd64
GOOS=windows GOARCH=amd64 gom build -ldflags="-w" -o ./pkg/${_BIN_NAME}_windows_amd64.exe
GOOS=darwin GOARCH=amd64 gom build -ldflags="-w" -o ./pkg/${_BIN_NAME}_darwin_amd64

#!/usr/bin/env bash

if [ $# != 1 ]; then
    echo "Usage: $0 [Code File Name]"
    exit 0
fi

_BIN_NAME=${1%.*}

rm ./linux/${_BIN_NAME}
GOOS=linux GOARCH=amd64 gom build -ldflags="-w" -o ./linux/${_BIN_NAME}
rm ./win/${_BIN_NAME}.exe
GOOS=windows GOARCH=amd64 gom build -ldflags="-w" -o ./win/${_BIN_NAME}.exe
rm ./osx/${_BIN_NAME}
GOOS=darwin GOARCH=amd64 gom build -ldflags="-w" -o ./osx/${_BIN_NAME}

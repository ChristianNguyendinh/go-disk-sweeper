#!/bin/bash

BASEDIR=$(dirname "$0")
FILES=( "$BASEDIR/sweeper.go" "$BASEDIR/info.go" "$BASEDIR/prompt.go" "$BASEDIR/windows.go" "$BASEDIR/linux.go" )

env GOOS=darwin GOARCH=amd64 go build -o $BASEDIR/dist/sweeper-darwin-amd64 ${FILES[@]}
if [ $? -eq 0 ]
then
    echo "MacOS 64-bit compilation SUCCESS"
else
    echo "MacOS 64-bit compilation FAILED!"
fi

env GOOS=linux GOARCH=386 go build -o $BASEDIR/dist/sweeper-linux-386 ${FILES[@]}
if [ $? -eq 0 ]
then
    echo "Linux 64-bit compilation SUCCESS"
else
    echo "Linux 64-bit compilation FAILED!"
fi

env GOOS=linux GOARCH=amd64 go build -o $BASEDIR/dist/sweeper-linux-amd64 ${FILES[@]}
if [ $? -eq 0 ]
then
    echo "Linux 32-bit compilation SUCCESS"
else
    echo "Linux 32-bit compilation FAILED!"
fi

env GOOS=windows GOARCH=amd64 go build -o $BASEDIR/dist/sweeper-windows-amd64.exe ${FILES[@]}
if [ $? -eq 0 ]
then
    echo "Windows 64-bit compilation SUCCESS"
else
    echo "Windows 64-bit compilation FAILED!"
fi

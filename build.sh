#!/bin/bash -
declare -r Name="ecr-builder"

for GOOS in darwin linux; do
    GO111MODULE=on GOOS=$GOOS GOARCH=amd64 go build -o bin/ecr-builder-$GOOS-amd64 *.go
done

#!/bin/bash

VERSION=$(cat version)
export $(cat env)

/home/ubuntu/.goenv/bin/goenv install $VERSION

/home/ubuntu/.goenv/versions/$VERSION/bin/go build -o build/zd cmd/zd/main.go
./build/zd

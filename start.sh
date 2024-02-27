#!/bin/sh

/home/ubuntu/.goenv/versions/1.21.6/bin/go build -o build/zd cmd/zd/main.go
./build/zd

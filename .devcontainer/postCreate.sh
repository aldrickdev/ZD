#!/bin/sh

# https://www.npmjs.com/package/json-server
echo "Installing JSON Server as a fake database"
npm install -g json-server

# https://github.com/cosmtrek/air/tree/master
echo "Install air for reloading on file changes"
go install github.com/cosmtrek/air@latest

#!/bin/bash

# Export Environment Variables
export $(cat env)

# Get runtime version
VERSION=$(cat version)

# Setup asdf variables
ASDF_BIN_DIR="/home/ubuntu/.asdf/bin"
ASDF_STARTUP_SCRIPT="/home/ubuntu/.asdf/asdf.sh"

# Load asdf
source $ASDF_STARTUP_SCRIPT

# Install the proper runtime version
$ASDF_BIN_DIR/asdf install golang $VERSION
$ASDF_BIN_DIR/asdf local golang $VERSION

# Load changes
source $ASDF_STARTUP_SCRIPT

# Build executable and run it
$ASDF_BIN_DIR/asdf exec go build -o build/zd cmd/zd/main.go
./build/zd

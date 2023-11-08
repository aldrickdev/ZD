# Variables
MOD_NAME=github.com/aldrickdev/gin_api_template

# Basic Commands
run:
	go run main.go

init:
	go mod init $(MOD_NAME)
	go mod tidy

# Devcontainer Commands
dc_up:
	devcontainer up

dc_up_r:
	devcontainer up --remove-existing-container

dc_open:
	devcontainer open

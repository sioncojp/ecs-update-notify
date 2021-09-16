#!/bin/bash -
GOOS=darwin GOARCH=amd64 go build -o bin/ecs-update-notify-darwin-amd64 cmd/ecs-update-notify/*.go
GOOS=darwin GOARCH=arm64 go build -o bin/ecs-update-notify-darwin-arm64 cmd/ecs-update-notify/*.go
GOOS=linux  GOARCH=amd64 go build -o bin/ecs-update-notify-linux-amd64  cmd/ecs-update-notify/*.go

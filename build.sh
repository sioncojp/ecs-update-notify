#!/bin/bash -
for GOOS in darwin linux; do
    GOOS=$GOOS GOARCH=amd64 go build -o bin/ecs-update-notify-$GOOS-amd64 cmd/ecs-update-notify/*.go
done
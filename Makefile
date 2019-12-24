$VERSION := $(shell git describe --tags)
$BUILD := $(shell git rev-parse --short HEAD)
$PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
$GOBASE := $(shell pwd)
$GOFILES := $(wildcard *.go)

# Use linker flags to provide version/build settings
$LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

# Redirect error output to a file, so we can show it in development mode.
$STDERR := /tmp/.$(PROJECTNAME)-stderr.txt

# PID file will keep the process id of the server
$PID := /tmp/.$(PROJECTNAME).pid

# Make is verbose in Linux. Make it silent.
$MAKEFLAGS += --silent


# 避免和同名文件冲突，改善性能 https://stackoverflow.com/questions/2145590/what-is-the-purpose-of-phony-in-a-makefile
.PHONY: test build

dev:
	./bin/air -c ./configs/air.conf

test:
	go test ./...

test-cover:
	mkdir codecoverage && ginkgo -r -cover -outputdir=./codecoverage/ -coverprofile=coverage.txt

build:
	go build -o ./bin/server ./cmd/server/main.go
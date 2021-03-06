# Make is verbose in Linux. Make it silent.
$MAKEFLAGS += --silent


# 避免和同名文件冲突，改善性能 https://stackoverflow.com/questions/2145590/what-is-the-purpose-of-phony-in-a-makefile
.PHONY:test build init

init:
	bash ./scripts/init.sh && go mod download && go mod verify

dev: export EVEREST_CONFIG_FILE_NAME=viper.local
dev:
	air -c ./configs/air.conf

test:
	go test ./...

test-cover:
	mkdir codecoverage && ginkgo -r -cover -outputdir=./codecoverage/ -coverprofile=coverage.txt

build:
	GOOS="linux" go build -o ./bin/server ./cmd/server/main.go

build-image:
	bash ./build/package/build.sh

migrate: export EVEREST_CONFIG_FILE_NAME=viper.local
migrate:
	go run cmd/database/main.go
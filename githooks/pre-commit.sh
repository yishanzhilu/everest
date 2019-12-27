#!/bin/bash

GOLINT=$GOPATH/bin/golint
GOIMPORTS=$GOPATH/bin/goimports

# Check for golint
if [[ ! -x "$GOLINT" ]]; then
  printf "\t\033[41mPlease install golint\033[0m (go get -u golang.org/x/lint/golint)"
  exit 1
fi

# Check for goimports
if [[ ! -x "$GOIMPORTS" ]]; then
  printf "\t\033[41mPlease install goimports\033[0m (go get golang.org/x/tools/cmd/goimports)"
  exit 1
fi

PASS=true
$GOLINT "-set_exit_status" ./...
if [[ $? == 1 ]]; then
  printf "\t\033[31mgolint \033[0m \033[0;30m\033[41mFAILURE!\033[0m\n"
  PASS=false
else
  printf "\t\033[32mgolint \033[0m \033[0;30m\033[42mpass\033[0m\n"
fi

go vet ./...
if [[ $? != 0 ]]; then
  printf "\t\033[31mgo vet \033[0m \033[0;30m\033[41mFAILURE!\033[0m\n"
  PASS=false
else
  printf "\t\033[32mgo vet \033[0m \033[0;30m\033[42mpass\033[0m\n"
fi

go test ./...
if [[ $? != 0 ]]; then
    printf "\t\033[31mgo test \033[0m \033[0;30m\033[41mFAILURE!\033[0m\n"
    PASS=false
else
    printf "\t\033[32mgo test \033[0m \033[0;30m\033[42mpass\033[0m\n"
fi
echo ""

if ! $PASS; then
  printf "\033[0;30m\033[41mCOMMIT FAILED\033[0m\n"
  exit 1
else
  printf "\033[0;30m\033[42mCOMMIT SUCCEEDED\033[0m\n"
fi

exit 0
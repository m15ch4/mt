#!/usr/bin/env bash

PACKAGE="micze.io/mt/cmd/alpha"

mkdir -p ./bin

go build -o bin $PACKAGE

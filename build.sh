#!/bin/sh

basedir=$(cd `dirname $0`; pwd)

export GOPATH=$(pwd)
export GOBIN=$(pwd)/bin

go install src/main.go


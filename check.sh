#!/bin/bash

set -e

_arg="./$1"

echo "Fmt ..."
go fmt $_arg

echo "Vet ..."
go vet $_arg

echo "Test ..."
go test $_arg


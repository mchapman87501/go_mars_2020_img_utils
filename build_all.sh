#!/bin/bash

# Build all of the commands and put them in a sandbox directory.
# Avoids polluting ${GOPATH}/bin, eh.
if [ ! -d sandbox ]; then
    mkdir sandbox
fi
go build -o sandbox ./...


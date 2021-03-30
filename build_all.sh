#!/bin/bash

# Build all of the commands and put them in a sandbox directory.
# Avoids polluting ${GOPATH}/bin, eh.
go build -o sandbox ./...


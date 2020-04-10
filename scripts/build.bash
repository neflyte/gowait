#!/usr/bin/env bash
echo "o  building..."
CGO_ENABLED=0 go build -i -installsuffix nocgo -pkgdir "${GOPATH}/pkg" -ldflags "-s -w -extldflags '-static'" -o bin/gowait ./cmd/gowait

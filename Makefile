# Makefile
#
lint:
	golangci-lint run

build:
	CGO_ENABLED=0 go build -i -installsuffix nocgo -pkgdir "$(shell go env GOPATH)/pkg" -ldflags "-s -w -extldflags '-static'" -o bin/gowait ./cmd/gowait
ifeq ($(shell uname -a),Linux)
	go get github.com/pwaller/goupx
	goupx -q bin/gowait
else
	upx -q bin/gowait
endif

clean:
	rm bin/gowait

build-docker: lint
	docker build --no-cache -t neflyte/gowait:latest .

# Makefile
#
.PHONY: lint build clean outdated ensure-fieldalignment check-fieldalignment autofix-fieldalignment build-docker deploy-test-kafka-helm undeploy-test-kafka-helm

# Check Make version (we need at least GNU Make 3.82). Fortunately,
# 'undefine' directive has been introduced exactly in GNU Make 3.82.
ifeq ($(filter undefine,$(.FEATURES)),)
$(error Unsupported Make version. \
    The build system does not work properly with GNU Make $(MAKE_VERSION), \
    please use GNU Make 3.82 or above.)
endif

# Set platform-specific build variables
ifeq ($(OS),Windows_NT)
SHELL=C:\Windows\system32\cmd.exe
.SHELLFLAGS=/C
endif

build:
ifeq ($(OS),Windows_NT)
	SET CGO_ENABLED=0
	go build -ldflags "-s -w" -o bin\gowait.exe .\cmd\gowait
	upx -q bin\gowait.exe
else
	CGO_ENABLED=0 go build -ldflags "-s -w" -o bin/gowait ./cmd/gowait
	upx -q bin/gowait
endif

lint: check-fieldalignment
	@golangci-lint --version
	golangci-lint run --timeout=10m --verbose

clean:
ifeq ($(OS),Windows_NT)
	IF EXIST bin RD /S /Q bin
else
	if [ -d bin ]; then rm -Rf bin; fi
endif

outdated:
ifeq ($(OS),Windows_NT)
	PUSHD %HOMEDRIVE%%HOMEPATH% && go install github.com/psampaz/go-mod-outdated@v0.9.0
else
	hash go-mod-outdated 2>/dev/null || { cd && go install github.com/psampaz/go-mod-outdated@v0.9.0; }
endif
	go list -json -u -m all | go-mod-outdated -direct -update

ensure-fieldalignment:
ifeq ($(OS),Windows_NT)
	PUSHD %HOMEDRIVE%%HOMEPATH% && go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
else
	hash fieldalignment 2>/dev/null || { cd && go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest; }
endif

check-fieldalignment: ensure-fieldalignment
	fieldalignment ./...

autofix-fieldalignment: ensure-fieldalignment
	fieldalignment -fix ./...

build-docker: lint
	docker build --no-cache -t neflyte/gowait:latest .

deploy-test-kafka-helm: build-docker
	helm upgrade kafka testdata/helm/v2/kafka/gowait-kafka --install --namespace default

undeploy-test-kafka-helm:
	helm delete --purge kafka

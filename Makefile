# Makefile
#
.PHONY: lint build clean build-docker deploy-test-kafka-helm undeploy-test-kafka-helm
build:
	CGO_ENABLED=0 go build -ldflags "-s -w" -o bin/gowait ./cmd/gowait
	upx -q bin/gowait

lint:
	go vet ./...
	golangci-lint run

clean:
	rm bin/gowait

build-docker: lint
	docker build --no-cache -t neflyte/gowait:latest .

deploy-test-kafka-helm: build-docker
	helm upgrade kafka testdata/helm/v2/kafka/gowait-kafka --install --namespace default

undeploy-test-kafka-helm:
	helm delete --purge kafka

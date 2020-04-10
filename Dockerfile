FROM golang:1.12-buster AS builder
RUN apt-get update --yes && apt-get install --yes upx-ucl
COPY . /src/gowait
WORKDIR /src/gowait
RUN go mod download
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w -extldflags '-static'" -o /bin/gowait ./cmd/gowait
RUN go get github.com/pwaller/goupx
RUN goupx -q /bin/gowait

FROM scratch
WORKDIR /usr/local/bin
COPY --from=builder /bin/gowait /usr/local/bin/gowait
ENTRYPOINT [ "/usr/local/bin/gowait" ]

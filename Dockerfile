FROM golang:1.16-buster AS builder
RUN apt-get update --yes && apt-get upgrade --yes && apt-get install --yes upx-ucl
COPY . /src/gowait
WORKDIR /src/gowait
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w" -o /bin/gowait ./cmd/gowait
RUN upx -q /bin/gowait

FROM scratch
WORKDIR /usr/local/bin
COPY --from=builder /bin/gowait /usr/local/bin/gowait
ENTRYPOINT [ "/usr/local/bin/gowait" ]

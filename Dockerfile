FROM golang:1.18 AS builder
RUN apt-get update --yes && apt-get install --yes upx-ucl
COPY . /src/gowait
WORKDIR /src/gowait
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w" -o /bin/gowait ./cmd/gowait
RUN upx -q /bin/gowait

FROM gcr.io/distroless/static
WORKDIR /usr/local/bin
COPY --from=builder /bin/gowait /usr/local/bin/gowait
ENTRYPOINT [ "/usr/local/bin/gowait" ]

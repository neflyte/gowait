@ECHO OFF
ECHO o  building...
set CGO_ENABLED=0
go build -i -installsuffix nocgo -pkgdir "C:/Users/alan/go/pkg" -ldflags "-s -w -extldflags '-static'" -o bin/gowait.exe ./cmd/gowait

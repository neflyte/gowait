#!/usr/bin/env bash
echo "o  vetting..."
go vet ./... || {
  echo "*  error running 'go vet'; aborting"
  exit 1
}
type -p golangci-lint &>/dev/null || {
  echo "o  golangci-lint not found; attempting to install it"
  OSNAME=$(uname -s)
  case OSNAME in
    Linux)
      curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(go env GOPATH)/bin" v1.24.0
      [ $? -eq 0 ] || {
        echo "*  error installing golangci-lint; aborting"
        exit 1
      }
      ;;
    Darwin)
      brew install golangci/tap/golangci-lint
      [ $? -eq 0 ] || {
        echo "*  error installing golangci-lint; aborting"
        exit 1
      }
      ;;
    *)
      echo "*  unknown OS: ${OSNAME}; aborting"
      exit 1
      ;;
  esac
}
echo "o  linting..."
golangci-lint run || {
  echo "*  error running 'golangci-lint'; aborting"
  exit 1
}
echo "done."

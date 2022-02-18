# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.4] - TBA
### Added
- Add `go vet` to lint commands

### Changed
- Upgrade `Shopify/sarama` from v1.26.4 to v1.31.1
- Upgrade `lib/pq` from v1.7.0 to v1.10.4
- Upgrade `neflyte/configmap` from 20200412 to v0.2.0
- Upgrade `sirupsen/logrus` from v1.6.0 to v1.8.1
- Upgrade `gopkg.in/yaml.v3` from 20200615 to 20210107
- Update minimum Golang version to v1.16
- Update Dockerfile to use gcr.io/distroless/static as the base image so CA certificates are available
- Update Dockerfile to use golang-bullseye image for building
- No longer display warnings if some configuration values are empty
- Refactor log format constants into the logger package
- Refactor logger configuration into the logger package
- Exit immediately if no URL was specified to wait for

### Removed
- Remove `goupx` in favour of `upx`
- Stop installing built artifacts during build (`-i -installsuffix cgo` parameters to `go build`)

## [0.1.3] - 2020-04-10
### Added
- HTTP waiter support
- Basic integration test script

### Changed
- Move from `*logrus.Entry` to `logrus.FieldLogger`
- Ensure secret strings are not empty before using them
- Ensure URL credential exists before cloning it for display
- Ensure YAML and JSON configuration files load correctly

## [0.1.2] - 2019-11-??
### Added
- ...

### Changed
- ...

## [0.1.1] - 2019-11-??
### Added
- ...

### Changed
- ...

## [0.1.0] - 2019-11-??
### Added
- TBD...

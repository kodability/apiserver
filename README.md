# API Server
[![Build Status](https://travis-ci.org/kodability/apiserver.svg?branch=develop)](https://travis-ci.org/kodability/apiserver)
[![Coverage Status](https://coveralls.io/repos/github/kodability/apiserver/badge.svg?branch=develop)](https://coveralls.io/github/kodability/apiserver?branch=develop)

## Build
### Without Docker
Clone project by running:
```bash
$ go get github.com/kodability/apiserver
```

Build an executable by running:
```bash
$ go build -v
# or
$ make build
```

### With Docker
Create a docker image before build.
```bash
$ make docker-image
```

Build linux-x64 executable using docker:
```bash
$ make docker-build
```

## Run
### Without Docker
```bash
$ go run main.go
# or
$ make run
```

### With Docker
Run application inside docker container:
```bash
$ make docker-run
```


## Development
### Visual Studio Code
By default, vscode settings are added in `.vscode` directory.

[gometalinter](https://github.com/alecthomas/gometalinter) is used for lint.

```bash
$ go get -u gopkg.in/alecthomas/gometalinter.v2
```

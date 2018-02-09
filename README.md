# tryout-runner

## Build
### Without Docker
Clone project into `$GOPATH/src/github.com/kodability/tryout-runner` directory.

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

Build system specific executables using docker:
```bash
# compile x64-linux
$ make docker-linux64
# compile x64-windows
$ make docker-win64
# compile x64-osx
$ make docker-osx
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

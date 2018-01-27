# tryout-runner

## Build
### Without Docker
Clone project into `$GOPATH/src/tryout-runner` directory.

Build an executable by running:
```bash
$ go build -v
```
or
```bash
$ make compile
```

### With Docker
Create a docker image before build.
```bash
$ make image
```

Run application inside docker container:
```bash
$ make run
```

Build system specific executables using docker:
```bash
# compile x64-linux
$ make compile-linux64
# compile x64-windows
$ make compile-win64
# compile x64-osx
$ make compile-osx
```


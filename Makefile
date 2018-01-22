BASE_DIR=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))

GOPATH=$(BASE_DIR)
BINARY_NAME=tryout-runner


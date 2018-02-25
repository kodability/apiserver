BASE_DIR=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))

NAME=tryout-runner
IMAGE_NAME=$(NAME)
BINARY_NAME=$(NAME)
WORK_DIR=/go/src/github.com/kodability/$(NAME)
COVER_FILE=$(BASE_DIR)/coverage.out
COVER_HTML=$(BASE_DIR)/coverage.html

BUILD_NDEBUG=-ldflags="-s -w"

docker-image:
	@docker build -t $(IMAGE_NAME) .
docker-sh:
	@docker run -it --rm -v "$(BASE_DIR)":$(WORK_DIR) -w $(WORK_DIR) $(IMAGE_NAME) /bin/bash
docker-build:
	@docker run --rm -v "$(BASE_DIR)":$(WORK_DIR) -w $(WORK_DIR) -e GOOS=linux -e GOARCH=amd64 golang:1.10 go build -v $(BUILD_NDEBUG)
docker-run:
	@docker run -it --rm -p 8080:8080 $(IMAGE_NAME) $(BINARY_IMAGE)

build:
	@go build -i -v
build-min:
	@go build -i -v $(BUILD_NDEBUG) && upx --brute $(BINARY_NAME)

run: build
	@./$(BINARY_NAME)
prodrun: build
	@BEEGO_RUNMODE=prod ./$(BINARY_NAME)

test:
	@cd tests && go test -v
cover:
	@cd tests && go test -covermode=count -coverprofile=$(COVER_FILE) -coverpkg=../...
cover-html: cover
	@go tool cover -html=$(COVER_FILE) -o $(COVER_HTML)
cover-std: cover
	@go tool cover -func=$(COVER_FILE)

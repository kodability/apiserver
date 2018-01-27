BASE_DIR=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))

NAME=tryout-runner
IMAGE_NAME=$(NAME)
BINARY_NAME=$(NAME)
WORK_DIR=/go/src/$(NAME)

docker-image:
	@docker build -t $(IMAGE_NAME) .
docker-sh:
	@docker run -it --rm $(IMAGE_NAME) /bin/bash
docker-linux64:
	@docker run --rm -v "$(BASE_DIR)":$(WORK_DIR) -w $(WORK_DIR) -e GOOS=linux -e GOARCH=amd64 golang:1.9 go build -v
docker-win64:
	@docker run --rm -v "$(BASE_DIR)":$(WORK_DIR) -w $(WORK_DIR) -e GOOS=windows -e GOARCH=amd64 golang:1.9 go build -v
docker-osx:
	@docker run --rm -v "$(BASE_DIR)":$(WORK_DIR) -w $(WORK_DIR) -e GOOS=darwin -e GOARCH=amd64 golang:1.9 go build -v
docker-run:
	@docker run -it --rm -p 8080:8080 $(IMAGE_NAME) $(BINARY_IMAGE)

build:
	@go build -v
run:
	@go run main.go

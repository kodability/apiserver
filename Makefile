BASE_DIR=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))

NAME=tryout-runner
IMAGE_NAME=$(NAME)
BINARY_NAME=$(NAME)
WORK_DIR=/go/src/$(NAME)
COVER_FILE=$(BASE_DIR)/coverage.out
COVER_HTML=$(BASE_DIR)/coverage.html

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
	@go build -i -v
run:
	@go build -i -v && ./$(BINARY_NAME)
prodrun:
	@go build -i -v && BEEGO_RUNMODE=prod ./$(BINARY_NAME)
test:
	@cd tests && go test
cover:
	@cd tests && go test -cover -coverprofile=$(COVER_FILE) -coverpkg=../...
cover-html: cover
	@go tool cover -html=$(COVER_FILE) -o $(COVER_HTML)
cover-std: cover
	@go tool cover -func=$(COVER_FILE)

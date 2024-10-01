APP_NAME = httpServer
DOCKER_IMAGE = my_http_server
DOCKER_CONTAINER = http_server_container
GO_FILES = $(wildcard *.go) $(wildcard http/*.go) $(wildcard storage/*.go)
TEST_FILE = tests.py


build:
	swag init
	go build -o $(APP_NAME) main.go


run: build
	./$(APP_NAME) -addr=:8080


docker-build:
	docker build -t $(DOCKER_IMAGE) .


docker-run: docker-build
	docker run --name $(DOCKER_CONTAINER) -p 8080:8080 $(DOCKER_IMAGE)


clean:
	rm -f $(APP_NAME)


docker-clean:
	docker rm -f $(DOCKER_CONTAINER)


docker-rm-image:
	docker rmi $(DOCKER_IMAGE)


test:
	python3 $(TEST_FILE)


.PHONY: all build run clean docker-build docker-run docker-clean docker-rm-image test
all: build

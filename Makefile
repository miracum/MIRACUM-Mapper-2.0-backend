.PHONY: all build clean

BINARY_NAME := ./miracummapper
MAIN_SOURCE := ./cmd/miracummapper
BUILD_FLAGS := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
BUILD_COMMAND := go build -ldflags="-w -s" -o $(BINARY_NAME) $(MAIN_SOURCE)
GENERATE_COMMAND := cd tools && go generate -tags tools && cd -

all: build

build:
	$(GENERATE_COMMAND)
	$(BUILD_FLAGS) $(BUILD_COMMAND)

clean:
	rm -f $(BINARY_NAME)

run: build
	$(BINARY_NAME)

run-docker: build
	docker compose up --build -d
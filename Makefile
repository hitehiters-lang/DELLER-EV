.PHONY: build run help

APP_NAME   := deller
BUILD_DIR  := build

help:
	@echo "Targets:"
	@echo "  build           - build app into $(BUILD_DIR)/"
	@echo "  run             - run the app (requires PostgreSQL)"

build:
	cd cmd/server && go build -o ../../$(BUILD_DIR)/$(APP_NAME).exe .
	cd cmd/server && set GOOS=linux&& set GOARCH=amd64&& go build -o ../../$(BUILD_DIR)/$(APP_NAME) .

run:
	swag init -g cmd/server/main.go --parseInternal
	go run cmd/server/main.go
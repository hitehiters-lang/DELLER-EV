.PHONY: build run help copy-images

APP_NAME   := deller
BUILD_DIR  := build
IMAGES_DIR := shared

help:
	@echo "Targets:"
	@echo "  build           - build app into $(BUILD_DIR)/"
	@echo "  run             - run the app (requires PostgreSQL)"
	@echo "  copy-images     - copy images folder to $(BUILD_DIR)/"

copy-images:
	@echo Copying images to $(BUILD_DIR)...
	xcopy /E /I /Y /Q $(IMAGES_DIR) $(BUILD_DIR)\$(IMAGES_DIR)

build: copy-images
	cd cmd/server && go build -o ../../$(BUILD_DIR)/$(APP_NAME).exe .
	cd cmd/server && set GOOS=linux&& set GOARCH=amd64&& go build -o ../../$(BUILD_DIR)/$(APP_NAME) .

run:
	swag init -g cmd/server/main.go --parseInternal
	go run cmd/server/main.go
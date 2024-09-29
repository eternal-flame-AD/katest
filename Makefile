APP_NAME = katest

OUT_DIR = out

BUILD_GOOS = linux darwin
BUILD_GOOS_EXE = windows
BUILD_GOARCH = amd64 arm64


matrix:
	@echo "GOOS: $(BUILD_GOOS)"
	@echo "GOARCH: $(BUILD_GOARCH)"
	@echo "OUT_DIR: $(OUT_DIR)"
	@echo "Build all: make all"
	
all:
	@for goarch in $(BUILD_GOARCH); do \
		for goos in $(BUILD_GOOS); do \
			echo "Building $$goos-$$goarch"; \
			GOOS=$$goos GOARCH=$$goarch \
				go build -ldflags "-s -w" \
				  -o $(OUT_DIR)/$(APP_NAME)-$$goos-$$goarch .; \
		done; \
		for goos in $(BUILD_GOOS_EXE); do \
			echo "Building $$goos-$$goarch"; \
			GOOS=$$goos GOARCH=$$goarch \
				go build -ldflags "-s -w" \
				  -o $(OUT_DIR)/$(APP_NAME)-$$goos-$$goarch.exe .; \
		done; \
	done
# Version variable (can be overridden with VERSION=v1.2.3 make release)
VERSION ?= v1.0.0

# Detect OS and architecture
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

# Platform-specific directories and binaries
MACX86_DIR = sokru-$(VERSION)_mac_x86
MACM1_DIR = sokru-$(VERSION)_mac_m1
LINUX_DIR = sokru-$(VERSION)_linux_amd64
LINUX_ARM_DIR = sokru-$(VERSION)_linux_arm64
WIN_DIR = sokru-$(VERSION)_win_amd64
WIN_ARM_DIR = sokru-$(VERSION)_win_arm64

# Build directories
BUILD_DIR = build
TEMP_DIR = $(BUILD_DIR)/temp

# Default target based on OS and architecture
ifeq ($(UNAME_S),Darwin)
    ifeq ($(UNAME_M),arm64)
        all: mac
    else
        all: macx86
    endif
else ifeq ($(UNAME_S),Linux)
    all: lin
else
    all: win
endif

macx86:
	GOOS=darwin GOARCH=amd64 go build -o sok main.go
	GOOS=darwin GOARCH=amd64 go build -o sokru main.go

mac:
	GOOS=darwin GOARCH=arm64 go build -o sok main.go
	GOOS=darwin GOARCH=arm64 go build -o sokru main.go

lin:
	GOOS=linux GOARCH=amd64 go build -o sok main.go
	GOOS=linux GOARCH=amd64 go build -o sokru main.go

linarm:
	GOOS=linux GOARCH=arm64 go build -o sok-arm main.go
	GOOS=linux GOARCH=arm64 go build -o sokru-arm main.go

win:
	GOOS=windows GOARCH=amd64 go build -o sok.exe main.go
	GOOS=windows GOARCH=amd64 go build -o sokru.exe main.go

winarm:
	GOOS=windows GOARCH=arm64 go build -o sok-arm.exe main.go
	GOOS=windows GOARCH=arm64 go build -o sokru-arm.exe main.go

release: clean-release
	@echo "Building Sokru $(VERSION) for all platforms..."
	@mkdir -p $(TEMP_DIR)

	# macOS x86_64
	@echo "Building macOS x86_64..."
	@mkdir -p $(TEMP_DIR)/$(MACX86_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(TEMP_DIR)/$(MACX86_DIR)/sok main.go
	GOOS=darwin GOARCH=amd64 go build -o $(TEMP_DIR)/$(MACX86_DIR)/sokru main.go
	cd $(TEMP_DIR) && zip -r ../$(MACX86_DIR).zip $(MACX86_DIR)

	# macOS ARM64 (M1/M2)
	@echo "Building macOS ARM64..."
	@mkdir -p $(TEMP_DIR)/$(MACM1_DIR)
	GOOS=darwin GOARCH=arm64 go build -o $(TEMP_DIR)/$(MACM1_DIR)/sok main.go
	GOOS=darwin GOARCH=arm64 go build -o $(TEMP_DIR)/$(MACM1_DIR)/sokru main.go
	cd $(TEMP_DIR) && zip -r ../$(MACM1_DIR).zip $(MACM1_DIR)

	# Linux AMD64
	@echo "Building Linux AMD64..."
	@mkdir -p $(TEMP_DIR)/$(LINUX_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(TEMP_DIR)/$(LINUX_DIR)/sok main.go
	GOOS=linux GOARCH=amd64 go build -o $(TEMP_DIR)/$(LINUX_DIR)/sokru main.go
	cd $(TEMP_DIR) && zip -r ../$(LINUX_DIR).zip $(LINUX_DIR)

	# Linux ARM64
	@echo "Building Linux ARM64..."
	@mkdir -p $(TEMP_DIR)/$(LINUX_ARM_DIR)
	GOOS=linux GOARCH=arm64 go build -o $(TEMP_DIR)/$(LINUX_ARM_DIR)/sok main.go
	GOOS=linux GOARCH=arm64 go build -o $(TEMP_DIR)/$(LINUX_ARM_DIR)/sokru main.go
	cd $(TEMP_DIR) && zip -r ../$(LINUX_ARM_DIR).zip $(LINUX_ARM_DIR)

	# Windows AMD64
	@echo "Building Windows AMD64..."
	@mkdir -p $(TEMP_DIR)/$(WIN_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(TEMP_DIR)/$(WIN_DIR)/sok.exe main.go
	GOOS=windows GOARCH=amd64 go build -o $(TEMP_DIR)/$(WIN_DIR)/sokru.exe main.go
	cd $(TEMP_DIR) && zip -r ../$(WIN_DIR).zip $(WIN_DIR)

	# Windows ARM64
	@echo "Building Windows ARM64..."
	@mkdir -p $(TEMP_DIR)/$(WIN_ARM_DIR)
	GOOS=windows GOARCH=arm64 go build -o $(TEMP_DIR)/$(WIN_ARM_DIR)/sok.exe main.go
	GOOS=windows GOARCH=arm64 go build -o $(TEMP_DIR)/$(WIN_ARM_DIR)/sokru.exe main.go
	cd $(TEMP_DIR) && zip -r ../$(WIN_ARM_DIR).zip $(WIN_ARM_DIR)

	@echo "Release build complete. Files created:"
	@ls -la $(BUILD_DIR)/*.zip
	@rm -rf $(TEMP_DIR)

clean:
	@rm -f sok* sokru*

clean-release:
	@rm -rf $(BUILD_DIR)

.PHONY: all clean clean-release mac macx86 lin linarm win winarm release

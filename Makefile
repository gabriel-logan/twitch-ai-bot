.PHONY: run build build_linux build_linux_32 build_windows build_windows_32 build_darwin clean test test_race

VERSION ?= dev
API_DIST=cmd/api
APP_NAME=twitch_ai_bot

BUILD_DIR=bin
OUTPUT_DIR=$(APP_NAME)-$(VERSION)

TEMPLATES=templates
PUBLIC=public
SYSTEM_PROMPT=system_prompt.txt
SETUP_FILE=setup.txt

run:
	go run $(API_DIST)/main.go

build: clean build_linux build_linux_32 build_windows build_windows_32 build_darwin

build_linux:
	$(call build_binary,linux,amd64,)

build_linux_32:
	$(call build_binary,linux,386,)

build_windows:
	$(call build_binary,windows,amd64,.exe)

build_windows_32:
	$(call build_binary,windows,386,.exe)

build_darwin:
	$(call build_binary,darwin,amd64,)

define build_binary
	mkdir -p $(BUILD_DIR)/$(1)/$(OUTPUT_DIR)

	GOOS=$(1) GOARCH=$(2) go build -o $(BUILD_DIR)/$(1)/$(OUTPUT_DIR)/$(APP_NAME)-$(VERSION)-$(2)$(3) $(API_DIST)/main.go

	cp .env $(BUILD_DIR)/$(1)/$(OUTPUT_DIR)/
	cp -r $(TEMPLATES) $(BUILD_DIR)/$(1)/$(OUTPUT_DIR)/
	cp $(SYSTEM_PROMPT) $(BUILD_DIR)/$(1)/$(OUTPUT_DIR)/${SYSTEM_PROMPT}

	cp $(SETUP_FILE) $(BUILD_DIR)/$(1)/$(OUTPUT_DIR)/${SETUP_FILE}
	mv $(BUILD_DIR)/$(1)/$(OUTPUT_DIR)/${SETUP_FILE} $(BUILD_DIR)/$(1)/$(OUTPUT_DIR)/setup-$(VERSION).txt

	cp -r $(PUBLIC) $(BUILD_DIR)/$(1)/$(OUTPUT_DIR)/
endef

clean:
	rm -rf bin

test:
	go test ./...

test_race: 
	go test -race ./...

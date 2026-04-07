.PHONY: run build_linux build_windows build_darwin build clean

API_DIST=cmd/api
BINARY_NAME=twitch_ai_bot

run:
	go run $(API_DIST)/main.go

build_linux:
	GOOS=linux GOARCH=amd64 go build -o bin/linux/$(BINARY_NAME) $(API_DIST)/main.go

build_windows:
	GOOS=windows GOARCH=amd64 go build -o bin/windows/$(BINARY_NAME).exe $(API_DIST)/main.go

build_darwin:
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/$(BINARY_NAME) $(API_DIST)/main.go

build: build_linux build_windows build_darwin
	
clean:
	rm -rf bin

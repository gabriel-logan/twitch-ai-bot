.PHONY: run build_linux build_windows build_darwin build clean test test_race

API_DIST=cmd/api
BINARY_NAME=twitch_ai_bot

run:
	go run $(API_DIST)/main.go

build_linux:
	GOOS=linux GOARCH=amd64 go build -o bin/linux/$(BINARY_NAME) $(API_DIST)/main.go
	cp .env bin/linux/.env
	cp -r templates bin/linux
	cp system_prompt.txt bin/linux/system_prompt.txt

build_windows:
	GOOS=windows GOARCH=amd64 go build -o bin/windows/$(BINARY_NAME).exe $(API_DIST)/main.go
	cp .env bin/windows/.env
	cp -r templates bin/windows
	cp system_prompt.txt bin/windows/system_prompt.txt

build_darwin:
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/$(BINARY_NAME) $(API_DIST)/main.go
	cp .env bin/darwin/.env
	cp -r templates bin/darwin
	cp system_prompt.txt bin/darwin/system_prompt.txt

build: build_linux build_windows build_darwin
	
clean: rm -rf bin

test: go test ./...

test_race: go test -race ./...

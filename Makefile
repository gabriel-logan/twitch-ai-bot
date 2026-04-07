.PHONY: run build clean

API_DIST=cmd/api
BINARY_NAME=twitch_ai_bot

run:
	go run $(API_DIST)/main.go

build:
	go build -o bin/$(BINARY_NAME) $(API_DIST)/main.go

clean:
	rm -rf bin

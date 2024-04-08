run: build
	@./bin/radarr-list

build:
	@go build -o bin/radarr-list main.go
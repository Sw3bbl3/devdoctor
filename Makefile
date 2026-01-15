.PHONY: build test clean install

# Build the binary
build:
	go build -o devdoctor ./cmd/devdoctor

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o dist/devdoctor-linux-amd64 ./cmd/devdoctor
	GOOS=linux GOARCH=arm64 go build -o dist/devdoctor-linux-arm64 ./cmd/devdoctor
	GOOS=darwin GOARCH=amd64 go build -o dist/devdoctor-darwin-amd64 ./cmd/devdoctor
	GOOS=darwin GOARCH=arm64 go build -o dist/devdoctor-darwin-arm64 ./cmd/devdoctor
	GOOS=windows GOARCH=amd64 go build -o dist/devdoctor-windows-amd64.exe ./cmd/devdoctor

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Clean build artifacts
clean:
	rm -f devdoctor
	rm -rf dist/

# Install to local bin
install:
	go install ./cmd/devdoctor

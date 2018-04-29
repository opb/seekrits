VERSION := $(shell git describe --tags --always --dirty="-dev")
LDFLAGS := -ldflags='-X "main.Version=$(VERSION)"'

clean:
	rm -rf dist/*

dist: clean
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/seekrits-$(VERSION)-darwin-amd64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o dist/seekrits-$(VERSION)-linux-amd64
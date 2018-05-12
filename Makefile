VERSION := $(shell git describe --tags --always --dirty="-dev")
LDFLAGS := -ldflags='-X "main.Version=$(VERSION)"'

release: gh-release dist
	echo $$(GITHUB_TOKEN)
	github-release release \
	--security-token $$(GITHUB_TOKEN) \
	--user opb \
	--repo seekrits \
	--tag $(VERSION) \
	--name $(VERSION

dist: clean
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -v -o dist/seekrits-$(VERSION)-darwin-amd64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -v -o dist/seekrits-$(VERSION)-linux-amd64

clean:
	rm -rf dist/*

gh-release:
	go get -u github.com/aktau/github-release
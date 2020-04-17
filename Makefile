SHELL = bash
PROJECT_ROOT := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
GIT_COMMIT := $(shell git rev-parse HEAD)
GO_PKGS := $(shell go list ./...)

VERSION := 1.0.0
PLATFORMS := darwin linux windows
os = $(word 1, $@)

.PHONY: build
build:
	go build -o mkpw

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	@mkdir -p dist
	@GOOS=$(os) GOARCH=amd64 go build -o dist/$(os)/mkpw github.com/romantomjak/mkpw
	@zip -q -X -j dist/mkpw_$(VERSION)_$(os)_amd64.zip dist/$(os)/mkpw
	@rm -rf dist/$(os)

.PHONY: release
release: windows linux darwin

.PHONY: clean
clean:
	@rm -f "$(PROJECT_ROOT)/mkpw"

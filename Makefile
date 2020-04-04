CGOFLAG ?= CGO_ENABLED=1
BINDIR ?= /usr/local/bin
BUILDIR ?= out

VERSION=1.0.0-alpha.1
BINARY=govm
PLATFORMS=darwin linux windows freebsd
ARCHITECTURES=386 amd64 arm64 ppc64le s390x

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

default: build

all: build install

dependencies:
	go mod vendor

build: dependencies
	$(CGOFLAG) go build $(LDFLAGS) -o $(BUILDIR)/$(BINARY)

build-all: dependencies
	$(foreach GOOS, $(PLATFORMS),\
    	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -v -o $(BUILDIR)/$(BINARY)-$(GOOS)-$(GOARCH))))

install: build
	@echo "\n** Make sure to run 'make install' as root! **\n"
	cp -f $(BUILDIR)/$(BINARY) $(BINDIR)/

clean:
	rm -rf out && rm -rf dist

release-no-publish:
	export VERSION=$(VERSION) && goreleaser release --skip-publish

.PHONY: clean build install build-all all

IMPORT_PATH := github.com/jayapriya90/chatter
GO          ?= go
GOOS        ?= linux
GOARCH      ?= amd64
BINDIR 		:= $(CURDIR)/bin
DISTDIR     := $(CURDIR)/dist 
LDFLAGS     := -w

.PHONY: install setup clean
install:
	GOBIN=$(BINDIR) $(GO) install -ldflags $(LDFLAGS) $(IMPORT_PATH)
	
clean:
	rm -rf $(BINDIR) $(DISTDIR)

setup: clean
	go get -u github.com/golang/dep/cmd/dep
	# go get -u github.com/pressly/goose/cmd/goose
	dep ensure

format:
	go fmt ./...
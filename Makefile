SHELL = /bin/bash

# --- Go(Golang) ------------------------------------------------------------------------------------
GOCMD=go
GOMOD=$(GOCMD) mod
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOTEST=$(GOCMD) test
GOBUILD=$(GOCMD) build
GOGENERATE=$(GOCMD) generate
GOINSTALL=$(GOCMD) install
GOFLAGS := -v 
LDFLAGS := -s -w

ifneq ($(shell go env GOOS),darwin)
LDFLAGS := -extldflags "-static"
endif

GOLANGCILINTCMD=golangci-lint
GOLANGCILINTRUN=$(GOLANGCILINTCMD) run

.PHONY: go-mod-tidy
go-mod-tidy:
	$(GOMOD) tidy

.PHONY: go-mod-update
go-mod-update:
	$(GOGET) -f -t -u ./...
	$(GOGET) -f -u ./...

.PHONY: go-fmt
go-fmt:
	$(GOFMT) ./...

.PHONY: go-lint
go-lint: go-fmt
	$(GOLANGCILINTRUN) $(GOLANGCILINT) ./...

.PHONY: go-test
go-test:
	$(GOTEST) $(GOFLAGS) ./...

.PHONY: go-generate
go-generate:
	$(GOGENERATE) $(GOFLAGS) ./...

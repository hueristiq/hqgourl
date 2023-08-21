
# Go(Golang) Options
GOCMD=go
GOMOD=$(GOCMD) mod
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOTEST=$(GOCMD) test
GOFLAGS := -v 
LDFLAGS := -s -w

# Golangci Options
GOLANGCILINTCMD=golangci-lint
GOLANGCILINTRUN=$(GOLANGCILINTCMD) run

ifneq ($(shell go env GOOS),darwin)
LDFLAGS := -extldflags "-static"
endif

.PHONY: tidy
tidy:
	$(GOMOD) tidy

.PHONY: update-deps
update-deps:
	$(GOGET) -f -t -u ./...
	$(GOGET) -f -u ./...

.PHONY: format
format:
	$(GOFMT) ./...

.PHONY: lint
lint:
	$(GOLANGCILINTRUN) ./...

.PHONY: test
test:
	$(GOTEST) $(GOFLAGS) ./...

.PHONY: generate
generate:
	go run cmd/schemesgen/main.go -p ./schemes/schemes.go
	go run cmd/tldsgen/main.go -p ./tlds/tlds.go
	go run cmd/unicodesgen/main.go -p ./unicodes/unicodes.go

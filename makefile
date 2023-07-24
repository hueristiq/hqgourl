generate:
	@go run cmd/schemesgen/main.go -p ./schemes/schemes.go && \
	go run cmd/tldsgen/main.go -p ./tlds/tlds.go && \
	go run cmd/unicodesgen/main.go -p ./unicodes/unicodes.go

lint:
	@golangci-lint run ./...

test:
	@go test ./...
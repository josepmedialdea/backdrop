BINARY := backdrop
BINDIR := bin
CMD := ./cmd/backdrop

.PHONY: build clean run vet test lint fmt

build:
	CGO_ENABLED=0 go build -o $(BINDIR)/$(BINARY) $(CMD)

clean:
	rm -rf $(BINDIR)

run: build
	./$(BINDIR)/$(BINARY)

vet:
	go vet ./...

test:
	go test ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .

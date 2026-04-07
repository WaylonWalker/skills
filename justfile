set quiet

# Build the skills binary
build:
    mkdir -p bin
    go build -ldflags "-s -w" -o bin/skills .

# Install locally
install:
    go install -ldflags "-s -w" .

# Run all tests
test:
    go test ./... -v -race

# Run tests with coverage
cover:
    go test ./... -coverprofile=coverage.out -race
    go tool cover -func=coverage.out

# Run linter
lint:
    golangci-lint run ./...

# Format code
fmt:
    gofumpt -w .
    goimports -w .

# Run all checks (test + lint + vet)
check: fmt
    go vet ./...
    just lint
    just test

# Clean build artifacts
clean:
    rm -rf bin coverage.out

# Run the binary (pass args after --)
run *args:
    go run . {{args}}

# Show help
help:
    @just --list

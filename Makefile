# Makefile for bibel tests

.PHONY: all build test quick-test smoke-test clean

all: build test

build:
	@echo "Building bibel..."
	@go build -o bibel cmd/bibel.go

test: build
	@echo "=== Running All Tests ==="
	@echo "1. Running Go date progression tests..."
	@go run build/test/test_dates.go
	@echo "2. Running smoke test..."
	@./build/test/smoke_test.sh

quick-test:
	@echo "Running quick test..."
	@cd build/test && ./quick_test.sh

smoke-test: build
	@echo "Running smoke test..."
	@./build/test/smoke_test.sh

date-test:
	@echo "Running date progression tests..."
	@go run build/test/test_dates.go

clean:
	@rm -f bibel
	@echo "Cleaned up"

help:
	@echo "Available targets:"
	@echo "  all        - Build and run all tests"
	@echo "  build      - Build bibel binary"
	@echo "  test       - Run all tests (CI + Go + smoke)"
	@echo "  quick-test - Quick functionality check"
	@echo "  smoke-test - Daily smoke test"
	@echo "  date-test  - Go date progression tests"
	@echo "  clean      - Remove built binary"
	@echo "  help       - Show this help"

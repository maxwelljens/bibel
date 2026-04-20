#!/bin/bash

# Quick version check test
# Returns SUCCESS or FAILED

echo "=== bibel Quick Version Check ==="

# Change to project root
cd "$(dirname "$0")/.." || exit 1

# Check if binary exists and runs
if [ ! -f ./bibel ]; then
  echo "Building bibel..."
  if ! go build -o ./bibel ./cmd/bibel.go; then
    echo "FAILED: bibel build failed"
    exit 1
  fi
fi

# Test basic functionality
if ! ./bibel --help >/dev/null 2>&1; then
  echo "FAILED: bibel --help failed"
  exit 1
fi

# Test reading modes
for mode in evangelion new_testament old_testament bible; do
  timeout 2 ./bibel --plain -r "$mode" >/dev/null 2>&1
  exit_code=$?
  # Exit code 124 is timeout (from TUI mode), 0 is success, others are errors
  if [ $exit_code -ne 0 ] && [ $exit_code -ne 124 ]; then
    echo "FAILED: Mode $mode not working (exit code: $exit_code)"
    exit 1
  fi
done

# Test output modes (skip TUI with --plain)
for flag in --plain --formatted --numbered; do
  ./bibel $flag -r evangelion >/dev/null 2>&1
  if [ $? -ne 0 ]; then
    echo "FAILED: Flag $flag not working"
    exit 1
  fi
done

echo "SUCCESS: All basic functionality working"
exit 0

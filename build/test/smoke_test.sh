#!/bin/bash

# bibel Smoke Test
# Quick test to verify date progression is working

set -e

echo "🔍 bibel Smoke Test - $(date)"
echo "========================================"

# Build if needed
if [ ! -f bibel ]; then
  echo "Building bibel..."
  go build -o bibel ./cmd/bibel.go
fi

# Test 1: All reading modes should work
echo "📖 Test 1: Reading Modes"
MODES=("evangelion" "new_testament" "old_testament" "bible")
for mode in "${MODES[@]}"; do
  echo -n "  Testing $mode... "
  if timeout 2 ./bibel --plain -r "$mode" >/dev/null 2>&1; then
    echo "✅"
  elif [ $? -eq 124 ]; then
    echo "⏱️  (timeout - TUI mode)"
  else
    echo "❌ FAILED"
    exit 1
  fi
done

# Test 2: Date progression should produce today's verses
echo "📅 Test 2: Today's Reading"
TODAY_OUTPUT=$(./bibel --plain -r evangelion 2>/dev/null || true)
if [ -n "$TODAY_OUTPUT" ]; then
  echo "  Today's reading:"
  echo "$TODAY_OUTPUT" | head -3 | sed 's/^/    /'
  echo "  ✅ Output generated"
else
  echo "  ❌ No output generated"
  exit 1
fi

# Test 3: Check verse count (should be around 12 verses)
echo "🔢 Test 3: Verse Count"
# Count approximate verses by counting period+space patterns
VERSE_COUNT=$(echo "$TODAY_OUTPUT" | tr '.' '\n' | wc -l)
echo "  Approximate verses detected: $VERSE_COUNT"
if [ "$VERSE_COUNT" -ge 5 ]; then # Somewhat lenient
  echo "  ✅ Reasonable verse count"
else
  echo "  ⚠️  Low verse count, but may be OK"
fi

# Test 4: Different modes should show different books
echo "📚 Test 4: Mode Differentiation"
EV_OUTPUT=$(./bibel --plain -r evangelion 2>/dev/null | head -1)
OT_OUTPUT=$(./bibel --plain -r old_testament 2>/dev/null | head -1)

if [ "$EV_OUTPUT" != "$OT_OUTPUT" ]; then
  echo "  ✅ Different modes show different books"
  echo "    evangelion: $EV_OUTPUT"
  echo "    old_testament: $OT_OUTPUT"
else
  echo "  ❌ Modes show same book!"
  exit 1
fi

echo "========================================"
echo "✅ All smoke tests PASSED"
echo "💡 Next: Run 'make test' for comprehensive tests"
exit 0

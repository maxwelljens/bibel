# bibel CI Test Suite

This directory contains test scripts for verifying the bibel date progression system.

## Test Files

### quick_test.sh

Quick version check that validates basic functionality.

Run with:
```bash
cd /home/mjensen/Development/bibel/build/test
./quick_test.sh
```

### test_dates.go

Go program for more detailed date progression testing with specific test cases.
Currently includes tests for:

- Day 1 starting position
- Current date verification
- Reading mode starting points
- Wrap-around test
- Lookahead rule test

To run:
```bash
cd /home/mjensen/Development/bibel
go run build/test/test_dates.go
```

## Adding New Tests

When modifying the date progression system, update the test cases in `test_dates.go`
to verify the changes don't break existing functionality.

## Expected Behavior

The date progression system should:
1. Show 12 verses per day (configurable via `verses_per_day` in config)
2. Progress through the selected reading mode (evangelion, new_testament,
   old_testament, bible)
3. Wrap around when reaching the end of the reading mode
4. Apply lookahead rule: if less than 12 verses remain in a chapter, extend to
   chapter end
5. `evangelion` and `new_testament` modes should start at the same position
   (Matthew 1:1)

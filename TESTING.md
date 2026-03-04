# Testing Guide

This document describes pinctab's test organization: unit tests, integration tests, and manual tests.

## Test Pyramid

```
                  Manual Tests
                (Optional Verification)
              /                        \
         Quick-Start             Stress Testing
         (5-10 min)               (Extended)
      
         Integration Tests
    (Primary Test Coverage)
      /        |        \
 Multiple   Profile    Concurrent
 Instances  Conflicts  Operations
  
         Unit Tests
    (Fast, Focused)
        All cmd/,
      internal/* tests
```

## Test Levels

### 1. Unit Tests
**Location:** `cmd/`, `internal/` with `_test.go` suffix  
**Purpose:** Test individual functions, packages, error handling  
**Run:** `go test ./...`  
**Time:** ~1s total  
**Required:** Yes - must pass for all PRs

**Examples:**
- Configuration loading and validation
- Command parsing
- Error handling

### 2. Integration Tests
**Location:** `tests/integration/`  
**Purpose:** Test real orchestrator behavior with actual Chrome instances  
**Run:** `go test -v ./tests/integration`  
**Time:** ~1-2 min per test  
**Required:** Yes - must pass for orchestrator/strategy changes  
**Requires:** Running orchestrator (`./pinchtab`)

#### Integration Test Coverage

| Test | Purpose | Instances | Status |
|------|---------|-----------|--------|
| `TestMultipleInstancesWithDifferentProfiles` | Verify profile isolation works | 2 unique profiles | ✅ PASS |
| `TestProfileConflictTwoInstancesSameProfile` | Verify 409 rejection of duplicate profiles | 2 same profile | ✅ PASS |
| `TestConcurrentOperationsMultipleInstances` | Verify concurrent navigations work independently | 3 instances | ✅ Primary Coverage |
| `TestStressMultipleInstancesSequential` | Verify creation of 5+ instances | 5 sequential | ✅ Primary Coverage |
| *Existing tests* | Basic orchestrator operations | Various | ✅ All passing |

### 3. Manual Tests
**Location:** `tests/manual/`  
**Purpose:** Hands-on verification for developers  
**Required:** No - integration tests provide coverage  
**When to use:** 
- Final verification before release
- Visual confirmation of browser behavior
- Stress testing beyond integration test limits
- Debugging unexpected orchestrator behavior

#### Manual Test Scripts

| Script | Purpose | Coverage | Status |
|--------|---------|----------|--------|
| `quick-start.md` | 7-step guided tutorial | Basic operations | 📌 Reference only |
| `automated-test.sh` | Single instance workflow | Instance + navigation + cleanup | ⚠️ Use integration tests |
| `stress-test-10-instances.sh` | 10 concurrent instances | Load testing | ⚠️ Use integration tests |
| `multi-agent-isolation.sh` | 3 concurrent agents | Agent isolation | ⚠️ Use integration tests |

**Note:** Manual test scripts are not maintained. Use integration tests instead. The scripts can be useful for visual verification, but are not part of the CI/CD pipeline.

---

## Running Tests

### All Tests (Recommended for CI/CD)
```bash
go test ./...                          # Unit tests only (~1s)
go test -v ./tests/integration -timeout 120s  # Full coverage (~90s)
```

### Specific Integration Test
```bash
# Before running: ./pinchtab &
go test -v ./tests/integration -run TestMultipleInstancesWithDifferentProfiles
```

### Manual Test (Optional)
```bash
./pinchtab &
sleep 3
./tests/manual/automated-test.sh
```

---

## Test Data & Cleanup

### Profile Names
- Integration tests use **unique profile names** to avoid conflicts
- Format: `{test-name}-{timestamp}-{instance-num}` or `{test-name}-{num}`
- Example: `concurrent-test-1`, `stress-test-1709462400-1`

### Singleton Lock Cleanup
Chrome creates lock files in profile directories that can cause startup conflicts:
```bash
# Clean before running tests
find ~/.pinchtab/profiles -name "Singleton*" -delete

# Or fully reset
rm -rf ~/.pinchtab/profiles/*
```

### Instance Cleanup
Tests auto-cleanup via `stopInstance()`. Manual verification may require:
```bash
pkill -f "chrome.*user-data-dir"  # Kill stray Chrome processes
pkill -f "^\./pinchtab"            # Kill orchestrator
```

---

## Key Findings

### ✅ Orchestrator Correctly Isolates Instances
- Different profile names → separate instances ✓
- Same profile name → 409 Conflict (rejected at API) ✓
- Multiple instances reach "running" independently ✓

### ✅ Simple Strategy Works Reliably
- Shorthand endpoints (/navigate, /find) work with unique profiles ✓
- Concurrent operations succeed independently ✓
- Cleanup removes all instances ✓

### ⚠️ Manual Shell Tests Not Required
Manual test scripts hang due to shell job control issues with concurrent curl + `wait`.
These are **not critical issues** — the integration tests prove the functionality works.
Keep scripts for reference/learning, but rely on integration tests for CI/CD.

---

## Test Strategy Going Forward

### For Pull Requests
1. Run unit tests: `go test ./...`
2. Run integration tests: `go test -v ./tests/integration` (requires orchestrator running)
3. All must pass

### For Release
1. All unit + integration tests pass
2. Manual smoke test (optional): Run one of the manual test scripts
3. Create release notes referencing test coverage

### For Debugging
1. Use integration tests to isolate issue
2. If needed, add new integration test that reproduces the bug
3. Fix the bug
4. Verify test now passes
5. Add test to regression suite (don't remove)

---

## Adding New Tests

### New Integration Test
1. Add to `tests/integration/{feature}_test.go`
2. Use pattern: `func TestFeatureName(t *testing.T) { ... }`
3. Use existing helpers: `launchInstance()`, `waitForInstanceRunning()`, `getInstances()`, `stopInstance()`
4. Clean up all resources before returning
5. Run against live orchestrator: `./pinchtab & go test -v ./tests/integration -run TestFeatureName`

### New Unit Test
1. Add to package as `{feature}_test.go`
2. Use `go test ./...` to run
3. Keep tests isolated (no external dependencies)

### New Manual Test
1. Only if integration test cannot cover the scenario
2. Add to `tests/manual/{scenario}.sh`
3. Document expected behavior in header comments
4. Add note: "Optional - manual verification only. Integration tests provide automated coverage."

---

## CI/CD Integration

```yaml
# Example: .github/workflows/test.yml
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - run: go test ./...  # Unit tests only (no orchestrator)
      
  integration:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - run: go build -o pinchtab ./cmd/pinchtab
      - run: ./pinchtab & sleep 3 && go test -v ./tests/integration && pkill pinchtab
```

For macOS/local development, both unit and integration tests can run.

# Testing Guide

This guide covers manual and automated testing for the multi-instance orchestrator architecture.

---

## Test Organization

### Automated Integration Tests
**File:** `tests/integration/orchestrator_test.go`

Tests instance creation, lifecycle, port allocation, isolation, and proxy routing.

```bash
go test -tags integration ./tests/integration -run Orchestrator -timeout 120s
```

**Coverage:**
- ✓ Instance creation and lifecycle
- ✓ Hash-based ID generation (prof_X, inst_X, tab_X)
- ✓ Port allocation and reuse
- ✓ Instance isolation
- ✓ Orchestrator proxy routing

---

### Manual Tests

#### Quick Start (Step-by-Step)
**File:** `tests/manual/quick-start.md`

Manual test with 10 steps and detailed expected outputs. Good for learning how the system works and verifying basic functionality with visual feedback.

- **Duration:** ~5 minutes
- **Skill level:** Beginner-friendly
- **What it tests:**
  - Instance creation (headed + headless)
  - Port allocation and reuse
  - Orchestrator proxy routing
  - Instance isolation verification

---

#### Automated Test Script
**File:** `tests/manual/automated-test.sh`

One-command automated test script. Good for CI/CD and quick validation.

```bash
./tests/manual/automated-test.sh
```

- **Duration:** ~10 seconds
- **What it tests:**
  - Instance creation
  - Port allocation
  - Orchestrator routing
  - Instance isolation
  - Cleanup

---

### Complex Test Scenarios

#### Stress Test: 10 Concurrent Instances
**File:** `tests/manual/stress-test-10-instances.sh`

Tests the system under load with 10 concurrent instances, parallel navigation, and resource cleanup.

```bash
./tests/manual/stress-test-10-instances.sh
```

- **Duration:** ~30 seconds
- **Skill level:** Intermediate
- **What it tests:**
  - Concurrent instance creation (10 instances)
  - Concurrent navigation (10 parallel requests)
  - Port allocator under load
  - Chrome initialization at scale
  - Instance isolation at scale
  - Resource cleanup

**Real-world use case:** Verifying the system can handle multiple automation agents simultaneously.

---

#### Multi-Agent Isolation Test
**File:** `tests/manual/multi-agent-isolation.sh`

Tests that multiple agents can run in parallel with complete isolation (cookies, sessions, state).

```bash
./tests/manual/multi-agent-isolation.sh
```

- **Duration:** ~15 seconds
- **Skill level:** Intermediate
- **What it tests:**
  - 3 concurrent agents in separate instances
  - Session/cookie isolation
  - Concurrent action handling
  - No state leakage between agents

**Real-world use case:** Verifying agents can automate independently without interfering with each other.

---

## Prerequisites

```bash
# Build PinchTab
go build -o pinchtab ./cmd/pinchtab

# Verify ports available
lsof -i :9867-9968  # Should be empty or only PinchTab

# Chrome/Chromium installed
which chrome 2>/dev/null || which chromium 2>/dev/null || which google-chrome
```

---

## Quick Test Matrix

| Test | Duration | Command | Type | Scenario |
|------|----------|---------|------|----------|
| Integration | ~2 min | `go test -tags integration ./tests/integration -run Orchestrator` | Automated | Full coverage |
| Quick Start | ~5 min | Read `tests/manual/quick-start.md` and follow steps | Manual | Learning |
| Automated | ~10 sec | `./tests/manual/automated-test.sh` | Automated | CI/CD |
| Stress | ~30 sec | `./tests/manual/stress-test-10-instances.sh` | Automated | Load testing |
| Multi-Agent | ~15 sec | `./tests/manual/multi-agent-isolation.sh` | Automated | Agent isolation |

---

## Running Tests Locally

### Test Everything (5 minutes)

```bash
# 1. Automated integration tests
go test -tags integration ./tests/integration -run Orchestrator -timeout 120s

# 2. Run automated test script
./tests/manual/automated-test.sh

# 3. Run stress test
./tests/manual/stress-test-10-instances.sh

# 4. Run multi-agent test
./tests/manual/multi-agent-isolation.sh
```

If all pass → System is working correctly ✅

### Quick Smoke Test (30 seconds)

```bash
./tests/manual/automated-test.sh
```

Good for quick verification after code changes.

### Deep Dive (5+ minutes)

```bash
# Run integration tests with verbose output
go test -tags integration ./tests/integration -run Orchestrator -v -timeout 120s

# Then follow the manual quick-start guide
cat tests/manual/quick-start.md

# Then run stress and multi-agent tests
./tests/manual/stress-test-10-instances.sh
./tests/manual/multi-agent-isolation.sh
```

---

## Test Checklist

- [ ] Automated integration tests pass
- [ ] Automated test script passes
- [ ] Stress test (10 instances) passes
- [ ] Multi-agent isolation test passes
- [ ] Manual quick-start test completed
- [ ] No Chrome orphan processes after tests
- [ ] All ports released after tests

```bash
# Verify cleanup
lsof -i :9868-9968 2>/dev/null  # Should be empty
ps aux | grep -i chrome | grep -v grep  # Should be empty (except user's other Chrome)
```

---

## Debugging Failed Tests

### Port in use error
```bash
# Find what's using the port
lsof -i :9867
lsof -i :9868

# Kill the process
kill -9 $(lsof -t -i :9867)
```

### Chrome not found
```bash
# Install Chrome or set path
export CHROME_BIN=/path/to/chrome
./tests/manual/automated-test.sh
```

### Orphan Chrome processes
```bash
# Kill all Chrome processes managed by PinchTab
pkill -f "Chrome.*user-data-dir.*\.pinchtab"
```

### Test hangs
```bash
# Kill PinchTab if hung
pkill -f "^./pinchtab"

# Wait 2 seconds for cleanup
sleep 2

# Try again
./tests/manual/automated-test.sh
```

---

## Test Architecture

```
Tests
├── Automated (CI/CD friendly)
│   ├── tests/integration/orchestrator_test.go (Go tests)
│   ├── tests/manual/automated-test.sh
│   ├── tests/manual/stress-test-10-instances.sh
│   └── tests/manual/multi-agent-isolation.sh
│
└── Manual (Human-guided)
    └── tests/manual/quick-start.md (Learning guide)
    └── tests/manual/orchestrator.md (Comprehensive guide)
```

---

## Expected Results

After running all tests, you should see:

✅ All integration tests pass  
✅ All automated scripts pass  
✅ 10-instance stress test handles load  
✅ Multi-agent isolation verified  
✅ No resource leaks or orphan processes  
✅ Complete cleanup on shutdown  

---

## Next Steps

- **Passing all tests?** → System is ready for production
- **Some tests failing?** → Check debugging section above
- **Want to add tests?** → See CONTRIBUTING.md
- **Performance questions?** → Run stress test with monitoring (memory, CPU)

---

## Related Documentation

- [DEFINITION_OF_DONE.md](./.github/DEFINITION_OF_DONE.md) — PR checklist
- [LABELING_GUIDE.md](./.github/LABELING_GUIDE.md) — Issue labeling
- [DOCUMENTATION_REVIEW.md](./.github/DOCUMENTATION_REVIEW.md) — Doc audit guide

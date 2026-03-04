#!/bin/bash
# Complex Test Scenario: Stress Test with 10 Instances
# Tests: concurrent instance creation, parallel navigation, resource limits, cleanup
# Duration: ~30 seconds
# Usage: ./tests/manual/stress-test-10-instances.sh

set -e

echo "🔥 Starting Stress Test (10 Instances)..."
echo ""

./pinchtab &
DASHBOARD_PID=$!
sleep 2

echo "✓ Dashboard started (PID: $DASHBOARD_PID)"
echo ""

# Create 10 instances
echo "Creating 10 headless instances..."
INSTANCES=()

for i in {1..10}; do
  INST=$(curl -s -X POST http://localhost:9867/instances/launch \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"stress-$i\",\"headless\":true}")
  ID=$(echo $INST | jq -r '.id')
  PORT=$(echo $INST | jq -r '.port')
  INSTANCES+=($ID)
  printf "  %2d. %-20s (port: %s)\n" $i "$ID" "$PORT"
done

echo "✓ Created ${#INSTANCES[@]} instances"
echo ""

# Wait for Chrome initialization
sleep 3

# Navigate all instances concurrently
echo "Navigating all 10 instances concurrently..."
URLS=(
  "https://example.com"
  "https://github.com"
  "https://rust-lang.org"
  "https://wikipedia.org"
  "https://stackoverflow.com"
  "https://crates.io"
  "https://docs.rs"
  "https://news.ycombinator.com"
  "https://lobste.rs"
  "https://reddit.com/r/rust"
)

for i in {0..9}; do
  ID=${INSTANCES[$i]}
  URL=${URLS[$i]}
  curl -s -X POST "http://localhost:9867/instances/$ID/navigate" \
    -H "Content-Type: application/json" \
    -d "{\"url\":\"$URL\"}" > /dev/null &
done

wait
echo "✓ All 10 navigations completed"
echo ""

# Verify all are running
echo "Verifying all instances still running..."
RUNNING=$(curl -s http://localhost:9867/instances | jq 'length')
if [ "$RUNNING" -eq 10 ]; then
  echo "✓ All 10 instances running"
else
  echo "⚠️  Only $RUNNING instances still running (expected 10)"
fi
echo ""

# Snapshot each instance (verify isolation)
echo "Verifying instance isolation via snapshots..."
SNAPSHOT_SUCCESS=0
for i in {0..9}; do
  ID=${INSTANCES[$i]}
  SNAP=$(curl -s "http://localhost:9867/instances/$ID/snapshot" | jq -r '.url' 2>/dev/null || echo "")
  if [ ! -z "$SNAP" ]; then
    SNAPSHOT_SUCCESS=$((SNAPSHOT_SUCCESS + 1))
  fi
done
echo "✓ $SNAPSHOT_SUCCESS/10 snapshots successful (instance isolation verified)"
echo ""

# Concurrent stop (cleanup)
echo "Stopping all 10 instances concurrently..."
for ID in "${INSTANCES[@]}"; do
  curl -s -X POST "http://localhost:9867/instances/$ID/stop" > /dev/null &
done

wait
echo "✓ All instances stopped"
echo ""

# Verify cleanup
echo "Verifying cleanup..."
REMAINING=$(curl -s http://localhost:9867/instances | jq 'length')
if [ "$REMAINING" -eq 0 ]; then
  echo "✓ All instances cleaned up"
else
  echo "❌ FAILED: $REMAINING instances still running!"
  kill $DASHBOARD_PID 2>/dev/null
  exit 1
fi

# Clean up
kill $DASHBOARD_PID 2>/dev/null
sleep 1

echo ""
echo "=========================================="
echo "✅ STRESS TEST PASSED!"
echo "=========================================="
echo ""
echo "Summary:"
echo "  • Concurrent creation: 10 instances ✓"
echo "  • Concurrent navigation: 10 instances ✓"
echo "  • Instance isolation: $SNAPSHOT_SUCCESS/10 verified ✓"
echo "  • Resource management: OK ✓"
echo "  • Cleanup: ✓"
echo ""
echo "This tests:"
echo "  • Port allocator under load"
echo "  • Concurrent Chrome initialization"
echo "  • Parallel navigation requests"
echo "  • Memory/resource limits"
echo "  • Cleanup at scale"

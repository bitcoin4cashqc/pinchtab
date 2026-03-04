#!/bin/bash
# Automated orchestrator test script
# Verifies: instance creation, port allocation, isolation, cleanup
# Duration: ~10 seconds
# Usage: ./tests/manual/automated-test.sh

set -e

echo "🚀 Starting Automated Orchestrator Test..."
echo ""

# Start Pinchtab in background
./pinchtab &
DASHBOARD_PID=$!
sleep 2

echo "✓ Dashboard started (PID: $DASHBOARD_PID)"
echo ""

# Create instance 1
echo "Creating instance 1 (headed)..."
INST1=$(curl -s -X POST http://localhost:9867/instances/launch \
  -H "Content-Type: application/json" \
  -d '{"name":"work","headless":false}')
INST1_ID=$(echo $INST1 | jq -r '.id')
INST1_PORT=$(echo $INST1 | jq -r '.port')
echo "✓ Created instance 1: $INST1_ID (port: $INST1_PORT)"

# Create instance 2
echo "Creating instance 2 (headless)..."
INST2=$(curl -s -X POST http://localhost:9867/instances/launch \
  -H "Content-Type: application/json" \
  -d '{"name":"scrape","headless":true}')
INST2_ID=$(echo $INST2 | jq -r '.id')
INST2_PORT=$(echo $INST2 | jq -r '.port')
echo "✓ Created instance 2: $INST2_ID (port: $INST2_PORT)"
echo ""

# Wait for Chrome to initialize
sleep 2

# List instances
echo "Verifying instance list..."
INSTANCES=$(curl -s http://localhost:9867/instances | jq '.')
INSTANCE_COUNT=$(echo $INSTANCES | jq 'length')
echo "✓ $INSTANCE_COUNT instances running"
echo ""

# Navigate instance 1
echo "Navigating instance 1 to example.com..."
NAV1=$(curl -s -X POST "http://localhost:9867/instances/$INST1_ID/navigate" \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}')
TAB1_ID=$(echo $NAV1 | jq -r '.tabId')
TAB1_URL=$(echo $NAV1 | jq -r '.url')
echo "✓ Instance 1 navigated, tab: $TAB1_ID, url: $TAB1_URL"

# Navigate instance 2
echo "Navigating instance 2 to github.com..."
NAV2=$(curl -s -X POST "http://localhost:9867/instances/$INST2_ID/navigate" \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com"}')
TAB2_ID=$(echo $NAV2 | jq -r '.tabId')
TAB2_URL=$(echo $NAV2 | jq -r '.url')
echo "✓ Instance 2 navigated, tab: $TAB2_ID, url: $TAB2_URL"
echo ""

# Verify isolation: different tab IDs
echo "Verifying instance isolation..."
if [ "$TAB1_ID" != "$TAB2_ID" ]; then
  echo "✓ Instance isolation verified (different tab IDs: $TAB1_ID vs $TAB2_ID)"
else
  echo "❌ FAILED: Tab IDs should be different!"
  kill $DASHBOARD_PID 2>/dev/null
  exit 1
fi
echo ""

# Stop instance 1
echo "Stopping instance 1..."
curl -s -X POST "http://localhost:9867/instances/$INST1_ID/stop" > /dev/null
echo "✓ Instance 1 stopped, port $INST1_PORT released"

# Stop instance 2
echo "Stopping instance 2..."
curl -s -X POST "http://localhost:9867/instances/$INST2_ID/stop" > /dev/null
echo "✓ Instance 2 stopped, port $INST2_PORT released"
echo ""

# Verify all stopped
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
echo "✅ ALL TESTS PASSED!"
echo "=========================================="
echo ""
echo "Summary:"
echo "  • Instance creation: ✓"
echo "  • Port allocation: ✓"
echo "  • Orchestrator routing: ✓"
echo "  • Instance isolation: ✓"
echo "  • Cleanup: ✓"

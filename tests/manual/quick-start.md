# Quick Start Manual Test (10 Steps)

Manual step-by-step test to verify basic orchestrator functionality with visual/real-time feedback.

**Duration:** ~5 minutes  
**Requirements:** PinchTab built, ports 9867-9968 available, Chrome installed

---

## Prerequisites

```bash
cd ~/dev/pinchtab
go build -o pinchtab ./cmd/pinchtab
```

---

## Test Steps

### 1. Start PinchTab

```bash
./pinchtab
```

**Expected output:**
```
INFO dashboard listening addr=127.0.0.1:9867
INFO port allocator initialized start=9868 end=9968
```

**Verify:**
- ✓ Dashboard listening on 9867
- ✓ Port allocator ready (9868-9968)

---

### 2. Create First Instance (Headed)

```bash
curl -X POST http://localhost:9867/instances/launch \
  -H "Content-Type: application/json" \
  -d '{
    "name":"work",
    "headless":false
  }'
```

**Expected response:**
```json
{
  "id": "inst_XXXXXXXX",
  "profileId": "prof_YYYYYYYY",
  "profileName": "work",
  "port": "9868",
  "headless": false,
  "status": "starting",
  "startTime": "2026-02-28T20:15:00Z"
}
```

**Verify:**
- ✓ Hash-based instance ID (inst_XXXXXXXX format)
- ✓ Hash-based profile ID (prof_YYYYYYYY format)
- ✓ Port auto-allocated (9868)
- ✓ Status is "starting" (not "running" yet)
- ✓ **Within 2 seconds:** Chrome window opens on screen

---

### 3. Create Second Instance (Headless)

```bash
curl -X POST http://localhost:9867/instances/launch \
  -H "Content-Type: application/json" \
  -d '{
    "name":"scrape",
    "headless":true
  }'
```

**Expected response:**
```json
{
  "id": "inst_ZZZZZZZZ",
  "profileId": "prof_WWWWWWWW",
  "profileName": "scrape",
  "port": "9869",
  "headless": true,
  "status": "starting",
  "startTime": "2026-02-28T20:15:05Z"
}
```

**Verify:**
- ✓ Different instance ID
- ✓ Different port (9869)
- ✓ Same port range (auto-allocation working)
- ✓ No Chrome window (headless=true)

---

### 4. List All Instances

```bash
curl http://localhost:9867/instances
```

**Expected response:**
```json
[
  {
    "id": "inst_XXXXXXXX",
    "profileId": "prof_YYYYYYYY",
    "profileName": "work",
    "port": "9868",
    "headless": false,
    "status": "running",
    "startTime": "2026-02-28T20:15:00Z"
  },
  {
    "id": "inst_ZZZZZZZZ",
    "profileId": "prof_WWWWWWWW",
    "profileName": "scrape",
    "port": "9869",
    "headless": true,
    "status": "running",
    "startTime": "2026-02-28T20:15:05Z"
  }
]
```

**Verify:**
- ✓ Both instances listed
- ✓ Both have "running" status
- ✓ Hash-based IDs on both

---

### 5. Navigate Instance 1 (via Orchestrator Proxy)

```bash
curl -X POST http://localhost:9867/instances/inst_XXXXXXXX/navigate \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com"
  }'
```

Replace `inst_XXXXXXXX` with actual ID from step 2.

**Expected response:**
```json
{
  "tabId": "tab_MMMMMMMM",
  "url": "https://example.com",
  "title": "Example Domain"
}
```

**Verify:**
- ✓ Hash-based tab ID (tab_MMMMMMMM format)
- ✓ Instance 1 Chrome window navigates to example.com
- ✓ URL matches

---

### 6. Navigate Instance 2 (via Orchestrator Proxy)

```bash
curl -X POST http://localhost:9867/instances/inst_ZZZZZZZZ/navigate \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://github.com"
  }'
```

Replace `inst_ZZZZZZZZ` with actual ID from step 3.

**Expected response:**
```json
{
  "tabId": "tab_NNNNNNNN",
  "url": "https://github.com",
  "title": "GitHub"
}
```

**Verify:**
- ✓ Different tab ID (tab_NNNNNNNN vs tab_MMMMMMMM)
- ✓ Instance 2 (headless) navigates silently
- ✓ Instance 1 window still shows example.com (not affected)

---

### 7. Get Snapshot from Instance 1

```bash
curl "http://localhost:9867/instances/inst_XXXXXXXX/snapshot" \
  -o snapshot.json
cat snapshot.json | jq '.url'
```

**Expected:**
- ✓ Returns full page snapshot of example.com
- ✓ Shows isolation: only sees inst_XXXXXXXX's content

---

### 8. Stop Instance 1

```bash
curl -X POST http://localhost:9867/instances/inst_XXXXXXXX/stop
```

**Expected response:**
```json
{
  "status": "stopped",
  "id": "inst_XXXXXXXX"
}
```

**Verify:**
- ✓ Chrome window closes
- ✓ Instance 2 still running (headless, invisible)
- ✓ Port 9868 released back to allocator

---

### 9. Create Third Instance (Reuses Released Port)

```bash
curl -X POST http://localhost:9867/instances/launch \
  -H "Content-Type: application/json" \
  -d '{
    "name":"test",
    "headless":true
  }'
```

**Expected:**
- ✓ New instance gets port 9868 (reused from step 8)
- ✓ New instance ID generated

---

### 10. Stop All Instances

```bash
curl -X POST http://localhost:9867/instances/inst_ZZZZZZZZ/stop
curl -X POST http://localhost:9867/instances/inst_TTTTTTTT/stop
```

**Verify:**
- ✓ All instances stopped
- ✓ All ports released back to allocator
- ✓ Dashboard still running on 9867

---

## Summary

✅ If all steps pass:
- Hash-based ID generation works
- Port allocation and reuse works
- Headed instances open Chrome windows
- Headless instances run silently
- Orchestrator proxy routes work
- Instance isolation verified
- Cleanup works correctly

**Expected total time:** ~5 minutes

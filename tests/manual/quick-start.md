# Quick Start Manual Test (7 Steps)

Manual step-by-step test to verify basic orchestrator functionality: instance creation, tab creation, and find endpoint.

**Duration:** ~5 minutes  
**Requirements:** PinchTab built, ports 9867-9968 available, Chrome installed  
**Updated:** 2026-03-04 (uses tab creation + find endpoint, not navigate)

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

### 5. Create Tab in Instance 1

Get the port for instance 1 from step 2 (9868). Then create a tab:

```bash
curl -X POST http://localhost:9868/tabs \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Expected response:**
```json
{
  "id": "tab_MMMMMMMM",
  "url": "about:blank",
  "title": ""
}
```

**Verify:**
- ✓ Hash-based tab ID (tab_MMMMMMMM format)
- ✓ Tab created successfully

---

### 6. Create Tab in Instance 2

Get the port for instance 2 from step 3 (9869). Then create a tab:

```bash
curl -X POST http://localhost:9869/tabs \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Expected response:**
```json
{
  "id": "tab_NNNNNNNN",
  "url": "about:blank",
  "title": ""
}
```

**Verify:**
- ✓ Different tab ID (tab_NNNNNNNN vs tab_MMMMMMMM)
- ✓ Instance isolation: each has unique tab

---

### 7. Use Find Endpoint (Search for Elements)

Now test the find endpoint on instance 1:

```bash
curl -X POST http://localhost:9868/find \
  -H "Content-Type: application/json" \
  -d '{
    "text": "example"
  }'
```

**Expected response:**
```json
{
  "refs": [
    {"ref": "e1", "text": "..."},
    ...
  ]
}
```

**Verify:**
- ✓ Find endpoint responds with element references
- ✓ Can search for elements by text

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

---

## Summary

✅ If all steps pass:
- Instance creation works (headed + headless)
- Port allocation and reuse works
- Tab creation per instance works
- Find endpoint works
- Instance isolation verified
- Cleanup works correctly

**Expected total time:** ~5 minutes

---



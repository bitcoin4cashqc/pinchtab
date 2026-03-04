# CLI Quick Reference

PinchTab CLI provides **management commands** for server control and configuration. Browser automation is done via the **HTTP API** (see [API Reference](endpoints.md)).

---

## Management Commands

### Server Status

```bash
pinchtab health
```

Returns server health and mode.

---

## Profiles

### List available profiles

```bash
pinchtab profiles
```

Output:
```
📋 Available Profiles:

  👤 default
  👤 work
  👤 shopping
```

---

## Instances

### List running instances

```bash
pinchtab instances
```

Output:
```
🚀 Running Instances:

  ▶️ inst-abc123 (port 9868, headless)
     → http://localhost:9868

  ▶️ inst-def456 (port 9869, headed)
     → http://localhost:9869
```

---

## Tabs

### List all open tabs (across all instances)

```bash
pinchtab tabs
```

Output:
```
📑 Open Tabs:

  [tab-1] GitHub - Build software better
       https://github.com/pinchtab/pinchtab

  [tab-2] Google Search
       https://google.com
```

---

## Connection

### Get instance URL for a profile

```bash
pinchtab connect myprofile
```

Output:
```
http://localhost:9868
```

Useful for getting the exact URL and port of a running instance:
```bash
INSTANCE_URL=$(pinchtab connect work)
curl "$INSTANCE_URL/snapshot?filter=interactive"
```

---

## Configuration

### Initialize config file

```bash
pinchtab config init
```

Creates default config at `~/.config/pinchtab/config.json` (macOS/Linux) or `%APPDATA%\pinchtab\config.json` (Windows).

### Display config

```bash
# JSON format
pinchtab config show

# YAML format
pinchtab config show --format yaml
```

### Set config values

```bash
# Single values
pinchtab config set server.port 9999
pinchtab config set chrome.headless false
pinchtab config set chrome.maxTabs 50
pinchtab config set orchestrator.strategy session
```

Supported keys:
- `server.port` — Server port
- `server.stateDir` — State directory
- `server.profileDir` — Profile directory
- `server.token` — API token
- `chrome.headless` — Headless mode (true/false)
- `chrome.maxTabs` — Max tabs per instance
- `chrome.noRestore` — Don't restore session (true/false)
- `orchestrator.strategy` — Allocation strategy (simple/session/explicit)
- `orchestrator.allocationPolicy` — Policy (fcfs/round_robin/random)
- `orchestrator.instancePortStart` — Starting port
- `orchestrator.instancePortEnd` — Ending port
- `timeouts.actionSec` — Action timeout (seconds)
- `timeouts.navigateSec` — Navigation timeout (seconds)

### Merge config with JSON

```bash
pinchtab config patch '{"chrome": {"headless": false, "maxTabs": 100}}'
```

### Validate config

```bash
pinchtab config validate
```

Checks:
- Required fields present
- Port ranges valid
- Strategy/policy enums valid
- Timeout values non-negative

---

## Help

### Show command help

```bash
pinchtab help
```

Displays all available commands and basic examples.

---

## Environment Variables

### Server Configuration

```bash
# Server port
BRIDGE_PORT=9867

# Server bind address
BRIDGE_BIND=127.0.0.1

# Config file location
BRIDGE_CONFIG=~/.config/pinchtab/config.json
```

### Browser Settings

```bash
# Headless mode
BRIDGE_HEADLESS=true

# Chrome binary path
CHROME_BINARY=/usr/bin/google-chrome

# Profile directory
BRIDGE_PROFILE=~/.config/pinchtab/chrome-profile
```

### Orchestrator (Dashboard Mode)

```bash
# Allocation strategy
PINCHTAB_STRATEGY=simple

# Instance selection policy
PINCHTAB_ALLOCATION_POLICY=fcfs

# Instance port range
INSTANCE_PORT_START=9868
INSTANCE_PORT_END=9968
```

### Security & Debugging

```bash
# API authentication token
BRIDGE_TOKEN=secret-key

# Stealth level (light/medium/full)
BRIDGE_STEALTH=light

# Block ads/images/media
BRIDGE_BLOCK_ADS=true
BRIDGE_BLOCK_IMAGES=false
BRIDGE_BLOCK_MEDIA=false
```

---

## Examples

### Start server with custom port

```bash
BRIDGE_PORT=9999 pinchtab
```

### Initialize and configure

```bash
# Create default config
pinchtab config init

# Customize it
pinchtab config set server.port 9999
pinchtab config set chrome.headless false

# View it
pinchtab config show --format yaml

# Start server
pinchtab
```

### List and check instances

```bash
# See all instances
pinchtab instances

# See all tabs
pinchtab tabs

# Check server
pinchtab health
```

### Get instance URL and use HTTP API

```bash
# Get URL
INSTANCE_URL=$(pinchtab connect work)

# Use HTTP API
curl "$INSTANCE_URL/snapshot?filter=interactive&compact=true" | jq .

curl -X POST "$INSTANCE_URL/navigate" \
  -d '{"url":"https://example.com"}' \
  -H "Content-Type: application/json"
```

---

## For Browser Automation

Use the **HTTP API** for browser control:

```bash
curl -X POST http://localhost:9867/navigate \
  -d '{"url":"https://example.com"}' \
  -H "Content-Type: application/json"

curl "http://localhost:9867/snapshot?filter=interactive&compact=true"

curl -X POST http://localhost:9867/action \
  -d '{"kind":"click","ref":"e5"}' \
  -H "Content-Type: application/json"
```

Or use a client library (Playwright, Puppeteer, Cypress).

See [API Reference](endpoints.md) for complete HTTP endpoints.

# Configuration

Complete reference for PinchTab configuration via config file, environment variables, and CLI commands.

## Quick Start

```bash
# Initialize config file
pinchtab config init

# Set values
pinchtab config set server.port 9999
pinchtab config set chrome.headless false

# View config
pinchtab config show --format yaml

# Validate
pinchtab config validate
```

---

## Configuration File

PinchTab uses a JSON config file at:
- **macOS/Linux:** `~/.config/pinchtab/config.json`
- **Windows:** `%APPDATA%\pinchtab\config.json`

Override location with `BRIDGE_CONFIG` environment variable:
```bash
BRIDGE_CONFIG=/custom/path/config.json pinchtab
```

### Config Structure

```json
{
  "port": "9867",
  "stateDir": "~/.config/pinchtab",
  "profileDir": "~/.config/pinchtab/chrome-profile",
  "headless": true,
  "maxTabs": 20,
  "noRestore": false,
  "timeoutSec": 15,
  "navigateSec": 30,
  "strategy": "simple",
  "allocationPolicy": "fcfs",
  "instancePortStart": 9868,
  "instancePortEnd": 9968
}
```

---

## Config Management (CLI)

### View configuration

```bash
# JSON
pinchtab config show

# YAML
pinchtab config show --format yaml
```

### Set individual values

```bash
pinchtab config set <key> <value>
```

Examples:
```bash
pinchtab config set server.port 9999
pinchtab config set chrome.headless false
pinchtab config set chrome.maxTabs 50
pinchtab config set orchestrator.strategy session
pinchtab config set orchestrator.allocationPolicy round_robin
pinchtab config set timeouts.actionSec 30
```

### Merge JSON object

```bash
pinchtab config patch '<json>'
```

Examples:
```bash
pinchtab config patch '{"chrome": {"headless": false, "maxTabs": 100}}'
pinchtab config patch '{"server": {"port": "9999"}}'
```

### Validate configuration

```bash
pinchtab config validate
```

Checks:
- Required fields (`port`)
- Port ranges (`instancePortStart` < `instancePortEnd`)
- Enum values (`strategy`, `allocationPolicy`)
- Non-negative timeouts

---

## Configuration Sections

### server

```json
{
  "port": "9867",
  "stateDir": "~/.config/pinchtab",
  "profileDir": "~/.config/pinchtab/chrome-profile",
  "token": "api-key",
  "cdpUrl": "ws://localhost:9222"
}
```

**port** — Server HTTP port (default: `9867`)  
**stateDir** — Directory for instance state (default: `~/.config/pinchtab`)  
**profileDir** — Directory for Chrome profiles (default: `~/.config/pinchtab/chrome-profile`)  
**token** — API authentication token (optional)  
**cdpUrl** — Custom Chrome DevTools Protocol URL (optional)  

### chrome

```json
{
  "headless": true,
  "maxTabs": 20,
  "noRestore": false
}
```

**headless** — Run Chrome in headless mode (default: `true`)  
**maxTabs** — Maximum tabs per instance (default: `20`)  
**noRestore** — Don't restore session on startup (default: `false`)  

### orchestrator

```json
{
  "strategy": "simple",
  "allocationPolicy": "fcfs",
  "instancePortStart": 9868,
  "instancePortEnd": 9968
}
```

**strategy** — Allocation strategy: `simple`, `session`, `explicit` (default: none)  
**allocationPolicy** — Instance selection: `fcfs`, `round_robin`, `random` (default: `fcfs`)  
**instancePortStart** — Starting port for instances (default: `9868`)  
**instancePortEnd** — Ending port for instances (default: `9968`)  

### timeouts

```json
{
  "timeoutSec": 15,
  "navigateSec": 30
}
```

**timeoutSec** — Default action timeout in seconds (default: `15`)  
**navigateSec** — Page navigation timeout in seconds (default: `30`)  

---

## Configuration Priority (Precedence)

## Environment Variables

### Port & Network

| Variable | Default | Description |
|---|---|---|
| `BRIDGE_PORT` | `9867` | HTTP server port |
| `BRIDGE_BIND` | `127.0.0.1` | Bind address (127.0.0.1 = localhost only, 0.0.0.0 = all interfaces) |

### Browser & Chrome

| Variable | Default | Description |
|---|---|---|
| `CHROME_BINARY` | Auto-detect | Path to Chrome/Chromium binary |
| `BRIDGE_HEADLESS` | `true` | Run Chrome headless (no visible window) |
| `BRIDGE_PROFILE` | Default profile | Chrome profile name (stored in `~/.pinchtab/profiles/{name}`) |

### Stealth & Detection

| Variable | Default | Description |
|---|---|---|
| `BRIDGE_STEALTH` | `light` | Stealth level: `light`, `medium`, `full` (higher = more bot detection bypass, slower) |

### Content Filtering

| Variable | Default | Description |
|---|---|---|
| `BRIDGE_BLOCK_ADS` | `false` | Block ad domains (speeds up loading) |
| `BRIDGE_BLOCK_IMAGES` | `false` | Block image loading |
| `BRIDGE_BLOCK_MEDIA` | `false` | Block video/audio resources |

### Security & Authentication

| Variable | Default | Description |
|---|---|---|
| `BRIDGE_TOKEN` | Disabled | API authentication token (if set, all requests must include `Authorization: Bearer {token}`) |

### Debugging & Logging

| Variable | Default | Description |
|---|---|---|
| `BRIDGE_DEBUG` | `false` | Enable debug logging (verbose output) |
| `BRIDGE_LOG_LEVEL` | `info` | Log level: `debug`, `info`, `warn`, `error` |

### Dashboard

| Variable | Default | Description |
|---|---|---|
| `BRIDGE_DASHBOARD_PORT` | Same as `BRIDGE_PORT` | Dashboard HTTP port (usually same as API server) |
| `BRIDGE_NO_DASHBOARD` | `false` | Disable dashboard (API-only mode) |

---

---

## Configuration Priority (Precedence)

If multiple sources set the same value:

1. **Environment variables** (highest priority)
2. **Config file**
3. **Built-in defaults** (lowest priority)

**Example:**
```bash
# Config file has: port = 9867
# But env var overrides it:
BRIDGE_PORT=9999 pinchtab  # Uses port 9999
```

---

## Usage Examples

### Quick Examples

```bash
# Default (headless, localhost:9867)
pinchtab

# Custom port
BRIDGE_PORT=9999 pinchtab

# Headed mode with work profile
BRIDGE_HEADLESS=false BRIDGE_PROFILE=work pinchtab

# Network-accessible (requires auth)
BRIDGE_BIND=0.0.0.0 BRIDGE_TOKEN=secret pinchtab
```

### Using Config File + Environment Variables

```bash
# Initialize config
pinchtab config init

# Override specific values with env vars
BRIDGE_PORT=9999 pinchtab
```

Environment variables take precedence over config file values.

---

## Related Documentation

- [CLI Quick Reference](cli-quick-reference.md) — CLI commands
- [API Reference](endpoints.md) — HTTP endpoints
- [Getting Started](../get-started.md) — Quick setup

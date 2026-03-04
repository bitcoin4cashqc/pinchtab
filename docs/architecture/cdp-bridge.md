# CDP Bridge: Navigation & Wait Strategies

How Pinchtab handles page readiness between navigation and operations at the Chrome DevTools Protocol level.

## The Problem

When you navigate to a URL, the page goes through several states before it's usable:

```
Page.navigate → DOM parsing → DOMContentLoaded → scripts execute → images/fonts load → load event
```

Different operations need different levels of readiness:
- **Snapshot/text/find** need the DOM parsed (a11y tree is ready)
- **Screenshot/PDF** need visual rendering complete (fonts, images loaded)
- **JavaScript eval** needs DOM interactive (scripts can run)

Waiting for full load on every operation wastes time. Waiting too little gives empty results.

## Wait Modes

Pinchtab supports four wait modes, set via the `waitFor` parameter:

| Mode | What it checks | Speed | Use case |
|------|---------------|-------|----------|
| `none` | Nothing — returns immediately after navigation | Fastest | When you'll poll yourself |
| `dom` | `document.readyState == interactive` | Fast (~200ms) | Snapshot, text, find, eval |
| `complete` | `document.readyState == complete` | Medium | Screenshot, PDF |
| `networkidle` | readyState complete + no URL changes for 500ms | Slow (~3s) | SPAs, heavy JS pages |

### Why not always use `networkidle`?

Single-page apps never truly go idle — analytics pings, WebSocket connections, and lazy-loading keep firing network requests. `networkidle` can time out on these pages. For most operations, `dom` is sufficient and much faster.

## Smart Defaults

Each endpoint picks the right wait mode automatically:

| Endpoint | Default wait | Why |
|----------|-------------|-----|
| `GET /snapshot` | `dom` | A11y tree is built from parsed DOM |
| `GET /text` | `dom` | Text extraction reads DOM nodes |
| `POST /find` | `dom` | Searches the a11y snapshot |
| `POST /evaluate` | `dom` | JS can execute once DOM is interactive |
| `GET /screenshot` | `complete` | Needs rendered pixels, fonts, images |
| `GET /pdf` | `complete` | Chrome's printToPDF waits for complete internally |

Override any default with `waitFor`:

```bash
# Force networkidle for a heavy SPA
curl "localhost:9867/snapshot?url=https://spa-app.com&waitFor=networkidle"
```

## The `url` Parameter

All read endpoints accept an optional `url` parameter. When provided, the handler navigates first, waits using the smart default, then performs the operation — all in one HTTP call.

```bash
# Without url (two calls)
curl -X POST localhost:9867/navigate -d '{"url":"https://example.com"}'
curl localhost:9867/snapshot

# With url (one call)
curl "localhost:9867/snapshot?url=https://example.com"
```

This works on all read endpoints:

```bash
# Snapshot
curl "localhost:9867/snapshot?url=https://example.com"

# Text extraction
curl "localhost:9867/text?url=https://example.com"

# Screenshot
curl "localhost:9867/screenshot?url=https://example.com"

# PDF
curl "localhost:9867/pdf?url=https://example.com"

# Find (POST with JSON body)
curl -X POST localhost:9867/find -d '{"query":"Sign In","url":"https://example.com"}'

# Evaluate (POST with JSON body)
curl -X POST localhost:9867/evaluate -d '{"expression":"document.title","url":"https://example.com"}'
```

## CDP Lifecycle Events

Under the hood, Pinchtab uses these Chrome DevTools Protocol mechanisms:

### NavigatePage

`Page.navigate` + polls `document.readyState` every 200ms until `interactive` or `complete`:

```
Page.navigate(url)
  └─ poll every 200ms: Runtime.evaluate("document.readyState")
       └─ "loading"     → keep polling
       └─ "interactive" → DOM ready, return
       └─ "complete"    → fully loaded, return
```

### waitForNavigationState

Extends the basic navigation with configurable wait strategies:

- **`dom`**: Single `document.readyState` check
- **`selector`**: `DOM.querySelector` + wait for element visibility (via `chromedp.WaitVisible`)
- **`networkidle`**: Polls readyState + URL stability — requires 2 consecutive "complete" checks with stable URL (~500ms apart, up to 3s)

### Auto-snapshot (find)

`POST /find` always takes a fresh accessibility snapshot before searching. This uses `Accessibility.getFullAXTree` to build the element cache, then runs semantic matching against it. No stale data.

## Implementation

The `ensureNavigated` helper in `internal/handlers/navigate_helper.go`:

```go
func (h *Handlers) ensureNavigated(ctx context.Context, url, waitFor, defaultWait string) error {
    if url == "" {
        return nil  // no-op: use current page
    }
    bridge.NavigatePage(ctx, url)           // navigate + poll readyState
    h.waitForNavigationState(ctx, wait, "") // apply wait strategy
}
```

Each handler calls this before its main operation, passing the appropriate default wait mode.

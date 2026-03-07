# Showcase

The examples below keep the same snippet structure used in the full example pages:

## Full PinchTab Server

1. start `pinchtab`
2. launch an instance
3. open a tab in that instance
4. operate on the routed instance and tab endpoints

### Launch An Instance

```bash
INST=$(curl -s -X POST http://127.0.0.1:9867/instances/launch \
  -H "Content-Type: application/json" \
  -d '{"name":"my-profile","mode":"headless"}' \
  | jq -r '.id')

echo "$INST"
# Response
{
  "id": "inst_944a07ad",
  "profileId": "prof_910b1739",
  "profileName": "my-profile",
  "port": "9871",
  "headless": true,
  "status": "starting",
  "startTime": "2026-03-07T19:29:54.066542+01:00",
  "attached": false
}
```

### List Tabs On An Instance

```bash
curl -s "http://127.0.0.1:9867/instances/$INST/tabs" | jq .
# CLI alternative
pinchtab tabs --instance $INST
# Response
{
  "tabs": [
    {
      "id": "E291F6815F61C58B0C9EA9129F960744",
      "title": "GitHub - pinchtab/pinchtab",
      "type": "page",
      "url": "https://github.com/pinchtab/pinchtab"
    }
  ]
}
```

### Snapshot On An Instance

```bash
curl -s "http://127.0.0.1:9867/instances/$INST/snapshot?filter=interactive" | jq '.nodes[:5]'
# CLI alternative
pinchtab snap -i -c --instance $INST
# Response
e0:link "Skip to content"
e1:link "GitHub Homepage"
e2:link "pinchtab"
e5:button "Search or jump to…"
...
```

Use the full server when you want:
- profiles
- multiple instances
- attach support
- orchestration and routing
- dashboard visibility


## PinchTab bridge`

`pinchtab bridge` is the single-instance runtime.

Start it directly:

```bash
pinchtab bridge
```

Then use the single-instance API directly.

### Health

```bash
curl -s http://127.0.0.1:9867/health | jq .
# CLI alternative
pinchtab health
# Response
{
  "status": "ok",
  "tabs": 1
}
```

### List Tabs

```bash
curl -s http://127.0.0.1:9867/tabs | jq .
# CLI alternative
pinchtab tabs
# Response
{
  "tabs": [
    {
      "id": "BD78E40ED7400A4B0E73B99415E1B9EA",
      "title": "GitHub - pinchtab/pinchtab",
      "type": "page",
      "url": "https://github.com/pinchtab/pinchtab"
    }
  ]
}
```

### Export A PDF

```bash
curl -s http://127.0.0.1:9867/pdf > smoke.pdf
ls -lh smoke.pdf
# CLI alternative
pinchtab pdf -o smoke.pdf
# Response
Saved smoke.pdf (1494657 bytes)
```

Use bridge mode when you explicitly want:
- one browser runtime
- one process exposing the browser API directly
- no dashboard
- no multi-instance control plane

### Navigate

```bash
curl -s -X POST http://127.0.0.1:9867/navigate \
  -H "Content-Type: application/json" \
  -d '{"url":"https://github.com/pinchtab/pinchtab"}' | jq .
# CLI alternative
pinchtab nav https://github.com/pinchtab/pinchtab
# Response
{
  "tabId": "BD78E40ED7400A4B0E73B99415E1B9EA",
  "title": "GitHub - pinchtab/pinchtab",
  "url": "https://github.com/pinchtab/pinchtab"
}
```

### Snapshot

```bash
curl -s "http://127.0.0.1:9867/snapshot?filter=interactive" | jq .
# CLI alternative (compact format)
pinchtab snap -i -c
# Response
{
  "nodes": [
    { "ref": "e0", "role": "link", "name": "Skip to content" },
    { "ref": "e1", "role": "link", "name": "GitHub Homepage" },
    { "ref": "e14", "role": "button", "name": "Search or jump to…" }
  ]
}
```

### Extract Text

```bash
curl -s http://127.0.0.1:9867/text | jq .
# CLI alternative
pinchtab text
# Response
{
  "text": "High-performance browser automation bridge and multi-instance orchestrator...",
  "title": "GitHub - pinchtab/pinchtab",
  "url": "https://github.com/pinchtab/pinchtab"
}
```

### Click By Ref

```bash
curl -s -X POST http://127.0.0.1:9867/action \
  -H "Content-Type: application/json" \
  -d '{"kind":"click","ref":"e14"}' | jq .
# CLI alternative
pinchtab snap -i > /dev/null
pinchtab click e14
# Response
{
  "success": true,
  "result": {
    "clicked": true
  }
}
```

### Screenshot

```bash
curl -s http://127.0.0.1:9867/screenshot > smoke.jpg
ls -lh smoke.jpg
# CLI alternative
pinchtab ss -o smoke.jpg
# Response
Saved smoke.jpg (55876 bytes)
```
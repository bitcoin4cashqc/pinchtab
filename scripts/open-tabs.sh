#!/usr/bin/env bash
set -euo pipefail

COUNT="${1:-50}"
TAB_URL="${2:-about:blank}"

if ! [[ "$COUNT" =~ ^[0-9]+$ ]] || [ "$COUNT" -le 0 ]; then
  echo "Usage: $0 [count] [url]" >&2
  echo "Example: $0 50 https://example.com" >&2
  exit 1
fi

if command -v pinchtab >/dev/null 2>&1; then
  PINCHTAB=(pinchtab)
elif [ -x "./pinchtab" ]; then
  PINCHTAB=(./pinchtab)
else
  echo "pinchtab binary not found in PATH and ./pinchtab is not executable" >&2
  exit 1
fi

if ! "${PINCHTAB[@]}" health >/dev/null 2>&1; then
  echo "Pinchtab server is not reachable at ${PINCHTAB_URL:-http://127.0.0.1:9867}" >&2
  echo "Start it first with: pinchtab" >&2
  exit 1
fi

opened=0
failed=0

for ((i = 1; i <= COUNT; i++)); do
  if output=$("${PINCHTAB[@]}" open "$TAB_URL" 2>&1); then
    opened=$((opened + 1))
    printf '[%d/%d] opened\n' "$i" "$COUNT"
  else
    failed=$((failed + 1))
    printf '[%d/%d] failed: %s\n' "$i" "$COUNT" "$output" >&2

    if printf '%s' "$output" | grep -qi 'tab limit reached'; then
      echo "Stopped: tab limit reached after $opened opened tabs." >&2
      break
    fi
  fi

done

echo "Done. opened=$opened failed=$failed target=$COUNT url=$TAB_URL"

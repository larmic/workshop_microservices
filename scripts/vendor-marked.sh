#!/usr/bin/env bash
# Holt marked (Markdown-Renderer der Vorbereitungsseite) in der unten
# festgelegten Version von npm nach pages/vorbereitung/vendor/. So hängt die
# Seite an keinem CDN, das im Firmennetz gesperrt sein kann.
#
#   scripts/vendor-marked.sh           lädt die Datei
#   scripts/vendor-marked.sh --check   prüft, ob Datei und Version zusammenpassen
#
# Die Datei liegt nicht im Repo (.gitignore). Der Pages-Deploy und das
# Guard-Gate js in der Pipeline rufen das Skript auf. Lokal einmal ausführen,
# wer die Vorbereitungsseite ansehen will. Renovate hebt MARKED_VERSION an.
set -euo pipefail

MARKED_VERSION=18.3.0

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TARGET="$ROOT/pages/vorbereitung/vendor/marked.min.js"

if [ "${1:-}" = "--check" ]; then
  if ! head -c 300 "$TARGET" 2>/dev/null | grep -q "marked v$MARKED_VERSION "; then
    echo "pages/vorbereitung/vendor/marked.min.js ist nicht marked v$MARKED_VERSION: scripts/vendor-marked.sh ausführen" >&2
    exit 1
  fi
  echo "marked v$MARKED_VERSION liegt in pages/vorbereitung/vendor/"
  exit 0
fi

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
(cd "$tmp" && npm pack --silent "marked@$MARKED_VERSION" >/dev/null && tar -xzf marked-*.tgz)
mkdir -p "$(dirname "$TARGET")"
# Bis Version 15 liegt der Browser-Build als marked.min.js im Paket, danach
# als lib/marked.umd.js. Beide stellen das globale `marked` bereit.
src="$tmp/package/marked.min.js"
[ -f "$src" ] || src="$tmp/package/lib/marked.umd.js"
cp "$src" "$TARGET"
"$0" --check

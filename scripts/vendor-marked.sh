#!/usr/bin/env bash
# Holt marked (Markdown-Renderer der Vorbereitungsseite) in der unten
# festgelegten Version von npm und legt die Datei ins Repo. So hängt die Seite
# an keinem CDN, das im Firmennetz gesperrt sein kann.
#
#   scripts/vendor-marked.sh           lädt die Datei neu
#   scripts/vendor-marked.sh --check   prüft, ob Datei und Version zusammenpassen
#
# Renovate hebt nur MARKED_VERSION an. Danach dieses Skript ausführen und die
# Datei mit committen; das Guard-Gate js schlägt sonst fehl.
set -euo pipefail

MARKED_VERSION=18.2.0

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TARGET="$ROOT/pages/vorbereitung/vendor/marked.min.js"

if [ "${1:-}" = "--check" ]; then
  if ! head -c 300 "$TARGET" 2>/dev/null | grep -q "marked v$MARKED_VERSION "; then
    echo "pages/vorbereitung/vendor/marked.min.js ist nicht marked v$MARKED_VERSION: scripts/vendor-marked.sh ausführen" >&2
    exit 1
  fi
  echo "marked v$MARKED_VERSION liegt im Repo"
  exit 0
fi

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
(cd "$tmp" && npm pack --silent "marked@$MARKED_VERSION" >/dev/null && tar -xzf marked-*.tgz)
mkdir -p "$(dirname "$TARGET")"
cp "$tmp/package/marked.min.js" "$TARGET"
"$0" --check

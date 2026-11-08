#!/bin/sh
# Guard-Gates: deterministische Prüfungen, die vor jedem Merge grün sein müssen.
# Dieselben Kommandos laufen lokal, im Git-Hook (.githooks/pre-commit) und in der
# Pipeline (.github/workflows/gates.yml, ein Schritt pro Gate).
#
# Aufruf aus beliebigem Verzeichnis:
#   scripts/gates.sh            alle Gates in der Reihenfolge der Pipeline
#   scripts/gates.sh fast       die schnellen Gates des pre-commit-Hooks
#   scripts/gates.sh gofmt vet  einzelne Gates
#
# Gates: gofmt vet arch openapi palette test js unicode html compose
#
# Fehlt lokal ein Werkzeug (go, node, docker compose), wird das Gate mit Hinweis
# übersprungen – die Pipeline ist das verbindliche Gate. In der Pipeline (CI=true)
# ist ein fehlendes Werkzeug ein Fehler.
set -eu

ROOT=$(cd "$(dirname "$0")/.." && pwd)
SERVICES="$ROOT/services"
HTML_VALIDATE_VERSION=11.19.0
skipped=""
HTML_FILES="pages/landing/index.html pages/vorbereitung/index.html pages/feedback/index.html pages/material/index.html services/dashboard/static/index.html services/slides/index.html"

# Werkzeug vorhanden? Sonst lokal überspringen, in CI scheitern.
need() {
  command -v "$1" >/dev/null 2>&1 || missing "$1" "$2"
}

missing() {
  if [ "${CI:-}" = "true" ]; then
    echo "Guard-Gate $2: $1 fehlt" >&2
    exit 1
  fi
  echo "Guard-Gate $2: übersprungen, $1 ist nicht installiert (die Pipeline prüft es)"
  skipped="$skipped $2"
  return 1
}

gate_gofmt() {
  need go gofmt || return 0
  unformatted=$(cd "$SERVICES" && gofmt -l .)
  if [ -n "$unformatted" ]; then
    echo "Guard-Gate gofmt: nicht formatiert (beheben mit: cd services && gofmt -w .):" >&2
    echo "$unformatted" >&2
    return 1
  fi
}

gate_vet() {
  need go vet || return 0
  (cd "$SERVICES" && go vet ./...)
}

archtest() {
  need go "$1" || return 0
  (cd "$SERVICES" && go test -count=1 -run "$2" ./archtest/)
}

gate_arch() { archtest arch 'TestArchitecture|TestComponent'; }
gate_openapi() { archtest openapi 'TestOpenAPI'; }
gate_palette() { archtest palette 'TestPalette'; }
# Alle drei Architektur-Tests in einem Lauf (für den Hook).
gate_archtest() { archtest archtest '.'; }

gate_test() {
  need go test || return 0
  # archtest läuft als eigene Gates (arch, openapi, palette).
  # shellcheck disable=SC2046 # Paketliste soll in einzelne Argumente zerfallen
  (cd "$SERVICES" && go test $(go list ./... | grep -v '/archtest$'))
}

gate_js() {
  need node js || return 0
  (cd "$ROOT" && node scripts/check-inline-js.mjs)
}

gate_unicode() {
  need node unicode || return 0
  (cd "$ROOT" && node scripts/check-hidden-unicode.mjs)
}

gate_html() {
  need npx html || return 0
  # shellcheck disable=SC2086 # HTML_FILES ist bewusst eine Wortliste
  (cd "$ROOT" && npx --yes "html-validate@$HTML_VALIDATE_VERSION" --config scripts/htmlvalidate.json $HTML_FILES)
}

gate_compose() {
  docker compose version >/dev/null 2>&1 || missing "docker compose" compose || return 0
  # Die Kombinationen aus Makefile, docs/vorbereitung.md und Dashboard (COMPOSE_FILES).
  (
    cd "$SERVICES"
    docker compose -f docker-compose.yml -f docker-compose.infra.yml -f docker-compose.reference.yml config -q
    CUSTOM_BOOKING_PATH=./booking/custom docker compose -f docker-compose.yml -f docker-compose.infra.yml \
      -f docker-compose.reference.yml -f docker-compose.custom.yml config -q
  )
}

run() {
  echo "==> Guard-Gate $1"
  "gate_$1"
}

[ $# -eq 0 ] && set -- all
requested="$*"
for arg in "$@"; do
  case "$arg" in
    all) for g in gofmt vet arch openapi palette test js unicode html compose; do run "$g"; done ;;
    fast) for g in gofmt vet archtest js unicode; do run "$g"; done ;;
    gofmt | vet | arch | openapi | palette | archtest | test | js | unicode | html | compose) run "$arg" ;;
    *)
      echo "Unbekanntes Gate: $arg (erlaubt: all fast gofmt vet arch openapi palette test js unicode html compose)" >&2
      exit 2
      ;;
  esac
done
if [ -n "$skipped" ]; then
  echo "Guard-Gates ($requested) ohne Fehler, aber übersprungen:$skipped"
else
  echo "Guard-Gates grün ($requested)."
fi

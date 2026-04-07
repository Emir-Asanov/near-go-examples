#!/usr/bin/env bash
# test-all.sh — Build and test every near-go-examples project in sequence.
# Run from the repo root: ./test-all.sh
set -eo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

export GVM_ROOT="$HOME/.gvm"
export PATH="$HOME/.gvm/gos/go1.25.4/bin:$HOME/.gvm/bin:/opt/homebrew/bin:$PATH"
export NEAR_GO="$HOME/bin/near-go"

PASS=()
FAIL=()

run_project() {
    local label="$1"
    local dir="$2"
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "  PROJECT: $label"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    if bash "$dir/test.sh"; then
        PASS+=("$label")
    else
        FAIL+=("$label")
        echo "✗ FAILED: $label"
    fi
}

# ── Projects with unit + integration tests ───────────────────────────────────
run_project "greeting"        "$ROOT/greeting"
run_project "counter"         "$ROOT/counter"
run_project "donation"        "$ROOT/donation"
run_project "auction"         "$ROOT/auction"
run_project "collections"     "$ROOT/collections"
run_project "self-update"     "$ROOT/self-update"
run_project "yield-resume"    "$ROOT/yield-resume"
run_project "guest-book/base" "$ROOT/guest-book/base"
run_project "guest-book/v2"   "$ROOT/guest-book/v2"

# ── Summary ──────────────────────────────────────────────────────────────────
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "  SUMMARY"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
for p in "${PASS[@]}"; do echo "  ✓ $p"; done
for p in "${FAIL[@]}"; do echo "  ✗ $p"; done
echo ""

if [ ${#FAIL[@]} -eq 0 ]; then
    echo "All ${#PASS[@]} projects passed."
    exit 0
else
    echo "${#FAIL[@]} project(s) failed."
    exit 1
fi

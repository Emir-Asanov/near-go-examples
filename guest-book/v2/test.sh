#!/usr/bin/env bash
set -eo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

export GVM_ROOT="$HOME/.gvm"
export PATH="$HOME/.gvm/gos/go1.25.4/bin:$HOME/.gvm/bin:/opt/homebrew/bin:$PATH"
export NEAR_GO="$HOME/bin/near-go"

echo "=== [guest-book/base] Building (needed for migration test) ==="
(builtin cd ../base && "$NEAR_GO" build)

echo ""
echo "=== [guest-book/v2] Building ==="
"$NEAR_GO" build

echo ""
echo "=== [guest-book/v2] Unit tests ==="
"$NEAR_GO" test project

echo ""
echo "=== [guest-book/v2] Copying base.wasm for migration test ==="
cp ../base/main.wasm integration-tests/base.wasm
echo "  Copied base/main.wasm → v2/integration-tests/base.wasm"

echo ""
echo "=== [guest-book/v2] Integration tests (includes migration from base) ==="
cargo test --manifest-path integration-tests/Cargo.toml

echo ""
echo "✓ All guest-book/v2 tests passed"

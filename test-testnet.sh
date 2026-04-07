#!/usr/bin/env bash
# test-testnet.sh — Deploy and smoke-test contracts on NEAR testnet.
#
# Prerequisites:
#   1. near-cli-rs installed:  cargo install near-cli-rs
#   2. Logged in:              near login
#   3. near-go CLI available:  ~/bin/near-go
#   4. Testnet account set:    export NEAR_ACCOUNT=yourname.testnet
#
# Usage:
#   NEAR_ACCOUNT=yourname.testnet ./test-testnet.sh
#   ./test-testnet.sh greeting          # run only one contract
#   ./test-testnet.sh greeting counter  # run several

set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

export GVM_ROOT="$HOME/.gvm"
export PATH="$HOME/.gvm/gos/go1.25.4/bin:$HOME/.gvm/bin:$PATH"
export NEAR_GO="$HOME/bin/near-go"

NEAR="${NEAR_ACCOUNT:-}"
if [ -z "$NEAR" ]; then
    echo "Error: set NEAR_ACCOUNT=yourname.testnet"
    exit 1
fi

# ── Helpers ───────────────────────────────────────────────────────────────────

# near-cli-rs call shortcuts
tx_call() {   # tx_call <contract> <method> <json_args> [deposit]
    local contract="$1" method="$2" args="$3" deposit="${4:-0 NEAR}"
    near contract call-function as-transaction "$contract" "$method" \
        json-args "$args" \
        prepaid-gas '100.0 Tgas' \
        attached-deposit "$deposit" \
        sign-as "$NEAR" \
        network-config testnet \
        sign-with-keychain send
}

view_call() { # view_call <contract> <method> <json_args>
    local contract="$1" method="$2" args="$3"
    near contract call-function as-read-only "$contract" "$method" \
        json-args "$args" \
        network-config testnet now
}

deploy() {    # deploy <contract_id> [init_method] [init_args]
    local contract="$1" init_method="${2:-}" init_args="${3:-{}}"
    near contract deploy "$contract" \
        use-file main.wasm \
        without-init-call \
        network-config testnet \
        sign-with-keychain send
    if [ -n "$init_method" ]; then
        tx_call "$contract" "$init_method" "$init_args"
    fi
}

# Generate a unique sub-account for each test run so deploys don't collide
unique_id() { date +%s | tail -c 6; }

PASS=(); FAIL=()

run() {
    local name="$1"
    echo ""; echo "━━━  $name  ━━━"
    if "test_$name"; then PASS+=("$name"); else FAIL+=("$name"); fi
}

# ── Contract tests ────────────────────────────────────────────────────────────

test_greeting() {
    local id="greeting-$(unique_id).$NEAR"
    echo "Contract: $id"
    cd "$ROOT/greeting"
    "$NEAR_GO" build

    # Create sub-account and deploy
    near account create-account fund-myself "$id" '0.1 NEAR' \
        autogenerate-new-keypair save-to-keychain \
        sign-as "$NEAR" network-config testnet sign-with-keychain send

    deploy "$id"
    tx_call "$id" "set_greeting" '{"greeting":"Hello from testnet!"}'

    local result
    result=$(view_call "$id" "get_greeting" '{}')
    echo "get_greeting → $result"
    echo "$result" | grep -q "Hello from testnet!" && echo "✓ greeting passed"
}

test_counter() {
    local id="counter-$(unique_id).$NEAR"
    echo "Contract: $id"
    cd "$ROOT/counter"
    "$NEAR_GO" build

    near account create-account fund-myself "$id" '0.1 NEAR' \
        autogenerate-new-keypair save-to-keychain \
        sign-as "$NEAR" network-config testnet sign-with-keychain send

    deploy "$id" "init" '{}'

    local v0; v0=$(view_call "$id" "get_num" '{}')
    echo "initial value → $v0"

    tx_call "$id" "increment" '{}'
    tx_call "$id" "increment" '{}'

    local v2; v2=$(view_call "$id" "get_num" '{}')
    echo "after 2 increments → $v2"
    echo "$v2" | grep -q "2" && echo "✓ counter passed"
}

test_donation() {
    local id="donation-$(unique_id).$NEAR"
    echo "Contract: $id"
    cd "$ROOT/donation"
    "$NEAR_GO" build

    near account create-account fund-myself "$id" '0.5 NEAR' \
        autogenerate-new-keypair save-to-keychain \
        sign-as "$NEAR" network-config testnet sign-with-keychain send

    deploy "$id" "init" "{\"beneficiary\":\"$NEAR\"}"
    tx_call "$id" "donate" '{}' '0.1 NEAR'

    local result; result=$(view_call "$id" "get_donation_for_account" "{\"account_id\":\"$NEAR\"}")
    echo "donation record → $result"
    echo "$result" | grep -q "account_id" && echo "✓ donation passed"
}

test_auction() {
    local id="auction-$(unique_id).$NEAR"
    echo "Contract: $id"
    cd "$ROOT/auction"
    "$NEAR_GO" build

    near account create-account fund-myself "$id" '1 NEAR' \
        autogenerate-new-keypair save-to-keychain \
        sign-as "$NEAR" network-config testnet sign-with-keychain send

    # end_time = now + 10 minutes (in milliseconds)
    local end_time; end_time=$(( $(date +%s) * 1000 + 600000 ))
    deploy "$id" "init" "{\"end_time\":$end_time,\"auctioneer\":\"$NEAR\"}"

    tx_call "$id" "bid" '{}' '0.1 NEAR'

    local bid; bid=$(view_call "$id" "get_highest_bid" '{}')
    echo "highest_bid → $bid"
    echo "$bid" | grep -q "bidder" && echo "✓ auction passed"
}

# ── Main ──────────────────────────────────────────────────────────────────────

TARGETS=("${@:-greeting counter donation auction}")

for t in "${TARGETS[@]}"; do
    run "$t"
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━  SUMMARY  ━━━━━━━━━━━━━━━━━━━━━━━━"
for p in "${PASS[@]}"; do echo "  ✓ $p"; done
for p in "${FAIL[@]}"; do echo "  ✗ $p"; done
echo ""
[ ${#FAIL[@]} -eq 0 ] && echo "All passed." && exit 0 || { echo "${#FAIL[@]} failed."; exit 1; }

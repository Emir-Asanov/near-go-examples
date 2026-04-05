# near-go-examples

Complete, working smart contract examples for [NEAR Protocol](https://near.org) written in Go using [near-sdk-go](https://github.com/vlmoon99/near-sdk-go).

These examples accompany the [NEAR Go documentation](https://docs.near.org/smart-contracts/).

---

## Examples

| Example | Description | Docs |
|---------|-------------|------|
| [greeting](./greeting/) | Hello World — basic get/set greeting | [Quickstart](https://docs.near.org/smart-contracts/quickstart) |
| [counter](./counter/) | Simple counter with increment/decrement/reset | [Unit Testing](https://docs.near.org/smart-contracts/testing/unit-test) |
| [donation](./donation/) | Accept donations and forward to a beneficiary | [Unit Testing](https://docs.near.org/smart-contracts/testing/unit-test) |
| [auction](./auction/) | Time-limited auction with bid and claim | [Quickstart](https://docs.near.org/smart-contracts/quickstart) |
| [cross-contract](./cross-contract/) | Query and call external contracts | [Cross-Contract Calls](https://docs.near.org/smart-contracts/anatomy/crosscontract) |
| [yield-resume](./yield-resume/) | Async oracle/AI assistant pattern | [Yield & Resume](https://docs.near.org/smart-contracts/anatomy/yield-resume) |
| [guest-book/base](./guest-book/base/) | Guest book v1 (base state structure) | [Upgrading Contracts](https://docs.near.org/smart-contracts/release/upgrade) |
| [guest-book/v2](./guest-book/v2/) | Guest book v2 with state migration | [Upgrading Contracts](https://docs.near.org/smart-contracts/release/upgrade) |
| [self-update](./self-update/) | Contract that can deploy its own update | [Upgrading Contracts](https://docs.near.org/smart-contracts/release/upgrade) |

---

## Prerequisites

### 1. Install near-go CLI

The `near-go` CLI is required for building and testing Go smart contracts. It wraps TinyGo and handles code generation from comment directives.

```bash
curl -LO https://github.com/vlmoon99/near-cli-go/releases/latest/download/install.sh && bash install.sh
```

Verify the installation:

```bash
near-go version
```

> **Note:** The installer uses [GVM (Go Version Manager)](https://github.com/moovweb/gvm) internally. If you see `GVM_ROOT not set` errors, run `source ~/.gvm/scripts/gvm` and try again.

### 2. Install NEAR CLI (for testnet)

```bash
curl --proto '=https' --tlsv1.2 -LsSf https://github.com/near/near-cli-rs/releases/latest/download/near-cli-rs-installer.sh | sh
```

Verify:

```bash
near --version
```

### 3. Install Rust (for integration tests only)

```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source "$HOME/.cargo/env"
```

---

## Unit Tests

Unit tests use `MockSystem` — a pure in-memory mock of the NEAR runtime. No network or wallet needed.

> **Important:** Go contracts use TinyGo-specific `//go:wasmimport` declarations and **cannot** be run with standard `go test`. Always use `near-go test`.

```bash
# Run unit tests for a single contract
cd greeting
near-go test project

# Run for all contracts at once (from repo root)
for dir in greeting counter donation auction self-update guest-book/base guest-book/v2 yield-resume; do
  echo "=== $dir ==="
  (cd "$dir" && near-go test project)
done
```

Expected output for each:

```
🧪 Running project tests...
✅ Tests passed!
```

---

## Integration Tests

Integration tests use [near-workspaces-rs](https://github.com/near/near-workspaces-rs) to spin up a local NEAR sandbox and deploy the contract into it. **No testnet account needed** — everything runs locally.

Available for: `greeting`, `counter`, `donation`, `auction`.

### Step 1 — Build the contract

```bash
cd greeting          # or counter / donation / auction
near-go build
# produces: main.wasm
```

### Step 2 — Run integration tests

```bash
cd integration-tests
cargo test
```

Full example for the auction contract:

```bash
cd auction
near-go build
cd integration-tests
cargo test -- --nocapture
```

All four contracts follow the same pattern.

---

## Deploying to Testnet

### Step 1 — Create a testnet account

Go to [https://testnet.mynearwallet.com](https://testnet.mynearwallet.com) and create an account (e.g. `myname.testnet`). You will receive 10 NEAR for free.

### Step 2 — Log in with NEAR CLI

```bash
near login
```

This opens a browser window. Authorize the CLI, then return to the terminal.

### Step 3 — Build the contract

```bash
cd greeting          # or any other example
near-go build
# produces: main.wasm
```

### Step 4 — Deploy

```bash
near contract deploy myname.testnet \
  use-file main.wasm \
  with-init-call init json-args '{}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' \
  network-config testnet \
  sign-with-keychain \
  send
```

> For contracts that require init arguments (auction, donation), adjust `json-args` accordingly — see the per-contract sections below.

---

## Per-Contract Testnet Guide

### greeting

```bash
cd greeting && near-go build

# Deploy with init
near contract deploy myname.testnet \
  use-file main.wasm \
  with-init-call init json-args '{}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' \
  network-config testnet sign-with-keychain send

# Call: set greeting
near contract call-function as-transaction myname.testnet set_greeting \
  json-args '{"greeting": "Hello, NEAR!"}' \
  prepaid-gas '30.0 Tgas' attached-deposit '0 NEAR' \
  sign-as myname.testnet \
  network-config testnet sign-with-keychain send

# View: get greeting
near contract call-function as-read-only myname.testnet get_greeting \
  json-args '{}' network-config testnet now
```

---

### counter

```bash
cd counter && near-go build

near contract deploy myname.testnet \
  use-file main.wasm \
  with-init-call init json-args '{}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' \
  network-config testnet sign-with-keychain send

# Increment
near contract call-function as-transaction myname.testnet increment \
  json-args '{}' prepaid-gas '30.0 Tgas' attached-deposit '0 NEAR' \
  sign-as myname.testnet network-config testnet sign-with-keychain send

# Get count
near contract call-function as-read-only myname.testnet get_count \
  json-args '{}' network-config testnet now
```

---

### donation

```bash
cd donation && near-go build

# Init requires a beneficiary account
near contract deploy myname.testnet \
  use-file main.wasm \
  with-init-call init json-args '{"beneficiary": "beneficiary.testnet"}' \
  prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' \
  network-config testnet sign-with-keychain send

# Donate 1 NEAR
near contract call-function as-transaction myname.testnet donate \
  json-args '{}' prepaid-gas '30.0 Tgas' attached-deposit '1 NEAR' \
  sign-as myname.testnet network-config testnet sign-with-keychain send

# View donations
near contract call-function as-read-only myname.testnet get_donations \
  json-args '{}' network-config testnet now
```

---

### auction

```bash
cd auction && near-go build

# end_time is a Unix timestamp in milliseconds; set it ~5 min from now
END_TIME=$(( $(date +%s%3N) + 300000 ))

near contract deploy myname.testnet \
  use-file main.wasm \
  with-init-call init json-args "{\"end_time\": $END_TIME, \"auctioneer\": \"myname.testnet\"}" \
  prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' \
  network-config testnet sign-with-keychain send

# Place a bid of 2 NEAR
near contract call-function as-transaction myname.testnet bid \
  json-args '{}' prepaid-gas '30.0 Tgas' attached-deposit '2 NEAR' \
  sign-as myname.testnet network-config testnet sign-with-keychain send

# View highest bid
near contract call-function as-read-only myname.testnet get_highest_bid \
  json-args '{}' network-config testnet now
```

---

### yield-resume

```bash
cd yield-resume && near-go build

near contract deploy myname.testnet \
  use-file main.wasm \
  with-init-call init json-args '{}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' \
  network-config testnet sign-with-keychain send

# Ask the assistant a question
near contract call-function as-transaction myname.testnet ask_assistant \
  json-args '{"prompt": "What is 2+2?"}' \
  prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' \
  sign-as myname.testnet network-config testnet sign-with-keychain send

# Check pending requests (note the request_id from the output above)
near contract call-function as-read-only myname.testnet get_pending_requests \
  json-args '{}' network-config testnet now

# Respond as external service (use the request_id from ask_assistant output)
near contract call-function as-transaction myname.testnet respond \
  json-args '{"request_id": 0, "response": "The answer is 4"}' \
  prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' \
  sign-as myname.testnet network-config testnet sign-with-keychain send
```

---

### guest-book (base → v2 migration)

This example demonstrates **state migration** between two contract versions.

```bash
# 1. Deploy base version
cd guest-book/base && near-go build
near contract deploy myname.testnet \
  use-file main.wasm \
  with-init-call init json-args '{}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' \
  network-config testnet sign-with-keychain send

# 2. Add some messages
near contract call-function as-transaction myname.testnet add_message \
  json-args '{"text": "Hello from v1!"}' \
  prepaid-gas '30.0 Tgas' attached-deposit '0 NEAR' \
  sign-as myname.testnet network-config testnet sign-with-keychain send

# 3. Build and deploy v2 (calls migrate() as init)
cd ../v2 && near-go build
near contract deploy myname.testnet \
  use-file main.wasm \
  with-init-call migrate json-args '{}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' \
  network-config testnet sign-with-keychain send

# 4. Verify messages survived migration
near contract call-function as-read-only myname.testnet get_messages \
  json-args '{}' network-config testnet now
```

---

### self-update

```bash
cd self-update && near-go build

near contract deploy myname.testnet \
  use-file main.wasm \
  with-init-call init json-args '{}' prepaid-gas '100.0 Tgas' attached-deposit '0 NEAR' \
  network-config testnet sign-with-keychain send

# The contract can redeploy itself by calling update_contract with wasm bytes.
# This is typically done programmatically; see main.go for details.
```

---

## Comment Directives Reference

The `near-go` CLI reads these directives to generate WASM exports and state management:

| Directive | Description |
|-----------|-------------|
| `// @contract:state` | Main state struct (one per project) |
| `// @contract:init` | Initialization method |
| `// @contract:view` | Read-only method — state is NOT saved after execution |
| `// @contract:mutating` | State-modifying method — state IS saved after execution |
| `// @contract:payable [min_deposit=X]` | Accepts attached NEAR tokens |
| `// @contract:promise_callback` | Callback that receives `promise.PromiseResult` as an argument |

---

## SDK Version

All examples use [near-sdk-go v0.1.1](https://github.com/vlmoon99/near-sdk-go/releases/tag/v0.1.1).

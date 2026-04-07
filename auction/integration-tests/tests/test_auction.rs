use near_workspaces::types::NearToken;
use serde_json::{json, Value};

/// Build the contract before running:
///   cd auction && near-go build
/// Then run from integration-tests/:
///   cargo test

// near-sdk-go double-encodes all return values; use two-step deserialization on every view.

fn five_minutes_from_now_ms() -> u64 {
    std::time::SystemTime::now()
        .duration_since(std::time::UNIX_EPOCH)
        .unwrap()
        .as_millis() as u64
        + 5 * 60 * 1000
}

async fn view_str(contract: &near_workspaces::Contract, method: &str) -> anyhow::Result<String> {
    let raw: String = contract.view(method).args_json(json!({})).await?.json()?;
    Ok(serde_json::from_str(&raw)?)
}

async fn view_bool(contract: &near_workspaces::Contract, method: &str) -> anyhow::Result<bool> {
    let raw: String = contract.view(method).args_json(json!({})).await?.json()?;
    Ok(serde_json::from_str(&raw)?)
}

async fn view_obj(contract: &near_workspaces::Contract, method: &str) -> anyhow::Result<Value> {
    let raw: String = contract.view(method).args_json(json!({})).await?.json()?;
    Ok(serde_json::from_str(&raw)?)
}

#[tokio::test]
async fn test_auction_init() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;
    let auctioneer = sandbox.dev_create_account().await?;

    contract
        .call("init")
        .args_json(json!({ "end_time": five_minutes_from_now_ms(), "auctioneer": auctioneer.id().to_string() }))
        .transact().await?.into_result()?;

    assert_eq!(view_str(&contract, "get_auctioneer").await?, auctioneer.id().to_string());
    assert!(!view_bool(&contract, "get_claimed").await?);
    Ok(())
}

#[tokio::test]
async fn test_auction_bid() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;
    let auctioneer = sandbox.dev_create_account().await?;

    contract
        .call("init")
        .args_json(json!({ "end_time": five_minutes_from_now_ms(), "auctioneer": auctioneer.id().to_string() }))
        .transact().await?.into_result()?;

    let bidder = sandbox.dev_create_account().await?;
    bidder
        .call(contract.id(), "bid")
        .args_json(json!({}))
        .deposit(NearToken::from_near(2))
        .transact().await?.into_result()?;

    let bid = view_obj(&contract, "get_highest_bid").await?;
    assert_eq!(bid["bidder"], bidder.id().to_string());
    Ok(())
}

#[tokio::test]
async fn test_auction_bid_must_be_higher() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;
    let auctioneer = sandbox.dev_create_account().await?;

    contract
        .call("init")
        .args_json(json!({ "end_time": five_minutes_from_now_ms(), "auctioneer": auctioneer.id().to_string() }))
        .transact().await?.into_result()?;

    let alice = sandbox.dev_create_account().await?;
    alice.call(contract.id(), "bid").args_json(json!({})).deposit(NearToken::from_near(2)).transact().await?.into_result()?;

    let bob = sandbox.dev_create_account().await?;
    let result = bob.call(contract.id(), "bid").args_json(json!({})).deposit(NearToken::from_near(2)).transact().await?;
    assert!(!result.is_success(), "Bid of same amount should fail");
    Ok(())
}

#[tokio::test]
async fn test_auction_cannot_claim_before_end() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;
    let auctioneer = sandbox.dev_create_account().await?;

    contract
        .call("init")
        .args_json(json!({ "end_time": five_minutes_from_now_ms(), "auctioneer": auctioneer.id().to_string() }))
        .transact().await?.into_result()?;

    let result = contract.call("claim").args_json(json!({})).transact().await?;
    assert!(!result.is_success(), "Claim before auction end should fail");
    Ok(())
}

#[tokio::test]
async fn test_auction_claim_after_end() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;
    let auctioneer = sandbox.dev_create_account().await?;

    // end_time=1 ms → auction already over
    contract
        .call("init")
        .args_json(json!({ "end_time": 1u64, "auctioneer": auctioneer.id().to_string() }))
        .transact().await?.into_result()?;

    contract.call("claim").args_json(json!({})).transact().await?.into_result()?;

    assert!(view_bool(&contract, "get_claimed").await?);
    Ok(())
}

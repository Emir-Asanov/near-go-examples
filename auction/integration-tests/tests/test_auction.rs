use near_workspaces::types::NearToken;
use serde_json::{json, Value};

/// Build the contract before running:
///   cd auction && near-go build
/// Then run from integration-tests/:
///   cargo test

/// Returns a timestamp 5 minutes from now in milliseconds.
fn five_minutes_from_now_ms() -> u64 {
    let now_ms = std::time::SystemTime::now()
        .duration_since(std::time::UNIX_EPOCH)
        .unwrap()
        .as_millis() as u64;
    now_ms + 5 * 60 * 1000
}

#[tokio::test]
async fn test_auction_init() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    let auctioneer = sandbox.dev_create_account().await?;
    let end_time = five_minutes_from_now_ms();

    contract
        .call("init")
        .args_json(json!({
            "end_time": end_time,
            "auctioneer": auctioneer.id().to_string()
        }))
        .transact()
        .await?
        .into_result()?;

    let auctioneer_result: String = contract
        .view("get_auctioneer")
        .args_json(json!({}))
        .await?
        .json()?;

    assert_eq!(auctioneer_result, auctioneer.id().to_string());

    let claimed: bool = contract
        .view("get_claimed")
        .args_json(json!({}))
        .await?
        .json()?;

    assert!(!claimed);
    Ok(())
}

#[tokio::test]
async fn test_auction_bid() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    let auctioneer = sandbox.dev_create_account().await?;
    let end_time = five_minutes_from_now_ms();

    contract
        .call("init")
        .args_json(json!({
            "end_time": end_time,
            "auctioneer": auctioneer.id().to_string()
        }))
        .transact()
        .await?
        .into_result()?;

    let bidder = sandbox.dev_create_account().await?;

    bidder
        .call(contract.id(), "bid")
        .args_json(json!({}))
        .deposit(NearToken::from_near(2))
        .transact()
        .await?
        .into_result()?;

    let highest_bid: Value = contract
        .view("get_highest_bid")
        .args_json(json!({}))
        .await?
        .json()?;

    assert_eq!(highest_bid["bidder"], bidder.id().to_string());
    Ok(())
}

#[tokio::test]
async fn test_auction_bid_must_be_higher() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    let auctioneer = sandbox.dev_create_account().await?;
    let end_time = five_minutes_from_now_ms();

    contract
        .call("init")
        .args_json(json!({
            "end_time": end_time,
            "auctioneer": auctioneer.id().to_string()
        }))
        .transact()
        .await?
        .into_result()?;

    let alice = sandbox.dev_create_account().await?;
    alice
        .call(contract.id(), "bid")
        .args_json(json!({}))
        .deposit(NearToken::from_near(2))
        .transact()
        .await?
        .into_result()?;

    // Bob bids the same amount — should fail
    let bob = sandbox.dev_create_account().await?;
    let result = bob
        .call(contract.id(), "bid")
        .args_json(json!({}))
        .deposit(NearToken::from_near(2))
        .transact()
        .await?;

    assert!(!result.is_success(), "Bid of same amount should fail");
    Ok(())
}

#[tokio::test]
async fn test_auction_cannot_claim_before_end() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    let auctioneer = sandbox.dev_create_account().await?;
    let end_time = five_minutes_from_now_ms();

    contract
        .call("init")
        .args_json(json!({
            "end_time": end_time,
            "auctioneer": auctioneer.id().to_string()
        }))
        .transact()
        .await?
        .into_result()?;

    let result = contract
        .call("claim")
        .args_json(json!({}))
        .transact()
        .await?;

    assert!(!result.is_success(), "Claim before auction end should fail");
    Ok(())
}

#[tokio::test]
async fn test_auction_claim_after_end() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    let auctioneer = sandbox.dev_create_account().await?;

    // Set end time in the past (1 ms since epoch) so the auction is already over
    let end_time_ms: u64 = 1;

    contract
        .call("init")
        .args_json(json!({
            "end_time": end_time_ms,
            "auctioneer": auctioneer.id().to_string()
        }))
        .transact()
        .await?
        .into_result()?;

    // Auction is already over — claim should succeed
    contract
        .call("claim")
        .args_json(json!({}))
        .transact()
        .await?
        .into_result()?;

    let claimed: bool = contract
        .view("get_claimed")
        .args_json(json!({}))
        .await?
        .json()?;

    assert!(claimed);
    Ok(())
}

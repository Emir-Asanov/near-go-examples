use near_workspaces::types::NearToken;
use serde_json::{json, Value};

/// Build the contract before running:
///   cd donation && near-go build
/// Then run from integration-tests/:
///   cargo test

// near-sdk-go double-encodes all return values; use two-step deserialization on every view.

#[tokio::test]
async fn test_donation_beneficiary() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    let beneficiary = sandbox.dev_create_account().await?;
    contract
        .call("init")
        .args_json(json!({ "beneficiary": beneficiary.id().to_string() }))
        .transact().await?.into_result()?;

    let raw: String = contract.view("get_beneficiary").args_json(json!({})).await?.json()?;
    let b: String = serde_json::from_str(&raw)?;
    assert_eq!(b, beneficiary.id().to_string());
    Ok(())
}

#[tokio::test]
async fn test_donation_records_donor() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    let beneficiary = sandbox.dev_create_account().await?;
    contract
        .call("init")
        .args_json(json!({ "beneficiary": beneficiary.id().to_string() }))
        .transact().await?.into_result()?;

    let donor = sandbox.dev_create_account().await?;
    donor
        .call(contract.id(), "donate")
        .args_json(json!({}))
        .deposit(NearToken::from_near(1))
        .transact().await?.into_result()?;

    let raw: String = contract.view("get_donations").args_json(json!({})).await?.json()?;
    let donations: Vec<Value> = serde_json::from_str(&raw)?;

    assert_eq!(donations.len(), 1);
    assert_eq!(donations[0]["donor"], donor.id().to_string());
    Ok(())
}

#[tokio::test]
async fn test_donation_multiple_donors() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    let beneficiary = sandbox.dev_create_account().await?;
    contract
        .call("init")
        .args_json(json!({ "beneficiary": beneficiary.id().to_string() }))
        .transact().await?.into_result()?;

    for _ in 0..3 {
        let donor = sandbox.dev_create_account().await?;
        donor
            .call(contract.id(), "donate")
            .args_json(json!({}))
            .deposit(NearToken::from_millinear(500))
            .transact().await?.into_result()?;
    }

    let raw: String = contract.view("get_donations").args_json(json!({})).await?.json()?;
    let donations: Vec<Value> = serde_json::from_str(&raw)?;
    assert_eq!(donations.len(), 3);
    Ok(())
}

use near_workspaces::types::NearToken;
use serde_json::{json, Value};

/// Build the contract before running:
///   cd donation && near-go build
/// Then run from integration-tests/:
///   cargo test

async fn deploy_and_init(
    beneficiary_id: &str,
) -> anyhow::Result<(
    near_workspaces::Worker<near_workspaces::network::Sandbox>,
    near_workspaces::Contract,
)> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    contract
        .call("init")
        .args_json(json!({ "beneficiary": beneficiary_id }))
        .transact()
        .await?
        .into_result()?;

    Ok((sandbox, contract))
}

#[tokio::test]
async fn test_donation_beneficiary() -> anyhow::Result<()> {
    let beneficiary = sandbox_account_id();
    let (_sandbox, contract) = deploy_and_init(&beneficiary).await?;

    let b: String = contract
        .view("get_beneficiary")
        .args_json(json!({}))
        .await?
        .json()?;

    assert_eq!(b, beneficiary);
    Ok(())
}

// Note: beneficiary must be a valid NEAR account on the sandbox.
// We use the contract's own account for simplicity.
fn sandbox_account_id() -> String {
    "beneficiary.sandbox".to_string()
}

#[tokio::test]
async fn test_donation_records_donor() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    // Create beneficiary account first
    let beneficiary = sandbox.dev_create_account().await?;

    contract
        .call("init")
        .args_json(json!({ "beneficiary": beneficiary.id().to_string() }))
        .transact()
        .await?
        .into_result()?;

    let donor = sandbox.dev_create_account().await?;

    donor
        .call(contract.id(), "donate")
        .args_json(json!({}))
        .deposit(NearToken::from_near(1))
        .transact()
        .await?
        .into_result()?;

    let donations: Vec<Value> = contract
        .view("get_donations")
        .args_json(json!({}))
        .await?
        .json()?;

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
        .transact()
        .await?
        .into_result()?;

    for _ in 0..3 {
        let donor = sandbox.dev_create_account().await?;
        donor
            .call(contract.id(), "donate")
            .args_json(json!({}))
            .deposit(NearToken::from_millinear(500))
            .transact()
            .await?
            .into_result()?;
    }

    let donations: Vec<Value> = contract
        .view("get_donations")
        .args_json(json!({}))
        .await?
        .json()?;

    assert_eq!(donations.len(), 3);
    Ok(())
}

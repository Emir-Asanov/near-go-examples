use base64::{engine::general_purpose::STANDARD as BASE64, Engine as _};
use serde_json::json;

/// Build the contract before running:
///   cd self-update && near-go build
/// Then run from integration-tests/:
///   cargo test

// near-sdk-go double-encodes all return values; use two-step deserialization on every view.

#[tokio::test]
async fn test_self_update_init() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;
    let owner = sandbox.dev_create_account().await?;

    contract
        .call("init")
        .args_json(json!({ "ownerId": owner.id().to_string() }))
        .transact().await?.into_result()?;

    let raw: String = contract.view("get_owner").args_json(json!({})).await?.json()?;
    let returned_owner: String = serde_json::from_str(&raw)?;
    assert_eq!(returned_owner, owner.id().to_string());
    Ok(())
}

#[tokio::test]
async fn test_non_owner_cannot_update() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;
    let owner = sandbox.dev_create_account().await?;
    let attacker = sandbox.dev_create_account().await?;

    contract
        .call("init")
        .args_json(json!({ "ownerId": owner.id().to_string() }))
        .transact().await?.into_result()?;

    let result = attacker
        .call(contract.id(), "update_contract")
        .args_json(json!({ "wasm_bytes": BASE64.encode(&[] as &[u8]) }))
        .transact().await?;

    assert!(!result.is_success(), "Non-owner should not be able to update contract");
    Ok(())
}

#[tokio::test]
async fn test_owner_can_update() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;
    let owner = sandbox.dev_create_account().await?;

    contract
        .call("init")
        .args_json(json!({ "ownerId": owner.id().to_string() }))
        .transact().await?.into_result()?;

    // Use a minimal WASM (just magic bytes) for the update test.
    // Deploying the full ~178KB WASM via a JSON arg would require fitting it
    // in the function arg limit; a minimal valid WASM magic header suffices
    // for testing the access-control logic.
    let minimal_wasm: &[u8] = &[0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00];
    let result = owner
        .call(contract.id(), "update_contract")
        .args_json(json!({ "wasm_bytes": BASE64.encode(minimal_wasm) }))
        .gas(near_workspaces::types::Gas::from_tgas(300))
        .transact().await?;

    assert!(result.is_success(), "Owner should be able to update contract: {:?}", result.failures());
    Ok(())
}

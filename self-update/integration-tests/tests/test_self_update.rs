use serde_json::json;

/// Build the contract before running:
///   cd self-update && near-go build
/// Then run from integration-tests/:
///   cargo test

#[tokio::test]
async fn test_self_update_init() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;
    let owner = sandbox.dev_create_account().await?;

    contract
        .call("init")
        .args_json(json!({ "owner_id": owner.id().to_string() }))
        .transact()
        .await?
        .into_result()?;

    let returned_owner: String = contract
        .view("get_owner")
        .args_json(json!({}))
        .await?
        .json()?;

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
        .args_json(json!({ "owner_id": owner.id().to_string() }))
        .transact()
        .await?
        .into_result()?;

    // Attacker tries to deploy a new version — should fail
    let result = attacker
        .call(contract.id(), "update_contract")
        .args_json(json!({ "wasm_bytes": [] }))
        .transact()
        .await?;

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
        .args_json(json!({ "owner_id": owner.id().to_string() }))
        .transact()
        .await?
        .into_result()?;

    // Owner redeploys the same wasm — this is the self-update mechanism
    let result = owner
        .call(contract.id(), "update_contract")
        .args_json(json!({ "wasm_bytes": wasm }))
        .gas(near_workspaces::types::Gas::from_tgas(300))
        .transact()
        .await?;

    assert!(result.is_success(), "Owner should be able to update contract");
    Ok(())
}

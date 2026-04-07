use near_workspaces::types::NearToken;
use serde_json::{json, Value};

/// Build both contracts before running:
///   cd guest-book/base && near-go build
///   cd guest-book/v2  && near-go build
/// The test.sh script copies base/main.wasm → v2/integration-tests/base.wasm automatically.
/// Then run from guest-book/v2/integration-tests/:
///   cargo test

// near-sdk-go double-encodes all return values; use two-step deserialization on every view.

async fn get_messages(contract: &near_workspaces::Contract) -> anyhow::Result<Vec<Value>> {
    let raw: String = contract.view("get_messages").args_json(json!({})).await?.json()?;
    Ok(serde_json::from_str(&raw)?)
}

#[tokio::test]
async fn test_v2_init() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;
    contract.call("init").args_json(json!({})).transact().await?.into_result()?;

    assert!(get_messages(&contract).await?.is_empty());
    Ok(())
}

#[tokio::test]
async fn test_v2_add_message_has_payment_field() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;
    contract.call("init").args_json(json!({})).transact().await?.into_result()?;

    let alice = sandbox.dev_create_account().await?;
    alice
        .call(contract.id(), "add_message")
        .args_json(json!({ "text": "Hello v2!" }))
        .deposit(NearToken::from_near(1))
        .transact().await?.into_result()?;

    let messages = get_messages(&contract).await?;
    assert_eq!(messages.len(), 1);
    assert_eq!(messages[0]["text"], "Hello v2!");
    assert!(messages[0]["payment"].is_string(), "Expected payment field in v2 message");
    Ok(())
}

/// Tests the full upgrade path:
/// 1. Deploy base contract, add messages with payments
/// 2. Redeploy v2 and call migrate()
/// 3. Verify messages survived with payments embedded
#[tokio::test]
async fn test_migration_from_base_to_v2() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let base_wasm = std::fs::read("base.wasm")?;
    let v2_wasm = std::fs::read("../main.wasm")?;

    // Step 1: Deploy base
    let contract = sandbox.dev_deploy(&base_wasm).await?;
    contract.call("init").args_json(json!({})).transact().await?.into_result()?;

    // Step 2: Add a message to base
    let alice = sandbox.dev_create_account().await?;
    alice
        .call(contract.id(), "add_message")
        .args_json(json!({ "text": "From v1" }))
        .deposit(NearToken::from_near(1))
        .transact().await?.into_result()?;

    let raw_base: String = contract.view("get_messages").args_json(json!({})).await?.json()?;
    let base_messages: Vec<Value> = serde_json::from_str(&raw_base)?;
    assert_eq!(base_messages.len(), 1);

    // Step 3: Upgrade to v2 and migrate
    contract.as_account().deploy(&v2_wasm).await?.into_result()?;
    contract.call("migrate").args_json(json!({})).transact().await?.into_result()?;

    // Step 4: Verify messages survived
    let v2_messages = get_messages(&contract).await?;
    assert_eq!(v2_messages.len(), 1, "Messages should survive migration");
    assert_eq!(v2_messages[0]["text"], "From v1");
    assert!(v2_messages[0]["payment"].is_string(), "Payment should be embedded after migration");
    Ok(())
}

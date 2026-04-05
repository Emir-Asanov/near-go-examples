use near_workspaces::types::NearToken;
use serde_json::{json, Value};

/// Tests for guest-book v2, including the state migration path.
///
/// Build both contracts before running:
///   cd guest-book/base && near-go build && cp main.wasm ../v2/integration-tests/base.wasm
///   cd guest-book/v2  && near-go build
/// Then run from guest-book/v2/integration-tests/:
///   cargo test

#[tokio::test]
async fn test_v2_init() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    contract
        .call("init")
        .args_json(json!({}))
        .transact()
        .await?
        .into_result()?;

    let messages: Vec<Value> = contract
        .view("get_messages")
        .args_json(json!({}))
        .await?
        .json()?;

    assert!(messages.is_empty());
    Ok(())
}

#[tokio::test]
async fn test_v2_add_message_has_payment_field() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    contract
        .call("init")
        .args_json(json!({}))
        .transact()
        .await?
        .into_result()?;

    let alice = sandbox.dev_create_account().await?;
    alice
        .call(contract.id(), "add_message")
        .args_json(json!({ "text": "Hello v2!" }))
        .deposit(NearToken::from_near(1))
        .transact()
        .await?
        .into_result()?;

    let messages: Vec<Value> = contract
        .view("get_messages")
        .args_json(json!({}))
        .await?
        .json()?;

    assert_eq!(messages.len(), 1);
    assert_eq!(messages[0]["text"], "Hello v2!");
    // v2 embeds payment directly in the message
    assert!(messages[0]["payment"].is_string(), "Expected payment field in v2 message");
    Ok(())
}

/// Tests the full upgrade path:
/// 1. Deploy base contract, add messages
/// 2. Redeploy v2 calling migrate() as init
/// 3. Verify messages survived migration with payment embedded
#[tokio::test]
async fn test_migration_from_base_to_v2() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let base_wasm = std::fs::read("base.wasm")?;
    let v2_wasm = std::fs::read("../main.wasm")?;

    // Step 1: Deploy base contract
    let contract = sandbox.dev_deploy(&base_wasm).await?;
    contract
        .call("init")
        .args_json(json!({}))
        .transact()
        .await?
        .into_result()?;

    // Step 2: Add messages to base
    let alice = sandbox.dev_create_account().await?;
    alice
        .call(contract.id(), "add_message")
        .args_json(json!({ "text": "From v1" }))
        .deposit(NearToken::from_near(1))
        .transact()
        .await?
        .into_result()?;

    let base_messages: Vec<Value> = contract
        .view("get_messages")
        .args_json(json!({}))
        .await?
        .json()?;
    assert_eq!(base_messages.len(), 1);

    // Step 3: Redeploy as v2 with migrate()
    contract
        .as_account()
        .deploy(&v2_wasm)
        .await?
        .into_result()?;

    contract
        .call("migrate")
        .args_json(json!({}))
        .transact()
        .await?
        .into_result()?;

    // Step 4: Verify messages survived
    let v2_messages: Vec<Value> = contract
        .view("get_messages")
        .args_json(json!({}))
        .await?
        .json()?;

    assert_eq!(v2_messages.len(), 1, "Messages should survive migration");
    assert_eq!(v2_messages[0]["text"], "From v1");
    assert!(v2_messages[0]["payment"].is_string(), "Payment should be embedded after migration");
    Ok(())
}

use near_workspaces::types::NearToken;
use serde_json::{json, Value};

/// Build the contract before running:
///   cd guest-book/base && near-go build
/// Then run from integration-tests/:
///   cargo test

#[tokio::test]
async fn test_guest_book_init() -> anyhow::Result<()> {
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

    assert!(messages.is_empty(), "Expected no messages after init");
    Ok(())
}

#[tokio::test]
async fn test_add_regular_message() -> anyhow::Result<()> {
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

    // Small deposit — not premium
    alice
        .call(contract.id(), "add_message")
        .args_json(json!({ "text": "Hello, NEAR!" }))
        .deposit(NearToken::from_millinear(1)) // 0.001 NEAR
        .transact()
        .await?
        .into_result()?;

    let messages: Vec<Value> = contract
        .view("get_messages")
        .args_json(json!({}))
        .await?
        .json()?;

    assert_eq!(messages.len(), 1);
    assert_eq!(messages[0]["text"], "Hello, NEAR!");
    assert_eq!(messages[0]["sender"], alice.id().to_string());
    assert_eq!(messages[0]["premium"], false);
    Ok(())
}

#[tokio::test]
async fn test_add_premium_message() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    contract
        .call("init")
        .args_json(json!({}))
        .transact()
        .await?
        .into_result()?;

    let bob = sandbox.dev_create_account().await?;

    // 0.01 NEAR — premium threshold
    bob
        .call(contract.id(), "add_message")
        .args_json(json!({ "text": "Premium message!" }))
        .deposit(NearToken::from_millinear(10)) // 0.01 NEAR
        .transact()
        .await?
        .into_result()?;

    let messages: Vec<Value> = contract
        .view("get_messages")
        .args_json(json!({}))
        .await?
        .json()?;

    assert_eq!(messages[0]["premium"], true);
    Ok(())
}

#[tokio::test]
async fn test_multiple_messages() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    contract
        .call("init")
        .args_json(json!({}))
        .transact()
        .await?
        .into_result()?;

    for i in 0..3 {
        let user = sandbox.dev_create_account().await?;
        user
            .call(contract.id(), "add_message")
            .args_json(json!({ "text": format!("Message {}", i) }))
            .deposit(NearToken::from_millinear(1))
            .transact()
            .await?
            .into_result()?;
    }

    let messages: Vec<Value> = contract
        .view("get_messages")
        .args_json(json!({}))
        .await?
        .json()?;

    assert_eq!(messages.len(), 3);
    Ok(())
}

#[tokio::test]
async fn test_payments_recorded() -> anyhow::Result<()> {
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
        .args_json(json!({ "text": "Paid" }))
        .deposit(NearToken::from_near(1))
        .transact()
        .await?
        .into_result()?;

    let payments: Vec<String> = contract
        .view("get_payments")
        .args_json(json!({}))
        .await?
        .json()?;

    assert_eq!(payments.len(), 1);
    assert_eq!(payments[0], "1000000000000000000000000");
    Ok(())
}

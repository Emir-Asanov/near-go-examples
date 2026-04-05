use near_workspaces::types::NearToken;
use serde_json::json;

/// Build the contract before running:
///   cd greeting && near-go build
/// Then run from integration-tests/:
///   cargo test

#[tokio::test]
async fn test_greeting_init() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    // Initialize the contract
    contract
        .call("init")
        .args_json(json!({}))
        .transact()
        .await?
        .into_result()?;

    // Default greeting should be "Hello"
    let greeting: String = contract
        .view("get_greeting")
        .args_json(json!({}))
        .await?
        .json()?;

    assert_eq!(greeting, "Hello");
    Ok(())
}

#[tokio::test]
async fn test_set_greeting() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    contract
        .call("init")
        .args_json(json!({}))
        .transact()
        .await?
        .into_result()?;

    // Create a caller account
    let caller = sandbox.dev_create_account().await?;

    // Set a new greeting
    caller
        .call(contract.id(), "set_greeting")
        .args_json(json!({ "greeting": "Howdy" }))
        .transact()
        .await?
        .into_result()?;

    let greeting: String = contract
        .view("get_greeting")
        .args_json(json!({}))
        .await?
        .json()?;

    assert_eq!(greeting, "Howdy");
    Ok(())
}

#[tokio::test]
async fn test_greeting_update_multiple_times() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    contract
        .call("init")
        .args_json(json!({}))
        .transact()
        .await?
        .into_result()?;

    let caller = sandbox.dev_create_account().await?;

    for msg in &["Hi", "Hello", "Привет"] {
        caller
            .call(contract.id(), "set_greeting")
            .args_json(json!({ "greeting": msg }))
            .transact()
            .await?
            .into_result()?;
    }

    let greeting: String = contract
        .view("get_greeting")
        .args_json(json!({}))
        .await?
        .json()?;

    assert_eq!(greeting, "Привет");
    Ok(())
}

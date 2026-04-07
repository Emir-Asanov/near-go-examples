use serde_json::json;

/// Build the contract before running:
///   cd greeting && near-go build
/// Then run from integration-tests/:
///   cargo test

// near-sdk-go double-encodes all return values (ReturnValue calls json.Marshal on an
// already-JSON string). Every view result arrives as a JSON string that wraps the real
// value, so we always deserialize in two steps:
//   1.  .json::<String>()         → unwrap the outer encoding
//   2.  serde_json::from_str(…)   → parse the actual type

#[tokio::test]
async fn test_greeting_init() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    contract
        .call("init")
        .args_json(json!({}))
        .transact()
        .await?
        .into_result()?;

    let raw: String = contract.view("get_greeting").args_json(json!({})).await?.json()?;
    let greeting: String = serde_json::from_str(&raw)?;
    assert_eq!(greeting, "Hello");
    Ok(())
}

#[tokio::test]
async fn test_set_greeting() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    contract.call("init").args_json(json!({})).transact().await?.into_result()?;

    let caller = sandbox.dev_create_account().await?;
    caller
        .call(contract.id(), "set_greeting")
        .args_json(json!({ "greeting": "Howdy" }))
        .transact()
        .await?
        .into_result()?;

    let raw: String = contract.view("get_greeting").args_json(json!({})).await?.json()?;
    let greeting: String = serde_json::from_str(&raw)?;
    assert_eq!(greeting, "Howdy");
    Ok(())
}

#[tokio::test]
async fn test_greeting_update_multiple_times() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    contract.call("init").args_json(json!({})).transact().await?.into_result()?;

    let caller = sandbox.dev_create_account().await?;
    for msg in &["Hi", "Hello", "Привет"] {
        caller
            .call(contract.id(), "set_greeting")
            .args_json(json!({ "greeting": msg }))
            .transact()
            .await?
            .into_result()?;
    }

    let raw: String = contract.view("get_greeting").args_json(json!({})).await?.json()?;
    let greeting: String = serde_json::from_str(&raw)?;
    assert_eq!(greeting, "Привет");
    Ok(())
}

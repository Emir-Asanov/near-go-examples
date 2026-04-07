use serde_json::{json, Value};

/// Build the contract before running:
///   cd collections && near-go build
/// Then run from integration-tests/:
///   cargo test

// near-sdk-go double-encodes all return values; use two-step deserialization on every view.

#[tokio::test]
async fn test_collections_init() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;
    contract.call("init").args_json(json!({})).transact().await?.into_result()?;

    let raw: String = contract.view("vector_get_all").args_json(json!({})).await?.json()?;
    let all: Vec<String> = serde_json::from_str(&raw)?;
    assert!(all.is_empty());
    Ok(())
}

#[tokio::test]
async fn test_vector_push_and_get() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;
    contract.call("init").args_json(json!({})).transact().await?.into_result()?;

    let caller = sandbox.dev_create_account().await?;
    caller.call(contract.id(), "vector_push").args_json(json!({ "value": "hello" })).transact().await?.into_result()?;
    caller.call(contract.id(), "vector_push").args_json(json!({ "value": "world" })).transact().await?.into_result()?;

    let raw_item: String = contract.view("vector_get").args_json(json!({ "index": 0 })).await?.json()?;
    let item: String = serde_json::from_str(&raw_item)?;
    assert_eq!(item, "hello");

    let raw_all: String = contract.view("vector_get_all").args_json(json!({})).await?.json()?;
    let all: Vec<String> = serde_json::from_str(&raw_all)?;
    assert_eq!(all.len(), 2);
    Ok(())
}

#[tokio::test]
async fn test_lookup_map_set_and_get() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;
    contract.call("init").args_json(json!({})).transact().await?.into_result()?;

    let caller = sandbox.dev_create_account().await?;
    caller.call(contract.id(), "lookup_map_set").args_json(json!({ "key": "alice", "value": "100" })).transact().await?.into_result()?;

    let raw_val: String = contract.view("lookup_map_get").args_json(json!({ "key": "alice" })).await?.json()?;
    let val: String = serde_json::from_str(&raw_val)?;
    assert_eq!(val, "100");

    let raw_contains: String = contract.view("lookup_map_contains").args_json(json!({ "key": "alice" })).await?.json()?;
    let contains: bool = serde_json::from_str(&raw_contains)?;
    assert!(contains);
    Ok(())
}

#[tokio::test]
async fn test_unordered_set_add_and_contains() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;
    contract.call("init").args_json(json!({})).transact().await?.into_result()?;

    let caller = sandbox.dev_create_account().await?;
    caller.call(contract.id(), "unordered_set_add").args_json(json!({ "value": "owner1.testnet" })).transact().await?.into_result()?;

    let raw_c: String = contract.view("unordered_set_contains").args_json(json!({ "value": "owner1.testnet" })).await?.json()?;
    let contains: bool = serde_json::from_str(&raw_c)?;
    assert!(contains);

    let raw_nc: String = contract.view("unordered_set_contains").args_json(json!({ "value": "unknown.testnet" })).await?.json()?;
    let not_contains: bool = serde_json::from_str(&raw_nc)?;
    assert!(!not_contains);
    Ok(())
}

#[tokio::test]
async fn test_pagination() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;
    contract.call("init").args_json(json!({})).transact().await?.into_result()?;

    let caller = sandbox.dev_create_account().await?;
    for i in 0..5 {
        caller.call(contract.id(), "vector_push").args_json(json!({ "value": format!("item{}", i) })).transact().await?.into_result()?;
    }

    let raw_page: String = contract.view("get_page").args_json(json!({ "fromIndex": 1, "limit": 3 })).await?.json()?;
    let page: Vec<Value> = serde_json::from_str(&raw_page)?;
    assert_eq!(page.len(), 3);

    let raw_last: String = contract.view("get_page").args_json(json!({ "fromIndex": 4, "limit": 10 })).await?.json()?;
    let last_page: Vec<Value> = serde_json::from_str(&raw_last)?;
    assert_eq!(last_page.len(), 1);
    Ok(())
}

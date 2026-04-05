use serde_json::{json, Value};

/// Build the contract before running:
///   cd collections && near-go build
/// Then run from integration-tests/:
///   cargo test

#[tokio::test]
async fn test_collections_init() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    contract
        .call("init")
        .args_json(json!({}))
        .transact()
        .await?
        .into_result()?;

    // After init all collections should be empty
    let all: Vec<String> = contract
        .view("vector_get_all")
        .args_json(json!({}))
        .await?
        .json()?;
    assert!(all.is_empty());
    Ok(())
}

#[tokio::test]
async fn test_vector_push_and_get() -> anyhow::Result<()> {
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

    caller
        .call(contract.id(), "vector_push")
        .args_json(json!({ "value": "hello" }))
        .transact()
        .await?
        .into_result()?;

    caller
        .call(contract.id(), "vector_push")
        .args_json(json!({ "value": "world" }))
        .transact()
        .await?
        .into_result()?;

    let item: String = contract
        .view("vector_get")
        .args_json(json!({ "index": 0 }))
        .await?
        .json()?;
    assert_eq!(item, "hello");

    let all: Vec<String> = contract
        .view("vector_get_all")
        .args_json(json!({}))
        .await?
        .json()?;
    assert_eq!(all.len(), 2);
    Ok(())
}

#[tokio::test]
async fn test_lookup_map_set_and_get() -> anyhow::Result<()> {
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

    caller
        .call(contract.id(), "lookup_map_set")
        .args_json(json!({ "key": "alice", "value": "100" }))
        .transact()
        .await?
        .into_result()?;

    let val: String = contract
        .view("lookup_map_get")
        .args_json(json!({ "key": "alice" }))
        .await?
        .json()?;
    assert_eq!(val, "100");

    let contains: bool = contract
        .view("lookup_map_contains")
        .args_json(json!({ "key": "alice" }))
        .await?
        .json()?;
    assert!(contains);
    Ok(())
}

#[tokio::test]
async fn test_unordered_set_add_and_contains() -> anyhow::Result<()> {
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

    caller
        .call(contract.id(), "unordered_set_add")
        .args_json(json!({ "value": "owner1.testnet" }))
        .transact()
        .await?
        .into_result()?;

    let contains: bool = contract
        .view("unordered_set_contains")
        .args_json(json!({ "value": "owner1.testnet" }))
        .await?
        .json()?;
    assert!(contains);

    let not_contains: bool = contract
        .view("unordered_set_contains")
        .args_json(json!({ "value": "unknown.testnet" }))
        .await?
        .json()?;
    assert!(!not_contains);
    Ok(())
}

#[tokio::test]
async fn test_pagination() -> anyhow::Result<()> {
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

    for i in 0..5 {
        caller
            .call(contract.id(), "vector_push")
            .args_json(json!({ "value": format!("item{}", i) }))
            .transact()
            .await?
            .into_result()?;
    }

    let page: Vec<Value> = contract
        .view("get_page")
        .args_json(json!({ "from_index": 1, "limit": 3 }))
        .await?
        .json()?;
    assert_eq!(page.len(), 3);

    let last_page: Vec<Value> = contract
        .view("get_page")
        .args_json(json!({ "from_index": 4, "limit": 10 }))
        .await?
        .json()?;
    assert_eq!(last_page.len(), 1);
    Ok(())
}

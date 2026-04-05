use serde_json::{json, Value};

/// Build the contract before running:
///   cd yield-resume && near-go build
/// Then run from integration-tests/:
///   cargo test
///
/// Note: yield/resume spans multiple blocks. These tests cover the observable
/// surface (ask_assistant, get_pending_requests, respond) but cannot fully
/// simulate the async callback execution in a sandbox without fast-forwarding
/// blocks (200 blocks ≈ 2 minutes timeout).

#[tokio::test]
async fn test_yield_resume_init() -> anyhow::Result<()> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    contract
        .call("init")
        .args_json(json!({}))
        .transact()
        .await?
        .into_result()?;

    let pending: Value = contract
        .view("get_pending_requests")
        .args_json(json!({}))
        .await?
        .json()?;

    // Should be an empty object/map
    assert!(pending.as_object().map(|m| m.is_empty()).unwrap_or(true));
    Ok(())
}

#[tokio::test]
async fn test_ask_assistant_creates_request() -> anyhow::Result<()> {
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

    let result: Value = caller
        .call(contract.id(), "ask_assistant")
        .args_json(json!({ "prompt": "What is NEAR?" }))
        .gas(near_workspaces::types::Gas::from_tgas(100))
        .transact()
        .await?
        .json()?;

    assert_eq!(result["request_id"], 0);
    assert_eq!(result["status"], "processing");

    let pending: Value = contract
        .view("get_pending_requests")
        .args_json(json!({}))
        .await?
        .json()?;

    // Request 0 should now be in the pending map
    assert!(
        pending.as_object().map(|m| m.contains_key("0")).unwrap_or(false),
        "Request 0 should be in pending requests"
    );
    Ok(())
}

#[tokio::test]
async fn test_respond_to_invalid_request_fails() -> anyhow::Result<()> {
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

    // Responding to a non-existent request should fail
    let result = caller
        .call(contract.id(), "respond")
        .args_json(json!({ "request_id": 999, "response": "Hello" }))
        .transact()
        .await?;

    assert!(!result.is_success(), "Responding to invalid request_id should fail");
    Ok(())
}

#[tokio::test]
async fn test_multiple_pending_requests() -> anyhow::Result<()> {
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

    for prompt in &["Q1", "Q2", "Q3"] {
        caller
            .call(contract.id(), "ask_assistant")
            .args_json(json!({ "prompt": prompt }))
            .gas(near_workspaces::types::Gas::from_tgas(100))
            .transact()
            .await?
            .into_result()?;
    }

    let pending: Value = contract
        .view("get_pending_requests")
        .args_json(json!({}))
        .await?
        .json()?;

    let count = pending.as_object().map(|m| m.len()).unwrap_or(0);
    assert_eq!(count, 3, "Expected 3 pending requests");
    Ok(())
}

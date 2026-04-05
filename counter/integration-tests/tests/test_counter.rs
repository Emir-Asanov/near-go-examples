use serde_json::json;

/// Build the contract before running:
///   cd counter && near-go build
/// Then run from integration-tests/:
///   cargo test

async fn deploy_and_init() -> anyhow::Result<(near_workspaces::Worker<near_workspaces::network::Sandbox>, near_workspaces::Contract)> {
    let sandbox = near_workspaces::sandbox().await?;
    let wasm = std::fs::read("../main.wasm")?;
    let contract = sandbox.dev_deploy(&wasm).await?;

    contract
        .call("init")
        .args_json(json!({}))
        .transact()
        .await?
        .into_result()?;

    Ok((sandbox, contract))
}

#[tokio::test]
async fn test_counter_init() -> anyhow::Result<()> {
    let (_sandbox, contract) = deploy_and_init().await?;

    let count: i64 = contract
        .view("get_count")
        .args_json(json!({}))
        .await?
        .json()?;

    assert_eq!(count, 0);
    Ok(())
}

#[tokio::test]
async fn test_counter_increment() -> anyhow::Result<()> {
    let (sandbox, contract) = deploy_and_init().await?;
    let caller = sandbox.dev_create_account().await?;

    caller.call(contract.id(), "increment").args_json(json!({})).transact().await?.into_result()?;
    caller.call(contract.id(), "increment").args_json(json!({})).transact().await?.into_result()?;

    let count: i64 = contract.view("get_count").args_json(json!({})).await?.json()?;
    assert_eq!(count, 2);
    Ok(())
}

#[tokio::test]
async fn test_counter_decrement() -> anyhow::Result<()> {
    let (sandbox, contract) = deploy_and_init().await?;
    let caller = sandbox.dev_create_account().await?;

    caller.call(contract.id(), "increment").args_json(json!({})).transact().await?.into_result()?;
    caller.call(contract.id(), "decrement").args_json(json!({})).transact().await?.into_result()?;

    let count: i64 = contract.view("get_count").args_json(json!({})).await?.json()?;
    assert_eq!(count, 0);
    Ok(())
}

#[tokio::test]
async fn test_counter_reset() -> anyhow::Result<()> {
    let (sandbox, contract) = deploy_and_init().await?;
    let caller = sandbox.dev_create_account().await?;

    caller.call(contract.id(), "increment").args_json(json!({})).transact().await?.into_result()?;
    caller.call(contract.id(), "increment").args_json(json!({})).transact().await?.into_result()?;
    caller.call(contract.id(), "reset").args_json(json!({})).transact().await?.into_result()?;

    let count: i64 = contract.view("get_count").args_json(json!({})).await?.json()?;
    assert_eq!(count, 0);
    Ok(())
}

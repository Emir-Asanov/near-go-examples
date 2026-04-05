package main

import (
	_ "embed"
	"strconv"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/promise"
	"github.com/vlmoon99/near-sdk-go/types"
)

//go:embed hello.wasm
var contractWasm []byte

// @contract:state
type Contract struct{}

// --- Transfer NEAR ---

type TransferTokenInput struct {
	To     string `json:"to"`
	Amount string `json:"amount"`
}

// @contract:payable min_deposit=1NEAR
func (c *Contract) ExampleTransferToken(input TransferTokenInput) error {
	amount, err := types.U128FromString(input.Amount)
	if err != nil {
		return err
	}

	promise.CreateBatch(input.To).
		Transfer(amount)

	return nil
}

// --- Function Call ---

type MessageInput struct {
	Message string `json:"message"`
}

// @contract:payable min_deposit=0.00001NEAR
func (c *Contract) ExampleFunctionCall() {
	gas := uint64(types.ONE_TERA_GAS * 10)
	accountId := "hello-nearverse.testnet"
	args := map[string]string{
		"message": "howdy",
	}
	promise.NewCrossContract(accountId).
		Gas(gas).
		Call("set_greeting", args).
		Then("example_function_call_callback", args)
}

// @contract:view
// @contract:promise_callback
func (c *Contract) ExampleFunctionCallCallback(input MessageInput, result promise.PromiseResult) MessageInput {
	env.LogString("Executing callback")
	env.LogString("Input Message : " + input.Message)

	if result.Success {
		env.LogString("Cross-contract call executed successfully")
		env.LogString("Promise Result Status --> " + strconv.FormatInt(int64(result.StatusCode), 10))
		if len(result.Data) > 0 {
			env.LogString("Batch call data: " + string(result.Data))
		}
	} else {
		env.LogString("Cross-contract call failed")
	}
	return input
}

// --- Create Sub Account ---

// @contract:payable min_deposit=0.001NEAR
func (c *Contract) ExampleCreateSubaccount(prefix string) {
	currentAccountId, err := env.GetCurrentAccountId()
	if err != nil {
		env.PanicStr("Failed to get current account")
	}

	subaccountId := prefix + "." + currentAccountId

	amount, err := types.U128FromString("1000000000000000000000") //0.001Ⓝ
	if err != nil {
		env.PanicStr("Bad amount format")
	}

	promise.CreateBatch(subaccountId).
		CreateAccount().
		Transfer(amount)
}

// --- Create .testnet / .near Account ---

type CreateAccountInput struct {
	AccountId string `json:"account_id"`
	PublicKey string `json:"public_key"`
}

// @contract:payable min_deposit=0.002NEAR
func (c *Contract) ExampleCreateAccount(args CreateAccountInput) {
	amount, _ := types.U128FromString("2000000000000000000000") // 0.002 NEAR
	gas := uint64(200 * types.ONE_TERA_GAS)

	createArgs := map[string]string{
		"new_account_id": args.AccountId,
		"new_public_key": args.PublicKey,
	}

	promise.CreateBatch("testnet").
		FunctionCall("create_account", createArgs, amount, gas)
}

// --- Deploy a Contract ---

// @contract:payable min_deposit=1.1NEAR
func (c *Contract) ExampleDeployContract(prefix string) {
	currentAccountId, _ := env.GetCurrentAccountId()
	subaccountId := prefix + "." + currentAccountId
	amount, _ := types.U128FromString("1100000000000000000000000") // 1.1Ⓝ

	promise.CreateBatch(subaccountId).
		CreateAccount().
		Transfer(amount).
		DeployContract(contractWasm)
}

// --- Add Keys ---

type AddKeysInput struct {
	Prefix    string `json:"prefix"`
	PublicKey string `json:"public_key"`
}

// @contract:payable min_deposit=0.001NEAR
func (c *Contract) ExampleAddKeys(input AddKeysInput) {
	currentAccountId, _ := env.GetCurrentAccountId()
	subaccountId := input.Prefix + "." + currentAccountId
	amount, _ := types.U128FromString("1000000000000000000000") // 0.001Ⓝ

	promise.CreateBatch(subaccountId).
		CreateAccount().
		Transfer(amount).
		AddFullAccessKey([]byte(input.PublicKey), 0)
}

// --- Delete Account ---

type DeleteAccountInput struct {
	Prefix      string `json:"prefix"`
	Beneficiary string `json:"beneficiary"`
}

type SelfDeleteInput struct {
	Beneficiary string `json:"beneficiary"`
}

// @contract:payable min_deposit=0.001NEAR
func (c *Contract) ExampleCreateDeleteAccount(input DeleteAccountInput) {
	currentAccountId, _ := env.GetCurrentAccountId()
	subaccountId := input.Prefix + "." + currentAccountId
	amount, _ := types.U128FromString("1000000000000000000000") // 0.001Ⓝ

	promise.CreateBatch(subaccountId).
		CreateAccount().
		Transfer(amount).
		DeleteAccount(input.Beneficiary)
}

// @contract:mutating
func (c *Contract) ExampleSelfDeleteAccount(input SelfDeleteInput) {
	currentAccountId, _ := env.GetCurrentAccountId()

	promise.CreateBatch(currentAccountId).
		DeleteAccount(input.Beneficiary)
}

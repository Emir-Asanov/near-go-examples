package main

import (
	"strconv"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/promise"
	"github.com/vlmoon99/near-sdk-go/types"
)

// @contract:state
type CrossContractExample struct{}

// @contract:init
func (c *CrossContractExample) Init() {
	env.LogString("CrossContractExample initialized")
}

// Query another contract — result is returned to caller immediately
// @contract:payable min_deposit=0.001NEAR
func (c *CrossContractExample) QueryGreeting() {
	helloAccount := "hello-nearverse.testnet"
	gas := uint64(10 * types.ONE_TERA_GAS)

	promise.NewCrossContract(helloAccount).
		Gas(gas).
		Call("get_greeting", map[string]string{}).
		Value()
}

// Query another contract and handle the response in a callback
// @contract:payable min_deposit=0.001NEAR
func (c *CrossContractExample) QueryGreetingWithCallback() {
	helloAccount := "hello-nearverse.testnet"
	gas := uint64(10 * types.ONE_TERA_GAS)

	promise.NewCrossContract(helloAccount).
		Gas(gas).
		Call("get_greeting", map[string]string{}).
		Then("on_query_greeting_response", map[string]string{})
}

// @contract:view
// @contract:promise_callback
func (c *CrossContractExample) OnQueryGreetingResponse(result promise.PromiseResult) string {
	if !result.Success {
		env.LogString("Query failed")
		return ""
	}
	env.LogString("Received greeting: " + string(result.Data))
	return string(result.Data)
}

// Send information to another contract with a callback
// @contract:payable min_deposit=0.00001NEAR
func (c *CrossContractExample) SetGreetingOnExternal(message string) {
	helloAccount := "hello-nearverse.testnet"
	gas := uint64(30 * types.ONE_TERA_GAS)

	promise.NewCrossContract(helloAccount).
		Gas(gas).
		Call("set_greeting", map[string]string{"greeting": message}).
		Then("on_set_greeting_callback", map[string]string{})
}

// @contract:view
// @contract:promise_callback
func (c *CrossContractExample) OnSetGreetingCallback(result promise.PromiseResult) {
	if result.Success {
		env.LogString("Greeting updated successfully on external contract")
	} else {
		env.LogString("Failed to update greeting on external contract")
	}

	env.LogString("Status: " + types.IntToString(result.StatusCode))
}

// CallbackInput is context data passed through to callback functions.
type CallbackInput struct {
	Data string `json:"data"`
}

// @contract:payable min_deposit=0.00001NEAR
func (c *CrossContractExample) CallExternalContract() {
	externalAccount := "hello-nearverse.testnet"
	gas := uint64(5 * types.ONE_TERA_GAS)

	args := map[string]string{
		"greeting": "New Greeting",
	}
	callbackArgs := map[string]string{
		"data": "saved_for_callback",
	}
	promise.NewCrossContract(externalAccount).
		Gas(gas).
		Call("set_greeting", args).
		Then("on_call_external_callback", callbackArgs).
		Value()
}

// @contract:view
// @contract:promise_callback
func (c *CrossContractExample) OnCallExternalCallback(input CallbackInput, result promise.PromiseResult) {
	env.LogString("Executing callback")
	env.LogString("Callback input data: " + input.Data)

	if result.Success {
		env.LogString("Cross-contract call executed successfully")
	} else {
		env.LogString("Cross-contract call failed")
	}
}

// @contract:payable min_deposit=0.00001NEAR
func (c *CrossContractExample) BatchCallsSameContract() {
	helloAccount := "hello-nearverse.testnet"
	gas := uint64(10 * types.ONE_TERA_GAS)
	amount, _ := types.U128FromString("0")
	callbackArgs := map[string]string{
		"data": "[Greeting One, Greeting Two]",
	}

	promise.NewCrossContract(helloAccount).
		Batch().
		Gas(gas).
		FunctionCall("set_greeting", map[string]string{
			"greeting": "Greeting One",
		}, amount, gas).
		FunctionCall("another_method", map[string]string{
			"arg1": "val1",
		}, amount, gas).
		Then(helloAccount).
		FunctionCall("on_batch_calls_callback", callbackArgs, amount, gas)

	env.LogString("Batch call created successfully")
}

// @contract:view
// @contract:promise_callback
func (c *CrossContractExample) OnBatchCallsCallback(input CallbackInput, result promise.PromiseResult) {
	env.LogString("Processing batch call results")
	env.LogString("Callback input: " + input.Data)

	env.LogString("Batch call success: " + strconv.FormatBool(result.Success))
	if len(result.Data) > 0 {
		env.LogString("Batch call data: " + string(result.Data))
	}
}

// @contract:payable min_deposit=0.00001NEAR
func (c *CrossContractExample) ParallelCallsDifferentContracts() {
	contractA := "hello-nearverse.testnet"
	contractB := "child.neargopromises1.testnet"

	promiseA := promise.NewCrossContract(contractA).
		Call("get_greeting", map[string]string{})

	promiseB := promise.NewCrossContract(contractB).
		Call("set_status", map[string]string{"message": "Hello, World!"})

	promiseA.Join([]*promise.Promise{promiseB}, "on_parallel_contracts_callback", map[string]string{
		"data": contractA + "," + contractB,
	}).Value()

	env.LogString("Parallel contract calls initialized")
}

// @contract:view
// @contract:promise_callback
func (c *CrossContractExample) OnParallelContractsCallback(input CallbackInput, results []promise.PromiseResult) {
	env.LogString("Processing results from multiple contracts")
	env.LogString("Callback input: " + input.Data)

	for i, result := range results {
		env.LogString("Processing result " + types.IntToString(i))
		env.LogString("Success: " + strconv.FormatBool(result.Success))
		if len(result.Data) > 0 {
			env.LogString("Data: " + string(result.Data))
		}
	}

	env.LogString("Processed " + types.IntToString(len(results)) + " contract responses")
}

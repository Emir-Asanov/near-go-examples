package main

import (
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

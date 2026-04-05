package main

import (
	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/promise"
)

// @contract:state
type SelfUpdateContract struct {
	OwnerId string `json:"owner_id"`
}

// @contract:init
func (c *SelfUpdateContract) Init(ownerId string) {
	c.OwnerId = ownerId
	env.LogString("SelfUpdateContract initialized. Owner: " + ownerId)
}

// UpdateContract deploys new WASM bytecode to this contract account.
// Only the owner is allowed to call this method.
//
// @contract:mutating
func (c *SelfUpdateContract) UpdateContract(wasmBytes []byte) {
	caller, err := env.GetPredecessorAccountID()
	if err != nil {
		env.PanicStr("Failed to get caller")
	}

	if caller != c.OwnerId {
		env.PanicStr("Only the owner can update the contract")
	}

	currentAccountId, _ := env.GetCurrentAccountId()

	// Deploy the new wasm to this contract's own account
	promise.CreateBatch(currentAccountId).
		DeployContract(wasmBytes)

	env.LogString("Contract update initiated by " + caller)
}

// @contract:view
func (c *SelfUpdateContract) GetOwner() string {
	return c.OwnerId
}

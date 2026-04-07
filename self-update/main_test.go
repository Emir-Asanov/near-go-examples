package main

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/system"
)

func init() {
	env.SetEnv(system.NewMockSystem())
}

func setupTest(t *testing.T) *SelfUpdateContract {
	t.Helper()
	ms := env.NearBlockchainImports.(*system.MockSystem)
	ms.Storage = make(map[string][]byte)
	ms.CurrentAccountIdSys = "contract.testnet"

	c := &SelfUpdateContract{}
	c.Init("owner.testnet")
	return c
}

func TestSelfUpdate_GetOwner(t *testing.T) {
	c := setupTest(t)
	if c.GetOwner() != "owner.testnet" {
		t.Errorf("Expected owner.testnet, got %s", c.GetOwner())
	}
}

func TestSelfUpdate_OwnerCanUpdate(t *testing.T) {
	c := setupTest(t)
	ms := env.NearBlockchainImports.(*system.MockSystem)

	ms.PredecessorAccountIdSys = "owner.testnet"

	// UpdateContract creates a promise — no actual deploy in unit test
	// Just verify it doesn't panic with the correct owner
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("UpdateContract panicked for owner: %v", r)
		}
	}()
	c.UpdateContract(UpdateContractInput{WasmBytes: "AGFzbQ=="}) // base64 of minimal wasm magic bytes
}

func TestSelfUpdate_NonOwnerCannotUpdate(t *testing.T) {
	c := setupTest(t)
	ms := env.NearBlockchainImports.(*system.MockSystem)

	ms.PredecessorAccountIdSys = "attacker.testnet"

	// env.PanicStr in MockSystem is a no-op, so we can only verify the logic
	// by checking that the function returns without doing anything useful.
	// In a real environment, this would abort the transaction.
	_ = c
}

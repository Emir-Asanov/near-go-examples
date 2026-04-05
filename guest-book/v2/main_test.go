package main

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/system"
	"github.com/vlmoon99/near-sdk-go/types"
)

func init() {
	env.SetEnv(system.NewMockSystem())
}

func setupTest(t *testing.T) *GuestBookV2 {
	t.Helper()
	ms := env.NearBlockchainImports.(*system.MockSystem)
	ms.Storage = make(map[string][]byte)
	c := &GuestBookV2{}
	c.Init()
	return c
}

func ms() *system.MockSystem {
	return env.NearBlockchainImports.(*system.MockSystem)
}

func TestV2_AddMessage(t *testing.T) {
	c := setupTest(t)
	ms().PredecessorAccountIdSys = "alice.testnet"
	small, _ := types.U128FromString("1000000000000000000000") // 0.001 NEAR
	ms().AttachedDepositSys = small

	c.AddMessage("Hello v2!")

	msgs := c.GetMessages()
	if len(msgs) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(msgs))
	}
	if msgs[0].Text != "Hello v2!" {
		t.Errorf("Expected 'Hello v2!', got '%s'", msgs[0].Text)
	}
	if msgs[0].Sender != "alice.testnet" {
		t.Errorf("Expected alice.testnet, got %s", msgs[0].Sender)
	}
	if msgs[0].Payment != small.String() {
		t.Errorf("Expected payment %s, got %s", small.String(), msgs[0].Payment)
	}
}

func TestV2_PremiumMessage(t *testing.T) {
	c := setupTest(t)
	ms().PredecessorAccountIdSys = "bob.testnet"
	premium, _ := types.U128FromString("10000000000000000000000") // 0.01 NEAR (threshold)
	ms().AttachedDepositSys = premium

	c.AddMessage("Premium!")

	msgs := c.GetMessages()
	if !msgs[0].Premium {
		t.Error("Expected message to be premium at 0.01 NEAR")
	}
}

func TestV2_NonPremiumMessage(t *testing.T) {
	c := setupTest(t)
	ms().PredecessorAccountIdSys = "carol.testnet"
	below, _ := types.U128FromString("9999999999999999999999") // just below threshold
	ms().AttachedDepositSys = below

	c.AddMessage("Not premium")

	msgs := c.GetMessages()
	if msgs[0].Premium {
		t.Error("Expected message to be non-premium below threshold")
	}
}

// Note: TestV2_Migrate is intentionally omitted from unit tests.
// Vector.Len is an in-memory field; unit tests bypass the state
// serialization/deserialization layer that provides it in production.
// Migration is fully covered by the integration tests in integration-tests/.

func TestV2_MultipleMessages(t *testing.T) {
	c := setupTest(t)
	small, _ := types.U128FromString("1000000000000000000000")

	for _, sender := range []string{"a.testnet", "b.testnet", "c.testnet"} {
		ms().PredecessorAccountIdSys = sender
		ms().AttachedDepositSys = small
		c.AddMessage("msg from " + sender)
	}

	msgs := c.GetMessages()
	if len(msgs) != 3 {
		t.Errorf("Expected 3 messages, got %d", len(msgs))
	}
}

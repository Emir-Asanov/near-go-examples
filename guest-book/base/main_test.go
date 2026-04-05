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

func setupTest(t *testing.T) *GuestBook {
	t.Helper()
	ms := env.NearBlockchainImports.(*system.MockSystem)
	ms.Storage = make(map[string][]byte)

	c := &GuestBook{}
	c.Init()
	return c
}

func TestGuestBook_AddMessage(t *testing.T) {
	c := setupTest(t)
	ms := env.NearBlockchainImports.(*system.MockSystem)

	ms.PredecessorAccountIdSys = "alice.testnet"
	small, _ := types.U128FromString("1000000000000000000000") // 0.001 NEAR (not premium)
	ms.AttachedDepositSys = small

	c.AddMessage("Hello, NEAR!")

	msgs := c.GetMessages()
	if len(msgs) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(msgs))
	}
	if msgs[0].Text != "Hello, NEAR!" {
		t.Errorf("Expected 'Hello, NEAR!', got '%s'", msgs[0].Text)
	}
	if msgs[0].Sender != "alice.testnet" {
		t.Errorf("Expected sender alice.testnet, got %s", msgs[0].Sender)
	}
	if msgs[0].Premium {
		t.Error("Message should not be premium for small deposit")
	}
}

func TestGuestBook_PremiumMessage(t *testing.T) {
	c := setupTest(t)
	ms := env.NearBlockchainImports.(*system.MockSystem)

	ms.PredecessorAccountIdSys = "bob.testnet"
	premium, _ := types.U128FromString("10000000000000000000000") // 0.01 NEAR (premium threshold)
	ms.AttachedDepositSys = premium

	c.AddMessage("Premium message!")

	msgs := c.GetMessages()
	if !msgs[0].Premium {
		t.Error("Message should be premium for 0.01 NEAR deposit")
	}
}

func TestGuestBook_MultipleMessages(t *testing.T) {
	c := setupTest(t)
	ms := env.NearBlockchainImports.(*system.MockSystem)

	small, _ := types.U128FromString("1000000000000000000000")

	for i, sender := range []string{"alice.testnet", "bob.testnet", "carol.testnet"} {
		ms.PredecessorAccountIdSys = sender
		ms.AttachedDepositSys = small
		c.AddMessage("Message " + types.IntToString(i))
	}

	msgs := c.GetMessages()
	if len(msgs) != 3 {
		t.Errorf("Expected 3 messages, got %d", len(msgs))
	}
}

func TestGuestBook_PaymentsRecorded(t *testing.T) {
	c := setupTest(t)
	ms := env.NearBlockchainImports.(*system.MockSystem)

	ms.PredecessorAccountIdSys = "alice.testnet"
	oneNear, _ := types.U128FromString("1000000000000000000000000")
	ms.AttachedDepositSys = oneNear

	c.AddMessage("Paid message")

	payments := c.GetPayments()
	if len(payments) != 1 {
		t.Fatalf("Expected 1 payment, got %d", len(payments))
	}
	if payments[0] != "1000000000000000000000000" {
		t.Errorf("Expected payment 1000000000000000000000000, got %s", payments[0])
	}
}

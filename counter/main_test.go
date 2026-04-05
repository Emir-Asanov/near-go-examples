package main

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/system"
)

func init() {
	env.SetEnv(system.NewMockSystem())
}

func setupTest(t *testing.T) *CounterContract {
	mockSys, ok := env.NearBlockchainImports.(*system.MockSystem)
	if !ok {
		t.Fatal("Environment is not set to MockSystem")
	}
	mockSys.Storage = make(map[string][]byte)

	contract := &CounterContract{}
	contract.Init()
	return contract
}

// ── Snippet I from docs (unit-test.md) ───────────────────────────────────────

func TestCounter_Increment(t *testing.T) {
	c := setupTest(t)

	c.Increment()
	c.Increment()

	if c.GetCount() != 2 {
		t.Errorf("Expected count 2, got %d", c.GetCount())
	}
}

func TestCounter_Decrement(t *testing.T) {
	c := setupTest(t)

	c.Increment()
	c.Decrement()

	if c.GetCount() != 0 {
		t.Errorf("Expected count 0, got %d", c.GetCount())
	}
}

func TestCounter_Reset(t *testing.T) {
	c := setupTest(t)

	c.Increment()
	c.Increment()
	c.Reset()

	if c.GetCount() != 0 {
		t.Errorf("Expected count 0 after reset, got %d", c.GetCount())
	}
}

// ── Additional tests ──────────────────────────────────────────────────────────

func TestCounter_InitValue(t *testing.T) {
	c := setupTest(t)
	if c.GetCount() != 0 {
		t.Errorf("Expected initial count 0, got %d", c.GetCount())
	}
}

func TestCounter_DecrementBelowZero(t *testing.T) {
	c := setupTest(t)

	c.Decrement()

	if c.GetCount() != -1 {
		t.Errorf("Expected count -1, got %d", c.GetCount())
	}
}

func TestCounter_IncrementMany(t *testing.T) {
	c := setupTest(t)

	for i := 0; i < 100; i++ {
		c.Increment()
	}
	if c.GetCount() != 100 {
		t.Errorf("Expected count 100, got %d", c.GetCount())
	}
}

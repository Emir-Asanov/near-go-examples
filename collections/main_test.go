package main

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/system"
)

func init() {
	env.SetEnv(system.NewMockSystem())
}

func setupTest(t *testing.T) *CollectionsContract {
	t.Helper()
	ms := env.NearBlockchainImports.(*system.MockSystem)
	ms.Storage = make(map[string][]byte)
	c := &CollectionsContract{}
	c.Init()
	return c
}

// ── Vector ────────────────────────────────────────────────────────────────────

func TestVector_PushAndGet(t *testing.T) {
	c := setupTest(t)
	c.VectorPush("alpha")
	c.VectorPush("beta")

	if got := c.VectorGet(0); got != "alpha" {
		t.Errorf("Expected alpha, got %s", got)
	}
	if got := c.VectorGet(1); got != "beta" {
		t.Errorf("Expected beta, got %s", got)
	}
}

func TestVector_GetAll(t *testing.T) {
	c := setupTest(t)
	c.VectorPush("a")
	c.VectorPush("b")
	c.VectorPush("c")

	all := c.VectorGetAll()
	if len(all) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(all))
	}
}

func TestVector_Pop(t *testing.T) {
	c := setupTest(t)
	c.VectorPush("only")
	val := c.VectorPop()
	if val != "only" {
		t.Errorf("Expected 'only', got %s", val)
	}
	if c.MyVector.Length() != 0 {
		t.Error("Vector should be empty after pop")
	}
}

func TestVector_Pagination(t *testing.T) {
	c := setupTest(t)
	for i := 0; i < 5; i++ {
		c.VectorPush("item")
	}

	page := c.GetPage(1, 3)
	if len(page) != 3 {
		t.Errorf("Expected 3 items, got %d", len(page))
	}
	page2 := c.GetPage(4, 10) // only 1 item left
	if len(page2) != 1 {
		t.Errorf("Expected 1 item, got %d", len(page2))
	}
	empty := c.GetPage(10, 5) // out of bounds
	if len(empty) != 0 {
		t.Errorf("Expected 0 items, got %d", len(empty))
	}
}

// ── LookupMap ─────────────────────────────────────────────────────────────────

func TestLookupMap_SetAndGet(t *testing.T) {
	c := setupTest(t)
	c.LookupMapSet("alice", "100")

	if got := c.LookupMapGet("alice"); got != "100" {
		t.Errorf("Expected 100, got %s", got)
	}
	if got := c.LookupMapGet("unknown"); got != "" {
		t.Errorf("Expected empty for missing key, got %s", got)
	}
}

func TestLookupMap_Contains(t *testing.T) {
	c := setupTest(t)
	c.LookupMapSet("alice", "100")

	if !c.LookupMapContains("alice") {
		t.Error("Expected alice to be present")
	}
	if c.LookupMapContains("bob") {
		t.Error("Expected bob to be absent")
	}
}

func TestLookupMap_Remove(t *testing.T) {
	c := setupTest(t)
	c.LookupMapSet("alice", "100")
	c.LookupMapRemove("alice")

	if c.LookupMapContains("alice") {
		t.Error("Expected alice to be removed")
	}
}

// ── UnorderedMap ──────────────────────────────────────────────────────────────

func TestUnorderedMap_SetAndGet(t *testing.T) {
	c := setupTest(t)
	c.UnorderedMapSet("k1", "v1")
	c.UnorderedMapSet("k2", "v2")

	if got := c.UnorderedMapGet("k1"); got != "v1" {
		t.Errorf("Expected v1, got %s", got)
	}
}

func TestUnorderedMap_KeysAndValues(t *testing.T) {
	c := setupTest(t)
	c.UnorderedMapSet("a", "1")
	c.UnorderedMapSet("b", "2")

	keys := c.UnorderedMapKeys()
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}
	values := c.UnorderedMapValues()
	if len(values) != 2 {
		t.Errorf("Expected 2 values, got %d", len(values))
	}
}

func TestUnorderedMap_Remove(t *testing.T) {
	c := setupTest(t)
	c.UnorderedMapSet("key", "val")
	c.UnorderedMapRemove("key")

	if got := c.UnorderedMapGet("key"); got != "" {
		t.Errorf("Expected empty after remove, got %s", got)
	}
}

// ── LookupSet ─────────────────────────────────────────────────────────────────

func TestLookupSet_AddAndContains(t *testing.T) {
	c := setupTest(t)
	c.LookupSetAdd("alice.testnet")

	if !c.LookupSetContains("alice.testnet") {
		t.Error("Expected alice.testnet in set")
	}
	if c.LookupSetContains("bob.testnet") {
		t.Error("Expected bob.testnet not in set")
	}
}

func TestLookupSet_Remove(t *testing.T) {
	c := setupTest(t)
	c.LookupSetAdd("alice.testnet")
	c.LookupSetRemove("alice.testnet")

	if c.LookupSetContains("alice.testnet") {
		t.Error("Expected alice.testnet to be removed")
	}
}

// ── UnorderedSet ──────────────────────────────────────────────────────────────

func TestUnorderedSet_AddAndAll(t *testing.T) {
	c := setupTest(t)
	c.UnorderedSetAdd("owner1")
	c.UnorderedSetAdd("owner2")

	all := c.UnorderedSetAll()
	if len(all) != 2 {
		t.Errorf("Expected 2 owners, got %d", len(all))
	}
}

func TestUnorderedSet_Contains(t *testing.T) {
	c := setupTest(t)
	c.UnorderedSetAdd("owner1")

	if !c.UnorderedSetContains("owner1") {
		t.Error("Expected owner1 in set")
	}
	if c.UnorderedSetContains("owner2") {
		t.Error("Expected owner2 not in set")
	}
}

func TestUnorderedSet_Remove(t *testing.T) {
	c := setupTest(t)
	c.UnorderedSetAdd("owner1")
	c.UnorderedSetRemove("owner1")

	if c.UnorderedSetContains("owner1") {
		t.Error("Expected owner1 to be removed")
	}
}

// ── TreeMap ───────────────────────────────────────────────────────────────────

func TestTreeMap_SetAndMinMax(t *testing.T) {
	c := setupTest(t)
	c.TreeMapSet("alice", 100)
	c.TreeMapSet("bob", 200)
	c.TreeMapSet("carol", 50)

	allKeys := c.TreeMapGetAllKeys()
	if len(allKeys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(allKeys))
	}
}

func TestTreeMap_MinMaxKeys(t *testing.T) {
	c := setupTest(t)
	c.TreeMapSet("b_user", 10)
	c.TreeMapSet("a_user", 20)
	c.TreeMapSet("c_user", 30)

	min := c.TreeMapGetMin()
	max := c.TreeMapGetMax()

	if min != "a_user" {
		t.Errorf("Expected min a_user, got %s", min)
	}
	if max != "c_user" {
		t.Errorf("Expected max c_user, got %s", max)
	}
}

// ── Nested Collections ────────────────────────────────────────────────────────

func TestNested_AddAndGetAssets(t *testing.T) {
	c := setupTest(t)
	c.AddAsset("alice.testnet", "nft-001")
	c.AddAsset("alice.testnet", "nft-002")
	c.AddAsset("bob.testnet", "nft-100")

	aliceAssets := c.GetUserAssets("alice.testnet")
	if len(aliceAssets) != 2 {
		t.Errorf("Expected 2 assets for alice, got %d", len(aliceAssets))
	}

	bobAssets := c.GetUserAssets("bob.testnet")
	if len(bobAssets) != 1 {
		t.Errorf("Expected 1 asset for bob, got %d", len(bobAssets))
	}

	emptyAssets := c.GetUserAssets("carol.testnet")
	if len(emptyAssets) != 0 {
		t.Errorf("Expected 0 assets for carol, got %d", len(emptyAssets))
	}
}

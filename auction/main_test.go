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

// blockMs converts milliseconds to the nanosecond value that MockSystem expects.
// GetBlockTimeMs() = BlockTimestampSys / 1_000_000, so we store ms * 1_000_000.
func blockMs(ms uint64) uint64 { return ms * 1_000_000 }

// auctionEndMs is the end time used across tests (milliseconds).
const auctionEndMs = uint64(9_000_000)

func setupTest(t *testing.T) *AuctionContract {
	t.Helper()
	ms := env.NearBlockchainImports.(*system.MockSystem)
	ms.Storage = make(map[string][]byte)
	ms.CurrentAccountIdSys = "auction.testnet"
	ms.BlockTimestampSys = blockMs(1000) // current time = 1 000 ms (well before end)

	c := &AuctionContract{}
	c.Init(InitInput{
		EndTime:    auctionEndMs,
		Auctioneer: "auctioneer.testnet",
	})
	return c
}

func mockSys() *system.MockSystem {
	return env.NearBlockchainImports.(*system.MockSystem)
}

// ── Init ─────────────────────────────────────────────────────────────────────

func TestAuction_Init(t *testing.T) {
	c := setupTest(t)

	if c.GetAuctioneer() != "auctioneer.testnet" {
		t.Errorf("Expected auctioneer.testnet, got %s", c.GetAuctioneer())
	}
	if c.GetAuctionEndTime() != auctionEndMs {
		t.Errorf("Expected end time %d, got %d", auctionEndMs, c.GetAuctionEndTime())
	}
	if c.GetClaimed() {
		t.Error("Auction should not be claimed at start")
	}
	// Highest bid should be initialised with the contract account and bid "1"
	if c.GetHighestBid().Bidder != "auction.testnet" {
		t.Errorf("Expected initial bidder auction.testnet, got %s", c.GetHighestBid().Bidder)
	}
}

// ── Bid ──────────────────────────────────────────────────────────────────────

func TestAuction_BidAccepted(t *testing.T) {
	c := setupTest(t)
	ms := mockSys()

	ms.PredecessorAccountIdSys = "alice.testnet"
	twoNear, _ := types.U128FromString("2000000000000000000000000")
	ms.AttachedDepositSys = twoNear

	if err := c.Bid(); err != nil {
		t.Fatalf("Bid should succeed: %v", err)
	}

	bid := c.GetHighestBid()
	if bid.Bidder != "alice.testnet" {
		t.Errorf("Expected alice.testnet as highest bidder, got %s", bid.Bidder)
	}
	if bid.Bid != "2000000000000000000000000" {
		t.Errorf("Expected bid 2000000000000000000000000, got %s", bid.Bid)
	}
}

func TestAuction_BidMustBeHigher(t *testing.T) {
	c := setupTest(t)
	ms := mockSys()

	ms.PredecessorAccountIdSys = "alice.testnet"
	twoNear, _ := types.U128FromString("2000000000000000000000000")
	ms.AttachedDepositSys = twoNear
	c.Bid() //nolint

	// Bob tries the same amount — must fail
	ms.PredecessorAccountIdSys = "bob.testnet"
	ms.AttachedDepositSys = twoNear
	if err := c.Bid(); err == nil {
		t.Error("Expected error: bid must be strictly higher than current")
	}
}

func TestAuction_MultipleBids(t *testing.T) {
	c := setupTest(t)
	ms := mockSys()

	// Alice bids 2 NEAR
	ms.PredecessorAccountIdSys = "alice.testnet"
	twoNear, _ := types.U128FromString("2000000000000000000000000")
	ms.AttachedDepositSys = twoNear
	c.Bid() //nolint

	// Bob outbids with 3 NEAR
	ms.PredecessorAccountIdSys = "bob.testnet"
	threeNear, _ := types.U128FromString("3000000000000000000000000")
	ms.AttachedDepositSys = threeNear
	if err := c.Bid(); err != nil {
		t.Fatalf("Bob's bid should succeed: %v", err)
	}

	if c.GetHighestBid().Bidder != "bob.testnet" {
		t.Errorf("Expected bob.testnet as highest bidder, got %s", c.GetHighestBid().Bidder)
	}
}

func TestAuction_BidAfterEnd_Rejected(t *testing.T) {
	c := setupTest(t)
	ms := mockSys()

	// Move time past auction end
	ms.BlockTimestampSys = blockMs(auctionEndMs + 1)

	ms.PredecessorAccountIdSys = "alice.testnet"
	twoNear, _ := types.U128FromString("2000000000000000000000000")
	ms.AttachedDepositSys = twoNear

	if err := c.Bid(); err == nil {
		t.Error("Expected error: auction has ended")
	}
}

// ── Claim ────────────────────────────────────────────────────────────────────

func TestAuction_CannotClaimBeforeEnd(t *testing.T) {
	c := setupTest(t)

	if err := c.Claim(); err == nil {
		t.Error("Expected error: auction has not ended yet")
	}
}

func TestAuction_ClaimAfterEnd(t *testing.T) {
	c := setupTest(t)
	ms := mockSys()

	// Place a bid first
	ms.PredecessorAccountIdSys = "alice.testnet"
	twoNear, _ := types.U128FromString("2000000000000000000000000")
	ms.AttachedDepositSys = twoNear
	c.Bid() //nolint

	// Move time past auction end
	ms.BlockTimestampSys = blockMs(auctionEndMs + 1)

	if err := c.Claim(); err != nil {
		t.Fatalf("Claim should succeed after auction ends: %v", err)
	}
	if !c.GetClaimed() {
		t.Error("Auction should be marked as claimed")
	}
}

func TestAuction_CannotClaimTwice(t *testing.T) {
	c := setupTest(t)
	ms := mockSys()

	ms.BlockTimestampSys = blockMs(auctionEndMs + 1)
	c.Claim() //nolint

	if err := c.Claim(); err == nil {
		t.Error("Expected error: auction already claimed")
	}
}

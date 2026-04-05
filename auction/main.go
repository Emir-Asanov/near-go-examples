package main

import (
	"errors"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/promise"
	"github.com/vlmoon99/near-sdk-go/types"
)

// Bid represents a single bid placed in an auction.
type Bid struct {
	Bidder string `json:"bidder"`
	Bid    string `json:"bid"`
}

type InitInput struct {
	EndTime    uint64 `json:"end_time"`
	Auctioneer string `json:"auctioneer"`
}

// @contract:state
type AuctionContract struct {
	HighestBid     Bid    `json:"highest_bid"`
	AuctionEndTime uint64 `json:"auction_end_time"`
	Auctioneer     string `json:"auctioneer"`
	Claimed        bool   `json:"claimed"`
}

// @contract:init
func (c *AuctionContract) Init(input InitInput) {
	currentAccount, _ := env.GetCurrentAccountId()
	c.HighestBid = Bid{
		Bidder: currentAccount,
		Bid:    "1",
	}
	c.AuctionEndTime = input.EndTime
	c.Auctioneer = input.Auctioneer
	env.LogString("Auction initialized")
}

// @contract:payable min_deposit=1yoctoNEAR
// @contract:mutating
func (c *AuctionContract) Bid() error {
	blockTime := env.GetBlockTimeMs()
	if blockTime >= c.AuctionEndTime {
		return errors.New("auction has ended")
	}

	deposit, err := env.GetAttachedDeposit()
	if err != nil {
		return errors.New("failed to get attached deposit")
	}

	currentBid, err := types.U128FromString(c.HighestBid.Bid)
	if err != nil {
		return errors.New("invalid current bid amount in state")
	}

	if deposit.Cmp(currentBid) <= 0 {
		return errors.New("you must place a higher bid")
	}

	caller, err := env.GetPredecessorAccountID()
	if err != nil {
		return errors.New("failed to get caller account")
	}

	prevBidder := c.HighestBid.Bidder
	c.HighestBid = Bid{
		Bidder: caller,
		Bid:    deposit.String(),
	}

	promise.CreateBatch(prevBidder).Transfer(currentBid)

	return nil
}

// @contract:mutating
func (c *AuctionContract) Claim() error {
	blockTime := env.GetBlockTimeMs()
	if blockTime <= c.AuctionEndTime {
		return errors.New("auction has not ended yet")
	}

	if c.Claimed {
		return errors.New("auction has already been claimed")
	}

	c.Claimed = true

	winningBid, err := types.U128FromString(c.HighestBid.Bid)
	if err != nil {
		return errors.New("invalid winning bid amount in state")
	}

	promise.CreateBatch(c.Auctioneer).Transfer(winningBid)

	return nil
}

// @contract:view
func (c *AuctionContract) GetHighestBid() Bid {
	return c.HighestBid
}

// @contract:view
func (c *AuctionContract) GetAuctionEndTime() uint64 {
	return c.AuctionEndTime
}

// @contract:view
func (c *AuctionContract) GetAuctioneer() string {
	return c.Auctioneer
}

// @contract:view
func (c *AuctionContract) GetClaimed() bool {
	return c.Claimed
}

// Base version of the Guest Book contract.
// Messages store premium status, sender, and text.
// Payments are tracked in a separate Vector.
package main

import (
	"github.com/vlmoon99/near-sdk-go/collections"
	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/types"
)

type PostedMessage struct {
	Premium bool   `json:"premium"`
	Sender  string `json:"sender"`
	Text    string `json:"text"`
}

// @contract:state
type GuestBook struct {
	Messages *collections.Vector[PostedMessage] `json:"messages"`
	Payments *collections.Vector[string]        `json:"payments"`
}

// @contract:init
func (c *GuestBook) Init() {
	c.Messages = collections.NewVector[PostedMessage]("m")
	c.Payments = collections.NewVector[string]("p")
	env.LogString("GuestBook initialized")
}

// @contract:payable min_deposit=1yoctoNEAR
// @contract:mutating
func (c *GuestBook) AddMessage(text string) {
	sender, _ := env.GetPredecessorAccountID()
	deposit, _ := env.GetAttachedDeposit()

	premiumThreshold, _ := types.U128FromString("10000000000000000000000") // 0.01 NEAR
	isPremium := deposit.Cmp(premiumThreshold) >= 0

	c.Messages.Push(PostedMessage{
		Premium: isPremium,
		Sender:  sender,
		Text:    text,
	})
	c.Payments.Push(deposit.String())

	env.LogString("Message added by " + sender)
}

// @contract:view
func (c *GuestBook) GetMessages() []PostedMessage {
	count := c.Messages.Length()
	result := make([]PostedMessage, 0, count)
	for i := uint64(0); i < count; i++ {
		msg, err := c.Messages.Get(i)
		if err != nil {
			continue
		}
		result = append(result, msg)
	}
	return result
}

// @contract:view
func (c *GuestBook) GetPayments() []string {
	count := c.Payments.Length()
	result := make([]string, 0, count)
	for i := uint64(0); i < count; i++ {
		p, err := c.Payments.Get(i)
		if err != nil {
			continue
		}
		result = append(result, p)
	}
	return result
}

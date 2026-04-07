// V2 of the Guest Book contract.
// Payment is now embedded in PostedMessage (no separate Payments vector).
// The Migrate() method reads the old state and transforms it to the new structure.
package main

import (
	encodingJson "encoding/json"

	"github.com/vlmoon99/near-sdk-go/collections"
	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/types"
)

// Updated PostedMessage now includes the payment field
type PostedMessage struct {
	Premium bool   `json:"premium"`
	Sender  string `json:"sender"`
	Text    string `json:"text"`
	Payment string `json:"payment"` // new field — payment is now embedded in the message
}

// @contract:state
type GuestBookV2 struct {
	Messages *collections.Vector[PostedMessage] `json:"messages"`
	// Payments vector removed — payment is now part of PostedMessage
}

// @contract:init
func (c *GuestBookV2) Init() {
	c.Messages = collections.NewVector[PostedMessage]("m")
	env.LogString("GuestBookV2 initialized")
}

// OldMessage is the previous PostedMessage structure (without Payment).
// Used only during migration.
type OldMessage struct {
	Premium bool   `json:"premium"`
	Sender  string `json:"sender"`
	Text    string `json:"text"`
}

// OldGuestBook represents the previous contract state.
// Used only during migration to read and transform existing data.
type OldGuestBook struct {
	Messages *collections.Vector[OldMessage] `json:"messages"`
	Payments *collections.Vector[string]     `json:"payments"`
}

// Migrate reads the old state, merges payments into messages, and writes the new state.
// Called explicitly after deploying V2 over V1 — not an init, but a regular mutating call.
//
// @contract:mutating
func (c *GuestBookV2) Migrate() {
	// Read the raw state bytes so we can decode as OldGuestBook and recover vector lengths.
	// (NewVector("m") starts with Len=0; we need the persisted length from the v1 state.)
	rawState, err := env.StateRead()
	if err != nil || len(rawState) == 0 {
		env.PanicStr("No state to migrate")
	}

	var old OldGuestBook
	if jsonErr := encodingJson.Unmarshal(rawState, &old); jsonErr != nil {
		env.PanicStr("Failed to deserialize old state")
	}

	count := old.Messages.Length()

	// Reset with the same prefix — Push() will overwrite element keys in-place.
	c.Messages = collections.NewVector[PostedMessage]("m")

	for i := uint64(0); i < count; i++ {
		oldMsg, _ := old.Messages.Get(i)
		payment, _ := old.Payments.Get(i)

		c.Messages.Push(PostedMessage{
			Premium: oldMsg.Premium,
			Sender:  oldMsg.Sender,
			Text:    oldMsg.Text,
			Payment: payment,
		})
	}

	// Clear old payments storage entries.
	old.Payments.Clear()

	env.LogString("Migration completed successfully")
}

// @contract:payable min_deposit=1yoctoNEAR
// @contract:mutating
func (c *GuestBookV2) AddMessage(text string) {
	sender, _ := env.GetPredecessorAccountID()
	deposit, _ := env.GetAttachedDeposit()

	premiumThreshold, _ := types.U128FromString("10000000000000000000000") // 0.01 NEAR
	isPremium := deposit.Cmp(premiumThreshold) >= 0

	c.Messages.Push(PostedMessage{
		Premium: isPremium,
		Sender:  sender,
		Text:    text,
		Payment: deposit.String(),
	})

	env.LogString("Message added by " + sender)
}

// @contract:view
func (c *GuestBookV2) GetMessages() []PostedMessage {
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

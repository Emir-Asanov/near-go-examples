package main

import (
	"github.com/vlmoon99/near-sdk-go/collections"
	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/promise"
	"github.com/vlmoon99/near-sdk-go/types"
)

type DonationRecord struct {
	Donor  string `json:"donor"`
	Amount string `json:"amount"`
}

// @contract:state
type DonationContract struct {
	Beneficiary string                               `json:"beneficiary"`
	Donations   *collections.Vector[DonationRecord] `json:"donations"`
}

// @contract:init
func (c *DonationContract) Init(beneficiary string) {
	c.Beneficiary = beneficiary
	c.Donations = collections.NewVector[DonationRecord]("d")
	env.LogString("DonationContract initialized for " + beneficiary)
}

// @contract:payable min_deposit=1yoctoNEAR
// @contract:mutating
func (c *DonationContract) Donate() string {
	donor, err := env.GetPredecessorAccountID()
	if err != nil {
		env.PanicStr("Failed to get donor account")
	}

	deposit, err := env.GetAttachedDeposit()
	if err != nil {
		env.PanicStr("Failed to get attached deposit")
	}

	zero, _ := types.U128FromString("0")
	if deposit.Cmp(zero) == 0 {
		env.PanicStr("Attached deposit must be greater than 0")
	}

	record := DonationRecord{
		Donor:  donor,
		Amount: deposit.String(),
	}
	c.Donations.Push(record)

	env.LogString("Donation from " + donor + ": " + deposit.String() + " yoctoNEAR → " + c.Beneficiary)

	// Forward funds to beneficiary
	promise.CreateBatch(c.Beneficiary).Transfer(deposit)

	return "Thank you, " + donor + "! Your donation of " + deposit.String() + " yoctoNEAR was sent."
}

// @contract:view
func (c *DonationContract) GetBeneficiary() string {
	return c.Beneficiary
}

// @contract:view
func (c *DonationContract) GetDonations() []DonationRecord {
	count := c.Donations.Length()
	result := make([]DonationRecord, 0, count)
	for i := uint64(0); i < count; i++ {
		record, err := c.Donations.Get(i)
		if err != nil {
			continue
		}
		result = append(result, record)
	}
	return result
}

// @contract:view
func (c *DonationContract) GetTotalDonations() string {
	total, _ := types.U128FromString("0")
	count := c.Donations.Length()
	for i := uint64(0); i < count; i++ {
		record, err := c.Donations.Get(i)
		if err != nil {
			continue
		}
		amount, err := types.U128FromString(record.Amount)
		if err != nil {
			continue
		}
		total, _ = total.Add(amount)
	}
	return total.String()
}

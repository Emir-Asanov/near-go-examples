package main

import (
	"github.com/vlmoon99/near-sdk-go/env"
)

// @contract:state
type CounterContract struct {
	Count int64 `json:"count"`
}

// @contract:init
func (c *CounterContract) Init() {
	c.Count = 0
	env.LogString("CounterContract initialized")
}

// @contract:view
func (c *CounterContract) GetCount() int64 {
	return c.Count
}

// @contract:mutating
func (c *CounterContract) Increment() {
	c.Count++
	env.LogString("Counter incremented")
}

// @contract:mutating
func (c *CounterContract) Decrement() {
	c.Count--
	env.LogString("Counter decremented")
}

// @contract:mutating
func (c *CounterContract) Reset() {
	c.Count = 0
	env.LogString("Counter reset")
}

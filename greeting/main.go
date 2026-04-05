package main

import (
	"github.com/vlmoon99/near-sdk-go/env"
)

// @contract:state
type GreetingContract struct {
	Greeting string `json:"greeting"`
}

// @contract:init
func (c *GreetingContract) Init() {
	c.Greeting = "Hello"
	env.LogString("GreetingContract initialized")
}

// @contract:view
func (c *GreetingContract) GetGreeting() string {
	return c.Greeting
}

// @contract:mutating
func (c *GreetingContract) SetGreeting(greeting string) {
	c.Greeting = greeting
	env.LogString("Saving greeting: " + greeting)
}

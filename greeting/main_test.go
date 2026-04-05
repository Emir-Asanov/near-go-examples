package main

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/system"
)

func init() {
	env.SetEnv(system.NewMockSystem())
}

func setupTest(t *testing.T) *GreetingContract {
	t.Helper()
	ms := env.NearBlockchainImports.(*system.MockSystem)
	ms.Storage = make(map[string][]byte)

	c := &GreetingContract{}
	c.Init()
	return c
}

func TestGreeting_Init(t *testing.T) {
	c := setupTest(t)
	if c.GetGreeting() != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", c.GetGreeting())
	}
}

func TestGreeting_SetAndGet(t *testing.T) {
	c := setupTest(t)

	c.SetGreeting("Howdy")
	if c.GetGreeting() != "Howdy" {
		t.Errorf("Expected 'Howdy', got '%s'", c.GetGreeting())
	}
}

func TestGreeting_UpdateTwice(t *testing.T) {
	c := setupTest(t)

	c.SetGreeting("Hello World")
	c.SetGreeting("Hey")
	if c.GetGreeting() != "Hey" {
		t.Errorf("Expected 'Hey', got '%s'", c.GetGreeting())
	}
}

func TestGreeting_EmptyString(t *testing.T) {
	c := setupTest(t)

	c.SetGreeting("")
	if c.GetGreeting() != "" {
		t.Errorf("Expected empty string, got '%s'", c.GetGreeting())
	}
}

package main

import (
	"os"
	"testing"
)

func Test__config__configure__ok(t *testing.T) {
	os.Args = []string{"parser", "-f", "customers.csv", "-r", "10000"}

	cfg := newConfig()
	err := cfg.configure()
	if err != nil {
		t.Errorf("configure() -> %v", err)
	}
}

package main

import (
	"os"
	"testing"
)

func Test__config__configure__ok(t *testing.T) {
	os.Setenv("API_URL", "url_of_api")

	cfg := newConfig()
	err := cfg.configure()
	if err != nil {
		t.Errorf("configure() -> %v", err)
	}
}

func Test__config__configure__fail(t *testing.T) {
	os.Setenv("API_URL", "")

	cfg := newConfig()
	err := cfg.configure()
	if err == nil {
		t.Errorf("configure() -> an error was expected")
	}

	if err.Error() != "API_URL environment variable is required" {
		t.Errorf("configure() -> invalid API_UR was expected")
	}
}

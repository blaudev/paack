package main

import (
	"errors"
	"os"
)

var (
	errApiUrlIsRequired = errors.New("API_URL environment variable is required")
)

type config struct {
	api string
}

func newConfig() *config {
	return &config{}
}

func (c *config) configure() error {
	c.api = os.Getenv("API_URL")
	if c.api == "" {
		return errApiUrlIsRequired
	}

	return nil
}

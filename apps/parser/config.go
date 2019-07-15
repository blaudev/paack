package main

import (
	"errors"
	"flag"
)

var (
	errInvalidArguments = errors.New("invalid arguments")
)

type config struct {
	csv  string
	rows int
}

func newConfig() *config {
	return &config{}
}

func (c *config) configure() error {
	csv := flag.String("f", "", "filepath of customers csv")
	rows := flag.Int("r", 0, "number of rows per file")
	flag.Parse()

	if *csv == "" {
		return errInvalidArguments
	}

	if *rows == 0 {
		return errInvalidArguments
	}

	c.csv = *csv
	c.rows = *rows

	return nil
}

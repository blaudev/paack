package main

import (
	"testing"
)

func Test__api__sendCustomers__ok(t *testing.T) {
	records := 1000115

	cs := make([]Customer, 0, records)
	for i := 0; i < records; i++ {
		c := Customer{
			ID: i,
			Name: "customer",
			Email: "email",
		}

		cs = append(cs, c)
	}

	api := newApi("http://localhost:5010/api")
	err := api.sendCustomers(cs)
	if err != nil {
		t.Error(err)
	}
}

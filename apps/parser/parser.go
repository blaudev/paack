package main

import (
	"strconv"
)

func parse(record []string) (Customer, error) {
	c := Customer{}
	id, err := strconv.Atoi(record[0])
	if err != nil {
		return c, err
	}

	c.ID = id
	c.Name = record[1]
	c.Email = record[2]

	return c, nil
}

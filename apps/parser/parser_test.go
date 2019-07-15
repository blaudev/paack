package main

import (
	"fmt"
	"strconv"
	"testing"
)

func Test__parser__parse__ok(t *testing.T) {
	items := make([][]string, 0, 3)
	for i := 1; i <= 3; i++ {
		items = append(items, []string{strconv.Itoa(i), fmt.Sprintf("customer %d", i), fmt.Sprintf("email@customer%d.com", i)})
	}

	for _, i := range items {
		c, err := parse(i)
		if err != nil {
			t.Errorf("parse() error = %v", err)
		}

		id, err := strconv.Atoi(i[0])
		if err != nil {
			t.Errorf("parse() -> %v", err)
		}

		if c.ID != id {
			t.Errorf("parse() -> id is not equals, expected %s, returned %d", i[0], c.ID)
		}

		if c.Name != i[1] {
			t.Errorf("parse() -> name is not equals, expected %s, returned %s", i[1], c.Name)
		}

		if c.Email != i[2] {
			t.Errorf("parse() -> email is not equals, expected %s, returned %s", i[2], c.Email)
		}
	}
}

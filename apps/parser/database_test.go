package main

import (
	"testing"
)

func Test__database__open__ok(t *testing.T) {
	db := newDatabase()
	defer db.close()

	err := db.open()
	if err != nil {
		t.Errorf("open() -> %v", err)
	}

	err = db.postgres.Ping()
	if err != nil {
		t.Errorf("open() -> %v", err)
	}
}

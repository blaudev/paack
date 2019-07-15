package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	cfg := newConfig()
	if err := cfg.configure(); err != nil {
		flag.PrintDefaults()
		os.Exit(2)
	}

	db := newDatabase()
	if err := db.open(); err != nil {
		log.Fatal(err)
	}
	defer db.close()

	in := newIntegrator()
	if err := in.dial(); err != nil {
		log.Fatal(err)
	}
	defer in.close()

	sf := newSplitter()
	if err := sf.split(cfg.csv, cfg.rows); err != nil {
		log.Fatal(err)
	}
	defer sf.close()

	for {
		rs, err := sf.next()
		if err == EOC {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		cs := make([]Customer, 0, cfg.rows)
		for _, r := range rs {
			c, err := parse(r)
			if err != nil {
				log.Fatal(err)
			}

			cs = append(cs, c)
		}

		cs, err = db.saveCustomers(cs)
		if err != nil {
			log.Fatal(err)
		}

		err = in.sendCustomers(cs)
		if err != nil {
			log.Fatal(err)
		}
	}
}

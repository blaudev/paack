package main

import (
	"log"
)

func main() {
	cfg := newConfig()
	if err := cfg.configure(); err != nil {
		log.Fatalln(err)
	}

	api := newApi(cfg.api)
	s := newService(api)
	svr := newServer(s)
	if err := svr.serve(); err != nil {
		log.Fatalln(err)
	}
	defer svr.close()
}

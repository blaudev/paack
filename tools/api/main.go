package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	apiHost = ":5010"
)

func main() {
	count := 1
	rand.Seed(time.Now().UTC().UnixNano())

	mux := http.NewServeMux()
	mux.HandleFunc("/api", func(writer http.ResponseWriter, request *http.Request) {
		c := Customer{}
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Fatalln(err)
		}

		err = json.Unmarshal(body, &c)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("CUSTOMER", c.ID)

		if rand.Intn(10) == 0 {
			writer.WriteHeader(http.StatusInternalServerError)
			log.Printf("%10d - %s - FALSE ERROR\n", count, request.Method)
			count++
			return
		}

		log.Printf("%10d - %s\n", count, request.Method)
		count++

		writer.WriteHeader(http.StatusOK)
	})

	panic(http.ListenAndServe(apiHost, mux))
}

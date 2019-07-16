package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	apiConcurrency     = 50
	numOfAttempts      = 10
	maxHttpConnections = 1000
	idleConnTimeout    = time.Second * 60
)

var (
	errSendingCustomers = errors.New("error sending customers to api")
)

type api struct {
	client *http.Client
	url    string
}

type apiResp struct {
	err error
}

func newApi(url string) *api {
	dtp, _ := http.DefaultTransport.(*http.Transport)
	dt := *dtp
	dt.MaxIdleConns = maxHttpConnections
	dt.MaxIdleConnsPerHost = maxHttpConnections
	dt.IdleConnTimeout = idleConnTimeout
	cli := &http.Client{Transport: &dt}

	return &api{
		client: cli,
		url:    url,
	}
}

func (a *api) sendCustomers(cs []Customer) error {
	jobs := make(chan Customer, len(cs))

	resp := make(chan apiResp, len(cs))
	defer close(resp)

	for i := 0; i < apiConcurrency; i++ {
		go a.worker(jobs, resp)
	}

	for _, c := range cs {
		jobs <- c
	}
	defer close(jobs)

	for i := 0; i < len(cs); i++ {
		if r := <-resp; r.err != nil {
			return r.err
		}
	}

	return nil
}

func (a *api) worker(inCh <-chan Customer, outCh chan<- apiResp) {
	for cs := range inCh {
		outCh <- a.job(cs)
	}
}

func (a *api) job(c Customer) apiResp {
	log.Printf("Sending customer %d to api\n", c.ID)

	method := http.MethodPost
	if c.Status == "to_update" {
		method = http.MethodPut
	}

	data, err := json.Marshal(c)
	if err != nil {
		return apiResp{err: err}
	}

	for i := 0; i < numOfAttempts; i++ {
		req, err := http.NewRequest(method, a.url, strings.NewReader(string(data)))
		if err != nil {
			continue
		}

		resp, err := a.client.Do(req)
		if err != nil {
			continue
		}

		if resp.StatusCode != http.StatusOK {
			continue
		}

		return apiResp{}
	}

	return apiResp{err: errSendingCustomers}
}

package main

import (
	"errors"
	"log"
	"net/rpc"
)

const (
	integratorConcurrency     = 1
	integratorRecordsPerCycle = 10000
	integratorHost            = "localhost:5002"
	integratorProcessFuncName = "Service.Process"
)

var (
	errSendingCustomers = errors.New("error sending customers to integrator")
)

// Request represents the response from the server
type Response struct {
	Ok bool
}

type integrator struct {
	client *rpc.Client
}

type integratorResp struct {
	err error
}

func newIntegrator() *integrator {
	return &integrator{}
}

func (in *integrator) dial() error {
	cl, err := rpc.DialHTTP("tcp", integratorHost)
	if err != nil {
		return err
	}

	in.client = cl

	return nil
}

func (in *integrator) sendCustomers(cs []Customer) error {
	recs := integratorRecordsPerCycle
	if len(cs) < recs {
		recs = len(cs)
	}

	size := len(cs) / recs
	if len(cs)%size != 0 {
		size++
	}

	jobs := make(chan []Customer, size)
	defer close(jobs)

	resp := make(chan integratorResp, size)
	defer close(resp)

	for i := 0; i < integratorConcurrency; i++ {
		go in.worker(jobs, resp)
	}

	for i := 0; i < size; i++ {
		cst := cs[i*recs : (i+1)*recs]
		if i+1 == size {
			cst = cs[i*recs:]
		}

		jobs <- cst
	}

	for i := 0; i < size; i++ {
		if r := <-resp; r.err != nil {
			return r.err
		}
	}

	return nil
}

func (in *integrator) worker(inCh <-chan []Customer, outCh chan<- integratorResp) {
	for cs := range inCh {
		outCh <- in.job(cs)
	}
}

func (in *integrator) job(cs []Customer) integratorResp {
	log.Printf("INTEGRATOR -> %d - %d\n", cs[0].ID, cs[len(cs)-1].ID)

	resp := &Response{}
	if err := in.client.Call(integratorProcessFuncName, cs, resp); err != nil {
		return integratorResp{err: err}
	}

	if !resp.Ok {
		return integratorResp{err: errSendingCustomers}
	}

	return integratorResp{}
}

func (in *integrator) close() error {
	return in.client.Close()
}

package main

import (
	"log"
)

// Request represents the response to the client
type Response struct {
	Ok bool
}

// Service is the integrator service
type Service struct {
	api *api
}

func newService(api *api) *Service {
	return &Service{
		api: api,
	}
}

func (sv *Service) Process(cs []Customer, resp *Response) error {
	log.Printf("Sending %d customers to api", len(cs))

	err := sv.api.sendCustomers(cs)
	if err != nil {
		return err
	}

	resp.Ok = true
	return nil
}

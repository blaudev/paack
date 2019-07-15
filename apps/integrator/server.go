package main

import (
	"net"
	"net/http"
	"net/rpc"
)

const (
	integratorHost = ":5002"
)

type server struct {
	listener net.Listener
	service  *Service
}

func newServer(sv *Service) *server {
	return &server{
		service: sv,
	}
}

func (sr *server) serve() error {
	if err := rpc.Register(sr.service); err != nil {
		return err
	}

	rpc.HandleHTTP()

	ls, err := net.Listen("tcp", integratorHost)
	if err != nil {
		return err
	}

	sr.listener = ls

	return http.Serve(ls, nil)
}

func (sr *server) close() error {
	return sr.listener.Close()
}

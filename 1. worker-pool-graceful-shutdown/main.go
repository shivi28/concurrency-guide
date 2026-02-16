package main

import (
	"fmt"
	"sync"
)

type Request struct {
	ID      int
	Payload any
	Resp    chan Response
}

type Response struct {
	ID     int
	Result any
	Err    error
}

type Processor func(*Request)

type Server struct {
	workers int
	process Processor

	wg   sync.WaitGroup
	done chan struct{}
}

func NewServer(workerCount int, processMethod Processor) *Server {
	return &Server{
		workers: workerCount,
		process: processMethod,
		done:    make(chan struct{}),
	}
}

func (s *Server) Start(reqQueue chan *Request) {
	s.wg.Add(s.workers)
	for i := 0; i < s.workers; i++ {
		go s.Worker(reqQueue)
	}
}

func (s *Server) Worker(reqQueue chan *Request) {
	defer s.wg.Done()
	for {
		select {
		case <-s.done:
			return
		case r, ok := <-reqQueue:
			if !ok {
				return
			}
			s.process(r)
			fmt.Printf("Processing done for Request ID %+v \n", r.ID)
		}
	}
}

func (s *Server) Stop() {
	close(s.done)
	s.wg.Wait()
}

func main() {
	reqQueue := make(chan *Request)

	processMethod := func(req *Request) {
		fmt.Printf("I am processing request with ID:%+v \n", req.ID)
		response := Response{
			ID:     req.ID,
			Result: "SUCCESS",
		}
		req.Resp <- response
	}

	server := NewServer(5, processMethod)

	server.Start(reqQueue)

	for i := 0; i < 10; i++ {
		reqQueue <- &Request{ID: i, Resp: make(chan Response, 1)}
	}

	server.Stop()
}

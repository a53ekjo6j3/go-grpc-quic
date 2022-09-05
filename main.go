package main

import (
	"grpc-quic/client"
	"grpc-quic/server"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.NewGrpc()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		client.NewClient()
	}()

	wg.Wait()
}

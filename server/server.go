package server

import (
	"context"
	"fmt"
	"log"

	grpcquic "grpc-quic/grpc-quic"
	proto "grpc-quic/helloword"

	"github.com/lucas-clemente/quic-go"
	"google.golang.org/grpc"
)

// Implements a server for Greeter service.
type server struct {
	proto.GreeterServer
}

func NewGrpc() {
	ql, err := quic.ListenAddr(":8081", grpcquic.GenerateTLSConfig(), nil)
	if err != nil {
		log.Fatalf("Failed to ListenAddr. %s", err.Error())
	}

	listener := grpcquic.Listen(ql)

	serv := grpc.NewServer()

	proto.RegisterGreeterServer(serv, &server{})

	log.Printf("Grpc server listening at %v", listener.Addr())
	if err := serv.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %s", err.Error())
	}
}

func (s *server) SayHello(ctx context.Context, req *proto.HelloRequest) (resp *proto.HelloReply, err error) {
	resp = &proto.HelloReply{
		Message: fmt.Sprintf("Hello! %s!", req.Name),
	}
	return
}

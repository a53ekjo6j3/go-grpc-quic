package client

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"time"

	grpcquic "grpc-quic/grpc-quic"
	proto "grpc-quic/helloword"

	"google.golang.org/grpc"
)

// Implements a client for Greeter service.
func NewClient() {
	http.HandleFunc("/", sayHello)
	log.Printf("Client listening at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func sayHello(w http.ResponseWriter, req *http.Request) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{grpcquic.NextProtos},
	}

	creds := grpcquic.NewCredentials(tlsConf)

	dialer := grpcquic.NewQuicDialer(tlsConf)
	grpcOpts := []grpc.DialOption{
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(creds),
	}

	conn, err := grpc.Dial("localhost:8081", grpcOpts...)
	if err != nil {
		log.Printf("Client error: %s", err.Error())
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Client error: %s", err.Error())
		}
	}(conn)

	cx := proto.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	result, err := cx.SayHello(ctx, &proto.HelloRequest{Name: req.URL.Query().Get("name")})
	if err == nil {
		w.Write([]byte(result.Message))
	}
}

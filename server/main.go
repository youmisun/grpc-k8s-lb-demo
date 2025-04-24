package main

import (
	"context"
	pb "demo/gen/go/proto/helloworld"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	hostname, _ := os.Hostname()
	log.Printf("Received: %v, hostname: %v", in.GetName(), hostname)
	return &pb.HelloReply{Message: fmt.Sprintf("hostname:%s, in:%s", hostname, in.Name)}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	gwmux := runtime.NewServeMux()
	err = pb.RegisterGreeterHandlerFromEndpoint(context.Background(), gwmux, lis.Addr().String(), []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}
	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}

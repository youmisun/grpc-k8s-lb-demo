package main

import (
	"context"
	pb "demo/gen/go/proto/helloworld"
	"flag"
	"fmt"
	"github.com/sercand/kuberesolver/v6"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"net/http"
	"time"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "k8s:///grpc-lb-server:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	resolver.Register(kuberesolver.NewBuilder(nil /*custom kubernetes client*/, "k8s"))
	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr,
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	http.HandleFunc("/grpc_client", func(w http.ResponseWriter, r *http.Request) {
		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		g, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		fmt.Fprint(w, fmt.Sprintf("Greeting: %s", g.GetMessage()))
	})
	log.Fatalln(http.ListenAndServe(":8081", nil))
}

package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "gprc-example/proto"
	"log"
)

const PORT = "9001"

func main() {
	conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	client := pb.NewSearchServiceClient(conn)
	resp, err := client.Search(context.Background(), &pb.SearchRequest{
		Request: "gRPC",
	})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}

	log.Printf("resp: %s", resp.GetResponse())
}

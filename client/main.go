package main

import (
	"context"
	"log"

	pb "github.com/sgrgug/kvstore/proto"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect %v", err)
	}

	defer conn.Close()

	c := pb.NewKVStoreClient(conn)

	// set a key-value pair
	_, err = c.Set(context.Background(), &pb.SetRequest{Key: "example-1", Value: "value-1"})

	if err != nil {
		log.Fatalf("couldnot set %v", err)
	}
	log.Printf("Set key 'example-1' to 'value-1'")

	// Get the value for key
	r, err := c.Get(context.Background(), &pb.GetRequest{Key: "example-1"})
	if err != nil {
		log.Fatalf("failed to get the value %v", err)
	}

	if r.GetFound() {
		log.Printf("Got the value for key 'example-1': %s", r.GetValue())
	} else {
		log.Printf("Key 'example-1' not found")
	}

	// delete the key
	_, err = c.Delete(context.Background(), &pb.DeleteRequest{Key: "example"})

	if err != nil {
		log.Fatalf("failed to delete: %v", err)
	}
	log.Println("Deleted key 'example'")

}

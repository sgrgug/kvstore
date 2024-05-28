package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/sgrgug/kvstore/proto"
	"google.golang.org/grpc"
)

type kvStoreServer struct {
	pb.UnimplementedKVStoreServer
	store map[string]string
	mu    sync.Mutex
}

func (s *kvStoreServer) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[req.GetKey()] = req.GetValue()
	return &pb.SetResponse{Success: true}, nil
}

func (s *kvStoreServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, found := s.store[req.GetKey()]
	return &pb.GetResponse{Value: value, Found: found}, nil
}

func (s *kvStoreServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, found := s.store[req.GetKey()]
	if found {
		delete(s.store, req.GetKey())
	}
	return &pb.DeleteResponse{Success: found}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	store := &kvStoreServer{store: make(map[string]string)}
	pb.RegisterKVStoreServer(s, store)
	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

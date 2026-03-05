package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "apa_aja/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {

	user := pb.UserResponse{
		Id:   req.Id,
		Name: "Budi",
		Age:  21,
	}

	fmt.Println("Request user ID:", req.Id)

	return &user, nil
}

func main() {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, &server{})

	fmt.Println("Server running on port 50051")

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
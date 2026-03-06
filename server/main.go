package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "apa_aja/proto"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
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

	wrappedGrpc := grpcweb.WrapServer(
		grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool { return true }),
	)

	go func() {
		httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if wrappedGrpc.IsGrpcWebRequest(r) || wrappedGrpc.IsGrpcWebSocketRequest(r) || wrappedGrpc.IsAcceptableGrpcCorsRequest(r) {
				wrappedGrpc.ServeHTTP(w, r)
				return
			}

			http.NotFound(w, r)
		})

		fmt.Println("gRPC-Web server running on port 8080")
		if err := http.ListenAndServe(":8080", httpHandler); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("gRPC server running on port 50051")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}

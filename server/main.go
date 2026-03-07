package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"apa_aja/database"
	pb "apa_aja/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	db := database.DB

	var user pb.UserResponse
	err := db.QueryRow("select id, name, age from user where id = ?", req.Id).Scan(&user.Id, &user.Name, &user.Age)
	if err != nil {
		fmt.Println("Gagal Mengambil Data User:", err)
		return nil, err
	}

	fmt.Println("Data User Berhasil Diambil")
	return &user, nil
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	db := database.DB

	res, err := db.Exec("insert into user (name, age) values (?, ?)", req.Name, req.Age)
	if err != nil {
		fmt.Println("Gagal Membuat User:", err)
		return nil, err
	}

	id, _ := res.LastInsertId()
	fmt.Println("User Berhasil Dibuat dengan ID:", id)

	return &pb.UserResponse{
		Id:   int32(id),
		Name: req.Name,
		Age:  req.Age,
	}, nil
}

func main() {
	database.InitDB()

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

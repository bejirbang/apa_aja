package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "apa_aja/proto"

	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.GetUser(ctx, &pb.UserRequest{Id: 1})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("User Name:", resp.Name)
	fmt.Println("Age:", resp.Age)
}

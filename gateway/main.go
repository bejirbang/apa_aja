package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	pb "apa_aja/proto"

	"google.golang.org/grpc"
)

func main() {

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {

		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}

		client := pb.NewUserServiceClient(conn)

		res, err := client.GetUser(context.Background(), &pb.UserRequest{Id: 1})
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(res)
	})

	log.Println("Gateway running on :8080")

	http.ListenAndServe(":8080", nil)
}
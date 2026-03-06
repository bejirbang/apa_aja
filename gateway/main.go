package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	pb "apa_aja/proto"

	"google.golang.org/grpc"
)

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {

		enableCORS(&w)

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
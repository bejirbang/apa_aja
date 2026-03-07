package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"strconv"

	pb "apa_aja/proto"

	"google.golang.org/grpc"
)

func enableCORS(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}

func main() {

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		if enableCORS(w, r) {
			return
		}

		// Ambil ID dari query parameter: /user?id=1
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		client := pb.NewUserServiceClient(conn)

		res, err := client.GetUser(context.Background(), &pb.UserRequest{Id: int32(id)})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(res)
	})

	http.HandleFunc("/user/create", func(w http.ResponseWriter, r *http.Request) {
		if enableCORS(w, r) {
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var reqData struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		client := pb.NewUserServiceClient(conn)

		res, err := client.CreateUser(context.Background(), &pb.CreateUserRequest{
			Name: reqData.Name,
			Age:  int32(reqData.Age),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(res)
	})

	log.Println("Gateway running on :8080")
	http.ListenAndServe(":8080", nil)
}

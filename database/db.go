package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("mysql", "root:@tcp(localhost:3306)/grpc_go")
	// DB, err = sql.Open("mysql", "USER:PASSWORD@tcp(localhost:3306)/grpc_go") (ini buat yang databasenya pakai password)
	if err != nil {
		log.Fatal(err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}
}

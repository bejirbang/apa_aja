package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:@tcp(localhost:3306)/grpc_go?parseTime=true"
	}

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("gagal koneksi ke database:", err)
	}

	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS user (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			age INT NOT NULL
		)
	`); err != nil {
		log.Fatal("gagal memastikan tabel user:", err)
	}
}

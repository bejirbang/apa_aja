package database

import (
	"database/sql"
	"log"
	"os"
	"strings"

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
			age INT NOT NULL,
			email VARCHAR(255) NULL,
			google_sub VARCHAR(255) NULL,
			avatar_url VARCHAR(512) NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		log.Fatal("gagal memastikan tabel user:", err)
	}

	ensureColumn("user", "email", "VARCHAR(255) NULL")
	ensureColumn("user", "google_sub", "VARCHAR(255) NULL")
	ensureColumn("user", "avatar_url", "VARCHAR(512) NULL")
	ensureColumn("user", "created_at", "TIMESTAMP DEFAULT CURRENT_TIMESTAMP")
	ensureIndex("user", "idx_user_google_sub", "(google_sub)", true)

	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS session (
			token VARCHAR(128) PRIMARY KEY,
			user_id INT NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_session_user_id (user_id)
		)
	`); err != nil {
		log.Fatal("gagal memastikan tabel session:", err)
	}
}

func ensureColumn(table, column, definition string) {
	query := "ALTER TABLE " + table + " ADD COLUMN " + column + " " + definition
	if _, err := DB.Exec(query); err != nil {
		if !strings.Contains(err.Error(), "Duplicate column name") {
			log.Println("gagal menambah kolom", column, "di", table+":", err)
		}
	}
}

func ensureIndex(table, indexName, definition string, unique bool) {
	createType := "INDEX"
	if unique {
		createType = "UNIQUE INDEX"
	}
	query := "CREATE " + createType + " " + indexName + " ON " + table + " " + definition
	if _, err := DB.Exec(query); err != nil {
		if !strings.Contains(err.Error(), "Duplicate key name") && !strings.Contains(err.Error(), "Duplicate index name") {
			log.Println("gagal menambah index", indexName, "di", table+":", err)
		}
	}
}

package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ryvasa/go-super-farmer/pkg/config"
)

type PostgreSQL struct {
	db *sql.DB
}

func ProvideDSN(cfg *config.Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=True&charset=utf8mb4&loc=Local", // Tambahkan charset dan loc
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Name,
	)
}

func NewPostgreSQL(dsn string) (*PostgreSQL, error) {
	var db *sql.DB
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Pastikan koneksi berhasil sebelum mengembalikan instance
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return &PostgreSQL{db: db}, nil
}

func (p *PostgreSQL) Close() {
	if p.db != nil {
		p.db.Close()
	}
}

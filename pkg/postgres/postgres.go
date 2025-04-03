package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

// psql -U postgres -d analytics
func OpenDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=analytics sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	return db, nil
}

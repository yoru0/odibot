package postgres

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type SupabaseStore struct {
	DB *pgx.Conn
}

func NewSupabaseStore() (*SupabaseStore, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	// Create table if not exist
	_, err = conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS results (
			id SERIAL PRIMARY KEY,
			num1 INTEGER,
			num2 INTEGER, 
			result INTEGER
		)
	`)
	if err != nil {
		return nil, err
	}
	log.Println("Table created")

	return &SupabaseStore{DB: conn}, nil
}

func (s *SupabaseStore) SaveResult(a, b, result int) error {
	_, err := s.DB.Exec(context.Background(),
		"INSERT INTO results (num1, num2, result) VALUES ($1, $2, $3)", a, b, result)
	return err
}

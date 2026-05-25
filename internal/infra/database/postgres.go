package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgresDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir o banco: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("banco de dados inacessível: %w", err)
	}

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("erro ao aplicar o schema e seed: %w", err)
	}

	log.Println("Postgres conectado com sucesso e tabelas prontas para o teste")
	return db, nil
}

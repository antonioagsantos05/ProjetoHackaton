package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // Driver do PostgreSQL
)

const (
	maxRetries = 5
	retryDelay = 5 * time.Second
)

// NewDBConnection cria e retorna uma nova conexão com o banco de dados PostgreSQL, com lógica de retry.
func NewDBConnection() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		host, user, password, dbname)

	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Printf("Tentativa %d/%d: erro ao abrir a conexão com o banco de dados: %v", i+1, maxRetries, err)
			time.Sleep(retryDelay)
			continue
		}

		err = db.Ping()
		if err != nil {
			log.Printf("Tentativa %d/%d: não foi possível conectar ao banco de dados (ping falhou): %v", i+1, maxRetries, err)
			db.Close() // Fecha a conexão falha
			time.Sleep(retryDelay)
			continue
		}

		log.Println("Conexão com o PostgreSQL estabelecida com sucesso!")
		return db, nil
	}

	return nil, fmt.Errorf("não foi possível conectar ao banco de dados após %d tentativas: %w", maxRetries, err)
}

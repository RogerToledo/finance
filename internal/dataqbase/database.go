package database

import (
	"database/sql"
	"fmt"

	"github.com/me/application/go/finance/internal/config"
	_ "github.com/lib/pq"
)

func NewDB() (*sql.DB, error) {
	connString := config.DB().StringConn

	conn, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to database: %v", err)
	}

	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("Error pinging database: %v", err)
	}

	return conn, nil

}
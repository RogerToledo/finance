package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/me/finance/internal/config"
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

package main

import (
	"context"
	"fmt"
	"net"

	"github.com/jackc/pgx/v5"
)

func getPostgresWaiter() waiter {
	host := envVar("DB_HOST", "localhost")
	port := envVar("DB_PORT", "5432")
	user := envVar("DB_USER", "postgres")
	password := envVar("DB_PASSWORD", "postgres")
	dbName := envVar("DB_NAME", "db")

	addr := fmt.Sprintf("%s:%s", host, port)
	dbConn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)

	return postgresWaiter{
		address:          addr,
		connectionString: dbConn,
	}
}

type postgresWaiter struct {
	address          string
	connectionString string
}

func (w postgresWaiter) waitFor() (bool, error) {
	conn, err := net.Dial("tcp", w.address)
	defer func(conn net.Conn) {
		if conn != nil {
			_ = conn.Close()
		}
	}(conn)
	if err != nil {
		return false, err
	}

	postgresConn, err := pgx.Connect(context.Background(), w.connectionString)
	if err != nil {
		return false, err
	}
	defer postgresConn.Close(context.Background())

	return true, nil
}

func (w postgresWaiter) name() string {
	return "Postgres"
}

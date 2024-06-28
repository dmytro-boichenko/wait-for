package main

import (
	"database/sql"
	"fmt"
	"net"

	_ "github.com/go-sql-driver/mysql"
)

func getMySQLWaiter() waiter {
	host := envVar("DB_HOST", "localhost")
	port := envVar("DB_PORT", "3306")
	user := envVar("DB_USER", "user")
	password := envVar("DB_PASSWORD", "password")
	dbName := envVar("DB_NAME", "db")

	addr := fmt.Sprintf("%s:%s", host, port)
	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbName)

	return mySQLWaiter{
		address:          addr,
		connectionString: dbConn,
	}
}

type mySQLWaiter struct {
	address          string
	connectionString string
}

func (w mySQLWaiter) waitFor() (bool, error) {
	conn, err := net.Dial("tcp", w.address)
	defer func(conn net.Conn) {
		if conn != nil {
			_ = conn.Close()
		}
	}(conn)

	if err != nil {
		return false, err
	}

	db, err := sql.Open("mysql", w.connectionString)
	defer func(db *sql.DB) {
		if db != nil {
			_ = db.Close()
		}
	}(db)

	if err != nil {
		return false, err
	}

	_, err = db.Exec("SHOW TABLES")
	if err != nil {
		return false, err
	}

	return true, nil
}

func (w mySQLWaiter) name() string {
	return "MySQL"
}

package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func Connect(host string, user string) (*sql.DB, error) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=postgres sslmode=disable", host, 5432, user)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(10)
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (d Database) Close() error {
	return d.DB.Close()
}

func (d Database) InitTables() error {
	err := d.InitTableTasks()
	return err
}

func (d Database) EmptyTables() error {
	err := d.EmptyTableTasks()
	return err
}

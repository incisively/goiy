package db

import (
	"database/sql"
	"database/sql/driver"
)

// Open mimics sql.Open function. However, it returns
// an sql.DB embedded with the db.database type, which
// implements db.DB interface.
func Open(driverName, dataSourceName string) (DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &database{db}, nil
}

type database struct {
	*sql.DB
}

func (db *database) Begin() (Tx, error) {
	return db.DB.Begin()
}

// DB matches the interface of sql.DB.
// This is useful for mocking within tests.
type DB interface {
	Begin() (Tx, error)
	Close() error
	Driver() driver.Driver
	Exec(query string, args ...interface{}) (sql.Result, error)
	Ping() error
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
}

// Tx matches the interface of sql.Tx.
// This is useful for mocking within tests.
type Tx interface {
	Commit() error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	Rollback() error
	Stmt(stmt *sql.Stmt) *sql.Stmt
}

package db

import (
	"database/sql"
	"database/sql/driver"
)

// MockDB implements DB.
//
// It is composed of fields to be used directly as return
// values of the functions which match the DB interface.
type MockDB struct {
	Beginf           func() (*Tx, error)
	Closef           func() error
	Driverf          func() driver.Driver
	Execf            func(query string, args ...interface{}) (sql.Result, error)
	Pingf            func() error
	Preparef         func(query string) (*sql.Stmt, error)
	Queryf           func(query string, args ...interface{}) (*sql.Rows, error)
	QueryRowf        func(query string, args ...interface{}) *sql.Row
	SetMaxIdleConnsf func(n int)
	SetMaxOpenConnsf func(n int)
}

func (m *MockDB) Begin() (*Tx, error)                                 { return m.Beginf() }
func (m *MockDB) Close() error                                        { return m.Closef() }
func (m *MockDB) Driver() driver.Driver                               { return m.Driverf() }
func (m *MockDB) Exec(q string, a ...interface{}) (sql.Result, error) { return m.Execf(q, a) }
func (m *MockDB) Ping() error                                         { return m.Pingf() }
func (m *MockDB) Prepare(q string) (*sql.Stmt, error)                 { return m.Preparef(q) }
func (m *MockDB) Query(q string, a ...interface{}) (*sql.Rows, error) { return m.Queryf(q, a) }
func (m *MockDB) QueryRow(q string, a ...interface{}) *sql.Row        { return m.QueryRowf(q, a) }
func (m *MockDB) SetMaxIdleConns(n int)                               { m.SetMaxIdleConnsf(n) }
func (m *MockDB) SetMaxOpenConns(n int)                               { m.SetMaxOpenConnsf(n) }

// MockTx implements Tx.
//
// It is composed of fields to be used directly as return
// values of the functions which match the Tx interface.
type MockTx struct {
	Commitf   func() error
	Execf     func(query string, args ...interface{}) (sql.Result, error)
	Preparef  func(query string) (*sql.Stmt, error)
	Queryf    func(query string, args ...interface{}) sql.Result
	QueryRowf func(query string, args ...interface{}) *sql.Row
	Rollbackf func() error
	Stmtf     func(stmt *sql.Stmt) *sql.Stmt
}

func (m *MockTx) Commit() error                                       { return m.Commitf() }
func (m *MockTx) Exec(q string, a ...interface{}) (sql.Result, error) { return m.Execf(q, a) }
func (m *MockTx) Prepare(q string) (*sql.Stmt, error)                 { return m.Preparef(q) }
func (m *MockTx) Query(q string, a ...interface{}) sql.Result         { return m.Queryf(q, a) }
func (m *MockTx) QueryRow(q string, a ...interface{}) *sql.Row        { return m.QueryRowf(q, a) }
func (m *MockTx) Rollback() error                                     { return m.Rollbackf() }
func (m *MockTx) Stmt(s *sql.Stmt) *sql.Stmt                          { return m.Stmt(s) }

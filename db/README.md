Database Package `db`
====================

The intention of this package is to generalise database interactions via interfaces.
The main intention is to allow for simple stubbing/mocking of database interactions for tests.

```go
// DB matches the interface of sql.DB.
// This is useful for mocking within tests.
type DB interface {
	Begin() (*sql.Tx, error)
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
	Query(string, ...interface{}) sql.Result
	QueryRow(string, ...interface{}) *sql.Row
	Rollback() error
	Stmt(stmt *sql.Stmt) *sql.Stmt
}
```

package dbtest

import (
	"bytes"
	"fmt"
	"os"

	"github.com/ConvertHQ/goincisively/db"
)

// TestDatabase is a struct, which embeds a db.DB interface
// It extends db.DB with helper methods for running
// functions wich require database transactions
type TestDatabase struct {
	db.DB
	QuoteIdentifier func(string) string
}

// NewTestDatabase returns a new TestDatabase value with the desired
// db.DB implementation and quote identifier function set.
func NewTestDatabase(db db.DB, quote func(string) string) TestDatabase {
	return TestDatabase{
		DB:              db,
		QuoteIdentifier: quote,
	}
}

// RollbackTx creates a transaction and passes it to the
// function `f`. It always results in performing a rollback
// on the transaction.
func (t TestDatabase) RollbackTx(f func(tx db.Tx) error) error {
	tx, err := t.Begin()
	if err != nil {
		return err
	}

	if err = f(tx); err != nil {
		tx.Rollback()
		return err
	}

	// always rollback the transaction
	return tx.Rollback()
}

// CommitTx creates a transaction and passes it to the function
// `f`. It will rollback the transaction if `f` produces an error.
// However, if `f` returns nil, it commits the transaction.
func (t TestDatabase) CommitTx(f func(tx db.Tx) error) error {
	tx, err := t.Begin()
	if err != nil {
		return err
	}

	if err = f(tx); err != nil {
		tx.Rollback()
		return err
	}

	// if `f` returns nil, commit the transaction
	return tx.Commit()
}

// WithFixtures takes a function `f` which requires a transaction
// and a sequence of paths to fixtures `fixs`. It first executes
// all of the fixtures on the transaction, before calling `f` on the same
// transaction. The function ends by calling rollback on the transaction.
func (t TestDatabase) WithFixtures(f func(tx db.Tx) error, fixs ...string) error {
	content, err := readAll(fixs...)
	if err != nil {
		return err
	}

	return t.RollbackTx(func(tx db.Tx) error {
		if _, err = tx.Exec(content); err != nil {
			return err
		}

		return f(tx)
	})
}

// RunFixtures execute sql fixtures on the database, within
// a transaction and then commits the changes to the database.
func (t TestDatabase) RunFixtures(fixs ...string) error {
	content, err := readAll(fixs...)
	if err != nil {
		return err
	}

	return t.CommitTx(func(tx db.Tx) (err error) {
		_, err = tx.Exec(content)
		return
	})
}

// Truncate performs a TRUNCATE query on the database for
// the provided tables, provided as a sequence of strings.
func (t TestDatabase) Truncate(tables ...string) error {
	if len(tables) == 0 {
		return fmt.Errorf("No tables provided")
	}

	query := "TRUNCATE "
	for i := 0; i < len(tables); i++ {
		query += fmt.Sprintf("%v", t.QuoteIdentifier(tables[i]))
		if i < len(tables)-1 {
			query += ", "
		}
	}

	return t.CommitTx(func(tx db.Tx) (err error) {
		_, err = tx.Exec(query)
		return
	})
}

// readAll takes a list of paths and reads each of the
// files they point to in to a single string.
func readAll(pths ...string) (string, error) {
	buf := bytes.Buffer{}
	for _, p := range pths {
		fi, err := os.Open(p)
		if err != nil {
			return buf.String(), err
		}

		if _, err = buf.ReadFrom(fi); err != nil {
			return buf.String(), err
		}
	}
	return buf.String(), nil
}

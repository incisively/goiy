package dbtest

import (
	"fmt"
	"testing"

	"github.com/ConvertHQ/goincisively/db"

	"launchpad.net/gocheck"
)

var _ = gocheck.Suite(DatabaseSuite{})

func Test(t *testing.T) { gocheck.TestingT(t) }

func Error() error { return fmt.Errorf("ERROR!") }

type DatabaseSuite struct{}

func (ds DatabaseSuite) TestRollbackTx(c *gocheck.C) {
	// Check behaves properly when error occurs on call to
	// DB.Begin.
	mdb := &MockDB{
		Beginf: func() (db.Tx, error) { return nil, Error() },
	}

	tdb := NewTestDatabase(mdb, func(s string) string { return s })

	err := tdb.RollbackTx(func(_ db.Tx) error { return nil })
	c.Check(err, gocheck.DeepEquals, Error())

	// Check behaves properly when error occurs within anonymous db
	// transaction function.
	var called bool
	mdb = rollbackTxDB(func() error {
		called = true
		return nil
	})

	tdb = NewTestDatabase(mdb, func(s string) string { return s })

	err = tdb.RollbackTx(func(_ db.Tx) error { return Error() })
	c.Check(err, gocheck.DeepEquals, Error())
	c.Check(called, gocheck.Equals, true)

	// Check rollback occurs even if all goes well
	called = false
	err = tdb.RollbackTx(func(_ db.Tx) error { return nil })
	c.Check(err, gocheck.IsNil)
	c.Check(called, gocheck.Equals, true)
}

func (ds DatabaseSuite) TestCommitTx(c *gocheck.C) {
	// Check behaves properly when error occurs on call to
	// DB.Begin()
	mdb := beginDB(func() (db.Tx, error) { return nil, Error() })

	tdb := NewTestDatabase(mdb, func(s string) string { return s })

	err := tdb.CommitTx(func(_ db.Tx) error { return nil })
	c.Check(err, gocheck.DeepEquals, Error())

	// Check behaves properly when error occurs within anonymous db
	// transaction function.
	var calledRollback bool
	mdb = rollbackTxDB(func() error {
		calledRollback = true
		return nil
	})

	tdb = NewTestDatabase(mdb, func(s string) string { return s })

	err = tdb.CommitTx(func(_ db.Tx) error { return Error() })
	c.Check(err, gocheck.DeepEquals, Error())
	c.Check(calledRollback, gocheck.Equals, true)

	// Check commit is called if all goes well
	var calledCommit, calledFunction bool
	calledRollback = false

	mdb = beginDB(func() (db.Tx, error) {
		return &MockTx{
			Rollbackf: func() error {
				calledRollback = true
				return nil
			},
			Commitf: func() error {
				calledCommit = true
				return nil
			},
		}, nil
	})

	tdb = NewTestDatabase(mdb, func(s string) string { return s })

	err = tdb.CommitTx(func(_ db.Tx) error {
		calledFunction = true
		return nil
	})

	c.Check(err, gocheck.IsNil)
	c.Check(calledRollback, gocheck.Equals, false)
	c.Check(calledFunction, gocheck.Equals, true)
	c.Check(calledCommit, gocheck.Equals, true)
}

// beginDB return a *MockDB with function `f` as
// the db.DB.Begin() function.
func beginDB(f func() (db.Tx, error)) *MockDB {
	return &MockDB{
		Beginf: f,
	}
}

// rollbaclTxDb return a mocked database, that return a mocked
// transaction on a call to begin that stubs out rollback
// with the provided function `f`.
func rollbackTxDB(f func() error) *MockDB {
	return beginDB(func() (db.Tx, error) {
		return &MockTx{
			Rollbackf: f,
		}, nil
	})
}

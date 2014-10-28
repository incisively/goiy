Database Test Package `db/dbtest`
=================================

The purpose of this package is for common database commands and database interaction mocking.

## Database Mocking

### MockDB

```go
type MockDB struct {
    ...
}
```

The `MockDB` is used for stubbing out calls to an `sql.DB` which complies with the incisively `goincisively/db.DB` interface.
The functionality of calls to a `DB` struct can be defined on the fly using anonymous functions/clojures.

When testing something that interacts with a `DB` type, for example, to make a query, the `MockDB` can be used in the following way.
```go
db := MockDB{
    Queryf: func(query string, args ...interface{}) (*sql.Rows, error) {
        fmt.Println(query, args)
        return nil, nil
    }
}

dbInteractingFunction(db)
```

Here we see a `MockDB` being created and passed to a function that interacts with a `db.DB` implementing type.
The purpose of this mock is to just print the arguments on a call to `DB.Query(...)` and return `nil` for each returned argument.

Note that for each function associated with the `db.DB` type, `MockDB` has a field that matches the type signature of these functions. 
Each field is named the same as the function in the interface, however, with an additional trailing `f`.
For example, `DB.Exec(...)` in a `MockDB` can be stubbed out by assigning to the field `MockDB.Execf`.

### MockTx

```go
type MockTx struct {
    ...
}
```

The `MockTx` type works in exactly the same fashion as a `MockDB`. It complies with the `db.Tx` interface.
Stubbing of each method can be achieved by assigning an anonymous function/clojure to the respective field.
For example, a call to `Tx.Rollback()` can be stubbed in the following way:
```go
var didRollback bool

transactionInteractingFunction(MockTx{
    Rollbackf: func() error {
        didRollback = true
        return nil
    }
})

if didRollback {
    fmt.Println(“A rollback occurred!”)
}
```
This will print the line `A rollback occurred!` if a call to `Tx.Rollback()` is made by the `transactionInteractingFunction(...)`.

## Database Interaction Testing

```go
type TestDatabase struct {
    db.DB
    QuoteIdentifier func(string) string
}
```

The TestDatabase struct embeds a db.DB implementation and extends it with common test functionality.

---

### `TestDatabase.RollbackTx(func(db.Tx) error) error`

```go
tdb.RollbackTx(func(tx db.Tx) (err error) {
    _, err = tx.Exec(“SELECT * FROM foo;”)
    return
})
```

Use this function to call a clojure that requires a database transaction `db.Tx`.
This function will always result in a call to `db.Tx.Rollback()`.

---

### `TestDatabase.CommitTx(func(db.Tx) error) error`

```go
tdb.CommitTx(func(tx db.Tx) (err error) {
    _, err = tx.Exec(“SELECT * FROM foo;”)
    return
})
```

Use this function to call a clojure that requires a database transaction `db.Tx`.
This function will result in a call to `db.Tx.Rollback()`, when the clojure returns a
non-nil error. However, if the clojure returns nil, it results in a call to `db.Tx.Commit()`


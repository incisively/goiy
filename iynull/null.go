package iynull

import (
	"database/sql"
	"encoding/json"
)

// String adds the ability to marshal and unmarshal to JSON, to a
// sql.NullString.
type String struct{ sql.NullString }

// NewString creates a new valid String from the provided string.
func NewString(s string) String {
	return String{sql.NullString{Valid: true, String: s}}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ns *String) UnmarshalJSON(data []byte) error {
	var val *string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	if val == nil {
		ns.NullString = sql.NullString{}
		return nil
	}
	ns.NullString.Valid, ns.NullString.String = true, *val
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (ns String) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

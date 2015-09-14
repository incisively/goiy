package iysql

import (
	"database/sql"
	"encoding/json"
)

// NullString adds the ability to marshal and unmarshal to JSON, to a
// sql.NullString.
type NullString struct{ *sql.NullString }

// NewNullString creates a new NullString from a string.
func NewNullString(s string) NullString {
	return NullString{&sql.NullString{Valid: true, String: s}}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ns *NullString) UnmarshalJSON(data []byte) error {
	var val *string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	if val == nil {
		ns.NullString = &sql.NullString{}
		return nil
	}
	ns.NullString = &sql.NullString{Valid: true, String: *val}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

package iysql

import (
	"database/sql"
	"encoding/json"
	"reflect"

	"testing"
)

func TestNewNullString(t *testing.T) {
	examples := []struct {
		In     string
		Valid  bool
		String string
	}{
		{In: "foo", Valid: true, String: "foo"},
		{In: "", Valid: true, String: ""},
	}

	for _, example := range examples {
		s := NewNullString(example.In)
		if s.Valid != example.Valid {
			t.Errorf("got %v, want %v", s.Valid, example.Valid)
		}

		if s.String != example.String {
			t.Errorf("got %v, want %v", s.String, example.String)
		}
	}
}

func TestNullString_Unmarshal(t *testing.T) {
	examples := []struct {
		In       string
		Expected NullString
	}{
		{In: `"foo"`, Expected: NullString{&sql.NullString{Valid: true, String: "foo"}}},
		{In: `""`, Expected: NullString{&sql.NullString{Valid: true, String: ""}}},
		{In: `null`, Expected: NullString{&sql.NullString{Valid: false, String: ""}}},
	}

	for i, example := range examples {
		var v NullString
		if err := json.Unmarshal([]byte(example.In), &v); err != nil {
			t.Errorf("[Example %d] %v", i, err)
		} else if !reflect.DeepEqual(v, example.Expected) {
			t.Errorf("[Example %d] got %#v, want %#v", i, v, example.Expected)
		}
	}
}

func TestNullString_Marshal(t *testing.T) {
	examples := []struct {
		In       NullString
		Expected string
	}{
		{In: NullString{&sql.NullString{Valid: true, String: "foo"}}, Expected: `"foo"`},
		{In: NullString{&sql.NullString{Valid: true, String: ""}}, Expected: `""`},
		{In: NullString{&sql.NullString{Valid: false, String: ""}}, Expected: `null`},
		{In: NullString{&sql.NullString{Valid: false, String: "ignored"}}, Expected: `null`},
	}

	for i, example := range examples {
		actual, err := json.Marshal(example.In)
		if err != nil {
			t.Error(err)
		}

		if string(actual) != example.Expected {
			t.Errorf("[Example %d] got %#v, want %#v", i, string(actual), example.Expected)
		}
	}
}

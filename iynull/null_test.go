package iynull

import (
	"database/sql"
	"encoding/json"
	"reflect"

	"testing"
)

func TestNewString(t *testing.T) {
	examples := []struct {
		In     string
		Valid  bool
		String string
	}{
		{In: "foo", Valid: true, String: "foo"},
		{In: "", Valid: true, String: ""},
	}

	for _, example := range examples {
		s := NewString(example.In)
		if s.Valid != example.Valid {
			t.Errorf("got %v, want %v", s.Valid, example.Valid)
		}

		if s.String != example.String {
			t.Errorf("got %v, want %v", s.String, example.String)
		}
	}
}

func TestString_Unmarshal(t *testing.T) {
	examples := []struct {
		In       string
		Expected String
	}{
		{In: `"foo"`, Expected: String{sql.NullString{Valid: true, String: "foo"}}},
		{In: `""`, Expected: String{sql.NullString{Valid: true, String: ""}}},
		{In: `null`, Expected: String{sql.NullString{Valid: false, String: ""}}},
	}

	for i, example := range examples {
		var v String
		if err := json.Unmarshal([]byte(example.In), &v); err != nil {
			t.Errorf("[Example %d] %v", i, err)
		} else if !reflect.DeepEqual(v, example.Expected) {
			t.Errorf("[Example %d] got %#v, want %#v", i, v, example.Expected)
		}
	}
}

func TestString_Marshal(t *testing.T) {
	examples := []struct {
		In       String
		Expected string
	}{
		{In: String{sql.NullString{Valid: true, String: "foo"}}, Expected: `"foo"`},
		{In: String{sql.NullString{Valid: true, String: ""}}, Expected: `""`},
		{In: String{sql.NullString{Valid: false, String: ""}}, Expected: `null`},
		{In: String{sql.NullString{Valid: false, String: "ignored"}}, Expected: `null`},
		{Expected: `null`}, // Zero value for String
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

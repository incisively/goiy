package iytime

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestDuration_UnmarshalJSON(t *testing.T) {
	examples := []struct {
		Data             string
		ExpectedErr      error
		ExpectedDuration Duration
	}{
		{Data: `"42ms"`, ExpectedDuration: Duration(42 * time.Millisecond)},
		{Data: `"2937s"`, ExpectedDuration: Duration(2937 * time.Second)},
		{Data: `"0ns"`, ExpectedDuration: Duration(0)},
		{Data: `"0m"`, ExpectedDuration: Duration(0)},
		{Data: `"3600s"`, ExpectedDuration: Duration(time.Hour)},
		{Data: `"1h"`, ExpectedDuration: Duration(time.Hour)},
		{Data: `"16m4s"`, ExpectedDuration: Duration(((16 * 60) + 4) * time.Second)},
		{Data: `"1hour"`, ExpectedErr: errors.New("time: unknown unit hour in duration 1hour")},
	}

	for i, example := range examples {
		var d Duration
		actualErr := json.Unmarshal([]byte(example.Data), &d)
		if (actualErr == nil && example.ExpectedErr != nil) ||
			(actualErr != nil && example.ExpectedErr == nil) ||
			(actualErr != nil && example.ExpectedErr != nil &&
				actualErr.Error() != example.ExpectedErr.Error()) {

			t.Errorf("example [%d] got %v, wanted %v", i+1, actualErr, example.ExpectedErr)
		}

		if d != example.ExpectedDuration {
			t.Errorf("example [%d] got %v, wanted %v", i+1, d, example.ExpectedDuration)
		}
	}
}

func TestDuration_MarshalJSON(t *testing.T) {
	examples := []struct {
		Duration       Duration
		ExpectedErr    error
		ExpectedString string
	}{
		{Duration: Duration(42 * time.Millisecond), ExpectedString: `"42ms"`},
		{Duration: Duration(2937 * time.Second), ExpectedString: `"48m57s"`},
		{Duration: Duration(0), ExpectedString: `"0"`},
		{Duration: Duration(0), ExpectedString: `"0"`},
		{Duration: Duration(time.Hour), ExpectedString: `"1h0m0s"`},
		{Duration: Duration(60 * time.Minute), ExpectedString: `"1h0m0s"`},
	}

	for i, example := range examples {
		actualByte, actualErr := json.Marshal(example.Duration)
		if (actualErr == nil && example.ExpectedErr != nil) ||
			(actualErr != nil && example.ExpectedErr == nil) ||
			(actualErr != nil && example.ExpectedErr != nil &&
				actualErr.Error() != example.ExpectedErr.Error()) {

			t.Errorf("example [%d] got %v, wanted %v", i+1, actualErr, example.ExpectedErr)
		}

		if string(actualByte) != example.ExpectedString {
			t.Errorf("example [%d] got %v, wanted %v", i+1, string(actualByte), example.ExpectedString)
		}
	}
}

func TestDuration_Hours(t *testing.T) {
	d := Duration(36000000000000)
	expected := 10.0
	if d.Hours() != expected {
		t.Errorf("got %v, wanted %v", d.Hours(), expected)
	}
}

func TestDuration_Minutes(t *testing.T) {
	d := Duration(3600000000000)
	expected := 60.0
	if d.Minutes() != expected {
		t.Errorf("got %v, wanted %v", d.Minutes(), expected)
	}
}

func TestDuration_Nanoseconds(t *testing.T) {
	d := Duration(36)
	expected := int64(36)
	if d.Nanoseconds() != expected {
		t.Errorf("got %v, wanted %v", d.Nanoseconds(), expected)
	}
}

func TestDuration_Seconds(t *testing.T) {
	d := Duration(360000000000)
	expected := 360.0
	if d.Seconds() != expected {
		t.Errorf("got %v, wanted %v", d.Seconds(), expected)
	}
}

func TestDuration_String(t *testing.T) {
	d := Duration(3600000)
	expected := "3.6ms"
	if d.String() != expected {
		t.Errorf("got %q, wanted %q", d.String(), expected)
	}
}

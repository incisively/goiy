package iyio

import (
	"errors"
	"reflect"
	"strings"

	"testing"
)

func TestMockReader_Read(t *testing.T) {
	r := NewMockReader(nil)
	r.readError = errors.New("a Read error")

	// It returns an error if one is set.
	_, actual := r.Read(nil)
	if !reflect.DeepEqual(r.readError, actual) {
		t.Errorf("expected %v got %v", r.readError, actual)
	}

	// It calls the underlying reader.
	r = NewMockReader(strings.NewReader("abc"))
	buf := make([]byte, 3)
	n, err := r.Read(buf)
	if !reflect.DeepEqual(nil, err) {
		t.Errorf("expected %v got %v", nil, err)
	}
	if n != 3 {
		t.Errorf("expected %v got %v", 3, n)
	}

	if !reflect.DeepEqual([]byte("abc"), buf) {
		t.Errorf("expected %v got %v", []byte("abc"), buf)
	}
}

func TestMockReader_SetReadError(t *testing.T) {
	r := NewMockReader(nil)
	expected := errors.New("a Read error")
	r.SetReadError(expected)
	if !reflect.DeepEqual(expected, r.readError) {
		t.Errorf("expected %v got %v", expected, r.readError)
	}
}

func TestMockReadCloser_Close(t *testing.T) {
	r := NewMockReadCloser(nil)
	expected := errors.New("a Close error")
	r.SetCloseError(expected)

	// It returns the error set by SetCloseError.
	actual := r.Close()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v got %v", expected, actual)
	}
}

func TestMockReadCloser_Closed(t *testing.T) {
	r := NewMockReadCloser(nil)
	if c := r.Closed(); c {
		t.Errorf("expected %v got %v", false, c)
	}

	r.Close()
	if c := r.Closed(); !c {
		t.Errorf("expected %v got %v", true, c)
	}
}

func TestMockWriter_Write(t *testing.T) {
	var p []byte
	input := [][]byte{
		nil,
		[]byte("foo"),
		p,
		[]byte{},
	}

	var w MockWriter
	for _, p := range input {
		w.Write(p)
	}

	// It stores copies of all the arguments passed into Write.
	if !reflect.DeepEqual(input, w.writeArgs) {
		t.Errorf("expected %v got %v", input, w.writeArgs)
	}
}

func TestMockWriter_Reset(t *testing.T) {
	var w MockWriter
	w.SetWriteN(10)
	w.SetWriteError(errors.New("error"))
	w.Write([]byte("foo"))
	w.WriteCalledArgs()

	w.Reset()
	if !reflect.DeepEqual(w, MockWriter{}) {
		t.Errorf("expected %v got %v", MockWriter{}, w)
	}
}

func TestMockWriter_WriteCalledArgs(t *testing.T) {
	var p []byte
	input := [][]byte{
		nil,
		[]byte("foo"),
		p,
		[]byte{},
	}

	var w MockWriter
	for _, p := range input {
		w.Write(p)

		// It returns a copy of the argument passed into Write.
		next := w.WriteCalledArgs()
		if !reflect.DeepEqual(p, next) {
			t.Errorf("expected %v got %v", p, next)
		}
	}

	// It panics if it's called one more time.
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic.")
		}
	}()
	w.WriteCalledArgs()
}

func TestMockWriter_WriteCalledN(t *testing.T) {
	var w MockWriter
	w.Write(nil)
	w.Write(nil)

	if w.WriteCalledN() != 2 {
		t.Errorf("expected %v got %v", 2, w.WriteCalledN())
	}
}

func TestMockWriter_SetWriteN(t *testing.T) {
	var w MockWriter
	w.SetWriteN(20)
	if n, _ := w.Write(nil); n != 20 {
		t.Errorf("expected %v got %v", 20, n)
	}
}

func TestMockWriter_SetWriteError(t *testing.T) {
	var w MockWriter
	expected := errors.New("an error")
	w.SetWriteError(expected)
	_, actual := w.Write(nil)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v got %v", expected, actual)
	}
}

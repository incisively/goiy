package iyio

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"testing"
)

func ExampleMockReader() {
	// MockReader makes it simple to force an error from an io.Reader.
	r := NewMockReader(nil)
	r.SetReadError(errors.New("a Read error"))
	_, err := io.Copy(os.Stdout, r)
	fmt.Println(err)

	fmt.Println()

	// Alternatively you can use it to track calls to Read.
	r = NewMockReader(strings.NewReader("hello world\n"))
	fmt.Println(io.Copy(os.Stdout, r))
	fmt.Printf("MockReader called %d times.", r.ReadCalledN())

	// Output: a Read error
	//
	// hello world
	// 12 <nil>
	// MockReader called 2 times.
}

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

func ExampleMockReadCloser() {
	// MockReadCloser makes it simple to force an error on Close.
	r := NewMockReadCloser(nil)
	r.SetCloseError(errors.New("a Close error"))
	fmt.Println(r.Close())
	fmt.Println(r.CloseCalled())

	// MockReadCloser embeds a MockReader, so all that functionality
	// is available too.
	r.SetReadError(errors.New("a Read error"))
	fmt.Println(io.Copy(os.Stdout, r))

	// Output: a Close error
	// true
	// 0 a Read error
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

func TestMockReadCloser_CloseCalled(t *testing.T) {
	r := NewMockReadCloser(nil)
	if c := r.CloseCalled(); c {
		t.Errorf("expected %v got %v", false, c)
	}

	r.Close()
	if c := r.CloseCalled(); !c {
		t.Errorf("expected %v got %v", true, c)
	}
}

func ExampleMockWriter() {
	// MockWriter makes it simple to force an error on Close.
	var w MockWriter

	w.SetWriteError(errors.New("a Write error"))

	n, err := fmt.Fprint(&w, "foo")
	fmt.Println(n, err)
	fmt.Printf("Write was called %d time.\n", w.WriteCalledN())

	fmt.Println()

	// Reset the MockWriter
	w.Reset()
	fmt.Printf("Write was called %d times.\n", w.WriteCalledN())

	n, err = fmt.Fprint(&w, "hello world")
	fmt.Println(n, err)

	fmt.Printf("Write was called %d time.\n", w.WriteCalledN())
	fmt.Println(string(w.WriteCalledArgs()))
	// w.WriteCalledArgs() - would panic.

	// Output: 3 a Write error
	// Write was called 1 time.
	//
	// Write was called 0 times.
	// 11 <nil>
	// Write was called 1 time.
	// hello world
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

func TestMockWriter_SetWriteError(t *testing.T) {
	var w MockWriter
	expected := errors.New("an error")
	w.SetWriteError(expected)
	_, actual := w.Write(nil)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %v got %v", expected, actual)
	}
}

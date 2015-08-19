// Package iyio provides convenience mock implementations of io.Readers,
// io.ReadClosers and io.Writers.
package iyio

import (
	"sync"

	"io"
)

// MockReader provides a mock implementation of an io.Reader, allowing
// the caller to override the error returned by Read.
//
// An underlying io.Reader can optionally be provided, in which case any
// calls to Read will call Read on the underlying reader. In this way,
// MockReader can be used as a shim, to track calls to Read on the
// underlying io.Reader.
//
// MockReader is safe for use by multiple goroutines.
type MockReader struct {
	r io.Reader

	mu          sync.Mutex
	readError   error
	readCalledN int
}

// NewMockReader returns a new MockReader that will wrap r. It is OK to
// pass nil in; the reader will simply never attempt to read anything
// into slices passed to Read.
func NewMockReader(r io.Reader) *MockReader {
	return &MockReader{r: r}
}

// Read first reads from the underlying reader into b. Then, if the
// MockReader has an error set for Read, the error value from the
// underlying reader is overwritten.
func (r *MockReader) Read(b []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.readCalledN++

	if r.r != nil {
		n, err = r.r.Read(b)
	}

	if r.readError != nil {
		err = r.readError
	}
	return n, err
}

// ReadCalledN returns the number of times Read has been called on
// the MockReader.
func (r *MockReader) ReadCalledN() int {
	return r.readCalledN
}

// SetReadError sets the error that Read will return.
func (r *MockReader) SetReadError(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.readError = err
}

// MockReadCloser turns an io.Reader into an io.ReadCloser.
//
// Using a MockReadCloser allows you to mock out errors from calls to
// Read or Close.
//
// MockReadCloser is safe for use by multiple goroutines.
type MockReadCloser struct {
	*MockReader

	mu         sync.Mutex
	closed     bool
	closeError error
}

// NewMockReadCloser returns a MockReadCloser that wraps r.
func NewMockReadCloser(r io.Reader) *MockReadCloser {
	return &MockReadCloser{MockReader: NewMockReader(r)}
}

// Close closes the MockReadCloser, returning the CloseError value.
func (r *MockReadCloser) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.closed = true
	return r.closeError
}

// SetCloseError sets the error that Close will return.
func (r *MockReadCloser) SetCloseError(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.closeError = err
}

// CloseCalled returns true if Close has been called.
func (r *MockReadCloser) CloseCalled() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.closed
}

// MockWriter provides a mock implementation of an io.Writer.
//
// The returned error for the Write method can be set for testing
// purposes, and each individual call to Write is stored for later
// inspection.
//
// MockWriter is safe for use by multiple goroutines.
type MockWriter struct {
	mu        sync.Mutex
	writeArgs [][]byte
	writeI    int

	writeError error
}

// Reset resets the state of the MockWriter.
func (w *MockWriter) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.writeArgs = nil
	w.writeI = 0
	w.writeError = nil
}

// Write stores a copy of p for later inspection. The returned values
// can be set using SetWriteN and SetWriteError.
func (w *MockWriter) Write(p []byte) (n int, err error) {
	var pp []byte
	w.mu.Lock()
	defer w.mu.Unlock()
	if p != nil {
		pp = make([]byte, len(p))
		n = copy(pp, p)
	}
	w.writeArgs = append(w.writeArgs, pp)
	return n, w.writeError
}

// WriteCalledArgs returns, in order, the arguments passed into calls to
// Write.
//
// WriteCalledArgs panics if there are no more call arguments left to
// return.
func (w *MockWriter) WriteCalledArgs() []byte {
	if len(w.writeArgs) <= w.writeI {
		panic("Write was not called.")
	}

	defer func() { w.writeI++ }()
	return w.writeArgs[w.writeI]
}

// WriteCalledN returns the number of times Write has been called on
// the MockWriter.
func (w *MockWriter) WriteCalledN() int {
	return len(w.writeArgs)
}

// SetWriteError sets the error returned by Write.
func (w *MockWriter) SetWriteError(err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.writeError = err
}

package iyhttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"testing"
)

func TestNewResponseWriterShim(t *testing.T) {
	in := httptest.NewRecorder()
	shim := NewResponseWriterShim(in)

	// It uses the same map as the original http.ResponseWriter
	in.Header().Add("foo", "bar")
	if shim.Header().Get("foo") != "bar" {
		t.Errorf("got %q, expected %q", shim.Header().Get("foo"), "bar")
	}
}

func TestResponseWriterShim_Write(t *testing.T) {
	in := httptest.NewRecorder()
	shim := NewResponseWriterShim(in)

	p := []byte("hello")
	n, err := shim.Write(p)
	if err != nil {
		t.Error(err)
	}

	if n != len(p) {
		t.Errorf("wrote %d, expected %d bytes", n, len(p))
	}

	// It writes the data to the underlying recorder.
	recData := shim.Body.String()
	if recData != string(p) {
		t.Errorf("got %q, expected %q", recData, string(p))
	}

	// It writes the data to the underlying http.ResponseWriter.
	rwData := in.Body.String()
	if rwData != string(p) {
		t.Errorf("got %q, expected %q", rwData, string(p))
	}
}

func TestResponseWriterShim_WriteHeader(t *testing.T) {
	in := httptest.NewRecorder()
	shim := NewResponseWriterShim(in)

	shim.WriteHeader(http.StatusInternalServerError)

	// It writes the status to the underlying recorder.
	recCode := shim.ResponseRecorder.Code
	if recCode != http.StatusInternalServerError {
		t.Errorf("got %d, expected %d", recCode, http.StatusInternalServerError)
	}

	// It writes the status to the underlying http.ResponseWriter.
	rwCode := in.Code
	if rwCode != in.Code {
		t.Errorf("got %q, expected %q", rwCode, in.Code)
	}
}

func TestResponseWriterShim_Dump(t *testing.T) {
	in := httptest.NewRecorder()
	shim := NewResponseWriterShim(in)

	shim.WriteHeader(http.StatusInternalServerError)
	shim.Header().Add("Content-Type", "application/json")
	shim.Header().Add("Content-Type", "foo")
	fmt.Fprint(shim, `{"hello": "world"}`)

	actual := shim.Dump()
	expected := `Content-Type: application/json,foo
{"hello": "world"}`

	if actual != expected {
		t.Errorf("got %s, expected %s", actual, expected)
	}
}

package iyhttp

import (
	"net/http"
	"net/http/httptest"
)

// ResponseWriterShim shims an http.ResponseWriter, capturing all
// modifications of headers, status code, or body.
//
// ResponseWriterShim is useful when you're passing an
// http.ResponseWriter into another API from your handler, and you want
// to capture any modifications to the response for later logging or
// processing.
//
// ResponseWriterShim implements the http.ResponseWriter interface.
//
// Example:
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//		shim := iyhttp.NewResponseWriterShim(w)
//
//		// DoSomething expects an http.ResponseWriter.
//		otherPkg.DoSomething(shim)
//
//		// Now we have insight into what DoSomething did with the
//		// http.ResponseWriter.
//		log.Println(shim.Dump())
//	}
type ResponseWriterShim struct {
	*httptest.ResponseRecorder
	errRec error
	w      http.ResponseWriter
}

// NewResponseWriterShim returns an initialialised shim.
func NewResponseWriterShim(w http.ResponseWriter) *ResponseWriterShim {
	rec := httptest.NewRecorder()
	rec.HeaderMap = w.Header()
	return &ResponseWriterShim{ResponseRecorder: rec, w: w}
}

// Header returns the header map that will be sent by WriteHeader.
// See net/http documentation for more information.
func (r *ResponseWriterShim) Header() http.Header {
	return r.ResponseRecorder.Header()
}

// Write writes the data to the connection as part of an HTTP reply.
// See net/http documentation for more information.
func (r *ResponseWriterShim) Write(p []byte) (int, error) {
	r.ResponseRecorder.Write(p)
	return r.w.Write(p)
}

// WriteHeader sends an HTTP response header with status code.
// See net/http documentation for more information.
func (r *ResponseWriterShim) WriteHeader(i int) {
	r.ResponseRecorder.WriteHeader(i)
	r.w.WriteHeader(i)
}

// Dump returns the captured response headers and body.
func (r *ResponseWriterShim) Dump() string {
	var data string
	for k, v := range r.Header() {
		data += k + ": "
		for i := 0; i < len(v); i++ {
			data += v[i]
			if i < len(v)-1 {
				data += ","
			}
		}
		data += "\n"
	}
	return data + r.Body.String()
}

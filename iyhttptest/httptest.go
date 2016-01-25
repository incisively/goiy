package iyhttptest

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/kr/pretty"
)

// A Server is an HTTP server listening on a system-chosen port on the
// local loopback interface, for use in end-to-end HTTP tests.
type Server struct{ *httptest.Server }

// NewServer returns a newly instantiated Server.
func NewServer(handler http.Handler) *Server {
	return &Server{
		Server: httptest.NewServer(handler),
	}
}

func (s *Server) fullURL(urlStr string, values url.Values) string {
	if len(urlStr) > 0 && !strings.HasPrefix(urlStr, "/") {
		urlStr = "/" + urlStr
	}

	// add optional query string parameters
	var qs string
	if values != nil && len(values) > 0 {
		qs = "?" + values.Encode()
	}
	return s.URL + urlStr + qs
}

// DoRequest sends an HTTP request to the test server at the urlStr
// endpoint, with the provided headers, query parameters and body.
//
// DoRequest currently uses http.DefaultClient to make the request, and
// panics if there is an error.
func (s *Server) DoRequest(method, urlStr string, headers map[string][]string, body io.Reader, values url.Values) Response {
	req, err := http.NewRequest(method, s.fullURL(urlStr, values), body)
	if err != nil {
		panic(err)
	}

	for name, vals := range headers {
		for _, val := range vals {
			req.Header.Add(name, val)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	return Response{Response: resp}
}

// Get is a shortcut for using the GET method with DoRequest.
func (s *Server) Get(urlStr string, headers map[string][]string, values url.Values) Response {
	return s.DoRequest("GET", urlStr, headers, nil, values)
}

// post is a shortcut for using the POST method with DoRequest.
func (s *Server) Post(urlStr string, headers map[string][]string, body io.Reader) Response {
	return s.DoRequest("POST", urlStr, headers, body, nil)
}

// Patch is a shortcut for using the PATCH method with DoRequest.
func (s *Server) Patch(urlStr string, headers map[string][]string, body io.Reader) Response {
	return s.DoRequest("PATCH", urlStr, headers, body, nil)
}

// DoJSON encodes the provided value into a JSON string, and then sends
// it to the test server using DoRequest.
func (s *Server) DoJSON(urlStr string, method string, headers map[string][]string, v interface{}) Response {
	// marshal v into JSON
	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(v); err != nil {
		panic(err)
	}
	return s.DoRequest(method, urlStr, headers, b, nil)
}

// PostJSON is a shortcut for using the POST method with DoJSON.
func (s *Server) PostJSON(urlStr string, headers map[string][]string, v interface{}) Response {
	return s.DoJSON(urlStr, "POST", headers, v)
}

// PatchJSON is a shortcut for using the PATCH method with DoJSON.
func (s *Server) PatchJSON(urlStr string, headers map[string][]string, v interface{}) Response {
	return s.DoJSON(urlStr, "PATCH", headers, v)
}

// PutJSON is a shortcut for using the PUT method with DoJSON.
func (s *Server) PutJSON(urlStr string, headers map[string][]string, v interface{}) Response {
	return s.DoJSON(urlStr, "PUT", headers, v)
}

// PostForm encodes the provided map into URL encoded form values and
// makes a request to the test server using DoRequest.
func (s *Server) PostForm(urlStr string, headers map[string][]string, form map[string][]string) Response {
	values := url.Values(form)
	b := bytes.NewBufferString(values.Encode())
	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = []string{"application/x-www-form-urlencoded"}
	}
	return s.DoRequest("POST", urlStr, headers, b, nil)
}

// Response is a http.Response object with some extra convenience
// methods. It is specifically for use with testing, and most methods
// accept a *testing.T value.
type Response struct{ *http.Response }

// NewResponse returns a newly instantiated Response.
func NewResponse(r *http.Response) *Response {
	return &Response{
		Response: r,
	}
}

// StatusIs checks if the HTTP status of the Response is equal to code.
// If it differs, StatusIs dumps a stacktrace and the response body.
func (r *Response) StatusIs(t *testing.T, code int) {
	if r.StatusCode != code {
		trace := make([]byte, 1024)
		runtime.Stack(trace, false)
		t.Fatalf("status %v is not %v\nbody: %v\n%v\n", r.StatusCode, code, r.BodyString(), string(trace))
	}
}

// HasHeader checks if the HTTP status of the Response contains all of
// the provided headers.
func (r *Response) HasHeader(t *testing.T, name string, values ...string) {
	rv, ok := r.Header[name]
	if !ok {
		t.Fatalf("Header %q missing from response\n", name)
	}

	sort.Strings(rv)
	sort.Strings(values)
	if !reflect.DeepEqual(rv, values) {
		t.Fatalf("expected %v, got %v\n", values, rv)
	}
}

// BodyString returns a string representation of the Response's body.
func (r *Response) BodyString() string {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}
	return string(b)
}

// JSONBodyEquals checks that the body of the response contains the
// values provided in expected.
func (r Response) JSONBodyEquals(t *testing.T, dest, expected interface{}) {
	// unmarshal into dest
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(b, dest); err != nil {
		trace := make([]byte, 1024)
		runtime.Stack(trace, false)
		t.Fatalf("Error unmarshaling response body.\nerror: %v\nbody: %v\ntrace: %v\n", err, string(b), string(trace))
	}

	if !reflect.DeepEqual(dest, expected) {
		trace := make([]byte, 1024)
		runtime.Stack(trace, false)

		format := "expected %# v\ngot %# v\ndiff: %v\nbody: %v\ntrace: %v\n"
		t.Fatalf(format,
			pretty.Formatter(expected),
			pretty.Formatter(dest),
			pretty.Diff(expected, dest),
			string(b),
			string(trace),
		)
	}
}

func (r Response) JSONBodyOnlyKeys(t *testing.T, keys []string) {
	// unmarshal body
	actual := map[string]interface{}{}
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(b, &actual); err != nil {
		t.Fatalf("error: %v\nbody: %v\n", err, string(b))
	}

	if len(actual) != len(keys) {
		t.Fatalf("expected %v keys in body: %v\n", keys, string(b))
	}

	for _, k := range keys {
		if _, ok := actual[k]; !ok {
			t.Fatalf("key %q not found\nbody: %v \n", k, string(b))
		}
	}
}

// NewGetRequest returns a new *http.Request, panicking if an error is
// returned.
func NewGetRequest(url string) *http.Request {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	r.RequestURI = url
	return r
}

// NewPostRequest returns a new *http.Request, panicking if an error is
// returned.
func NewPostRequest(url string, body io.Reader) *http.Request {
	r, err := http.NewRequest("POST", url, body)
	if err != nil {
		panic(err)
	}
	return r
}

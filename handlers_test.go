package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmorgan1321/SpaceRep/internal/test"
)

type TestServer struct {
}

func (ts TestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func TestServerCardIndexHandlerDoesX(t *testing.T) {
	// test setup
	expectedBody := "Hello"
	handler := TestServer{}
	recorder := httptest.NewRecorder()

	// create the request for the server to process
	url := "http://localhost.com/api/vi/cards/"
	req, err := http.NewRequest("GET", url, nil)
	test.Assert(t, err == nil, "unexpected error")

	// test that server handles request properly
	handler.ServeHTTP(recorder, req)
	test.AssertEQ(t, expectedBody, recorder.Body.String(), "handler response mismatch")
}

func TestServerIndexReturnsStatusOk(t *testing.T) {
	// test setup
	handler := TestServer{}
	recorder := httptest.NewRecorder()

	// create the request for the server to process
	url := "http://localhost.com:8080"
	req, err := http.NewRequest("GET", url, nil)
	test.Assert(t, err == nil, "unexpected error")

	// test that server handles request properly
	handler.ServeHTTP(recorder, req)
	test.ExpectEQ(t, recorder.Code, http.StatusOK, "code mismatch")
}

// func TestClientFuncDoesYWithX(t *testing.T) {
// 	// test setup
// 	expectedBody := "Hello"
// 	server := httptest.NewServer(TestServer{})
// 	defer server.Close()

// 	// create client request to server
// 	url := "http://localhost.com/api/vi/cards/"
// 	resp, err := http.DefaultClient.Get(url)
// 	test.Assert(t, err == nil, "unexpected error")

// 	// test client response to server data
// 	// ie, normally I'd do something with the "Hello" like test that I could
// 	// count the num characters in it or something.
// 	b, err := ioutil.ReadAll(io.LimitReader(resp.Body, 2^20))
// 	test.Assert(t, err == nil, "unexpected error")
// 	test.AssertEQ(t, expectedBody, string(b), "handler response mismatch")
// }

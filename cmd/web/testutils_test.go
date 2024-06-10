package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// returns an instance of our application struct containing mocked
// dependencies.
func newTestApplication(_ *testing.T) *application {
	return &application{
		infoLog:  log.New(io.Discard, "", 0),
		errorLog: log.New(io.Discard, "", 0),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(_ *testing.T, h http.Handler) *testServer {
	return &testServer{httptest.NewTLSServer(h)}
}

// get makes a GET request to a given url path using the test server client,
// and returns the response status code, headers and body.
func (ts *testServer) get(t *testing.T, urlPath string) (
	int,
	http.Header,
	string,
) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)
}

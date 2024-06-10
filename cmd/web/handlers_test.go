package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eXvimmer/lets_go/internal/assert"
)

func TestPing(t *testing.T) {
	app := application{
		infoLog:  log.New(io.Discard, "INFO:\t", 0),
		errorLog: log.New(io.Discard, "INFO:\t", 0),
	}
	ts := httptest.NewTLSServer(app.routes())
	defer ts.Close()
	rs, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, rs.StatusCode, http.StatusOK)
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}

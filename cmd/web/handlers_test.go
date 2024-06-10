package main

import (
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/eXvimmer/lets_go/internal/assert"
)

func TestPing(t *testing.T) {
	app := application{
		infoLog:  log.New(io.Discard, "INFO:\t", 0),
		errorLog: log.New(io.Discard, "INFO:\t", 0),
	}
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	code, _, body := ts.get(t, "/ping")
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}

package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tm := time.Date(2024, 6, 10, 14, 54, 0, 0, time.UTC)
	hd := humanDate(tm)
	want := "10 Jun 2024 at 14:54:00"
	if hd != want {
		t.Errorf("got %q; want %q", hd, want)
	}
}

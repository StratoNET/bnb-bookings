package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler

	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		// this would be correct, so do nothing
	default:
		t.Errorf("type is NOT http.Handler, it is %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler

	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		// this would be correct, so do nothing
	default:
		t.Errorf("type is NOT http.Handler, it is %T", v)
	}
}

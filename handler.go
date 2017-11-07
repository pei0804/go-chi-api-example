package main

import (
	"io"
	"net/http"

	"github.com/pei0804/goapi/httputil"
)

type handler func(w http.ResponseWriter, r *http.Request) error

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	runHandler(w, r, h, errHandler)
}

type errFunc func(w http.ResponseWriter, r *http.Request, status int, err error)

func errHandler(w http.ResponseWriter, r *http.Request, status int, err error) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	io.WriteString(w, err.Error())
}

func runHandler(w http.ResponseWriter, r *http.Request,
	fn func(w http.ResponseWriter, r *http.Request) error, errFunc errFunc) {
	err := fn(w, r)
	if e, ok := err.(*httputil.HTTPError); ok {
		errFunc(w, r, e.Status, err)
	} else if err != nil {
		errFunc(w, r, http.StatusInternalServerError, err)
	}
}

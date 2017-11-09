package httputil

import (
	"io"
	"net/http"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	RunHandler(w, r, h, ErrHandler)
}

type ErrFunc func(w http.ResponseWriter, r *http.Request, status int, err error)

func ErrHandler(w http.ResponseWriter, r *http.Request, status int, err error) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	io.WriteString(w, err.Error())
}

func RunHandler(w http.ResponseWriter, r *http.Request,
	fn func(w http.ResponseWriter, r *http.Request) error, errFunc ErrFunc) {
	err := fn(w, r)
	if e, ok := err.(*HTTPError); ok {
		errFunc(w, r, e.Status, err)
	} else if err != nil {
		errFunc(w, r, http.StatusInternalServerError, err)
	}
}

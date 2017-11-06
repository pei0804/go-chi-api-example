package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/zenazn/goji"
)

func main() {
	member := NewMember()
	goji.Get("/", handler(member.List))
	goji.Serve()
}

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

func NewMember() *Member {
	return &Member{}
}

type Member struct {
}

func (m *Member) List(w http.ResponseWriter, r *http.Request) error {
	type j struct {
		ID int
	}
	js := j{ID: 1}
	return JSON(w, http.StatusOK, js)
}

func JSON(w http.ResponseWriter, status int, i interface{}) (err error) {
	b, err := json.Marshal(i)
	if err != nil {
		return
	}
	return JSONBlob(w, status, b)
}

func JSONBlob(w http.ResponseWriter, status int, b []byte) (err error) {
	return Blob(w, status, "application/json", b)
}

func Blob(w http.ResponseWriter, status int, contentType string, b []byte) (err error) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	_, err = w.Write(b)
	return
}

type HTTPError struct {
	Status  int
	Message string
}

func (he *HTTPError) Error() string {
	return fmt.Sprintf("status=%d, message=%v", he.Status, he.Message)
}

func runHandler(w http.ResponseWriter, r *http.Request,
	fn func(w http.ResponseWriter, r *http.Request) error, errFunc errFunc) {
	err := fn(w, r)
	if e, ok := err.(*HTTPError); ok {
		errFunc(w, r, e.Status, err)
	} else if err != nil {
		errFunc(w, r, http.StatusInternalServerError, err)
	}
}

func NewHTTPError(status int, message ...interface{}) *HTTPError {
	he := &HTTPError{Status: status}
	return he
}

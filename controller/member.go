package controller

import (
	"net/http"

	"github.com/pei0804/goapi/httputil"
)

func NewMember() *Member {
	return &Member{}
}

type Member struct {
}

func (m *Member) List(w http.ResponseWriter, r *http.Request) error {
	type j struct {
		ID   int
		Name string
	}
	js := j{ID: 1, Name: "hoge"}
	return httputil.JSONPretty(w, http.StatusOK, js)
}

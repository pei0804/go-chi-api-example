package controller

import (
	"github.com/pei0804/goapi/httputil"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func NewRouter() *web.Mux {
	m := goji.DefaultMux
	member := NewMember()
	m.Get("/", httputil.Handler(member.List))
	return m
}

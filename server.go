package main

import (
	"github.com/pei0804/goapi/controller"
	"github.com/zenazn/goji"
)

func main() {
	member := controller.NewMember()
	goji.Get("/", handler(member.List))
	goji.Serve()
}

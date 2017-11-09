package main

import (
	"github.com/pei0804/goapi/controller"
	"github.com/zenazn/goji"
)

func main() {
	_ = controller.NewRouter()
	goji.Serve()
}

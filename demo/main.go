package main

import (
	"fmt"
	"net/http"

	"github.com/lucasmenendez/framegowork"
)

func echo(w http.ResponseWriter, r *http.Request, params framegowork.Params) {
	fmt.Fprintf(w, "Hello, "+params["msg"])
}

func middleware(w http.ResponseWriter, r *http.Request, next framegowork.NextHandler) {
	fmt.Fprintf(w, "Hey!\n")
	next.Exec(w, r)
}

func main() {
	framegowork := framegowork.New()

	framegowork.SetPort(9999)
	framegowork.GET("/echo/:msg", echo, middleware)
	framegowork.Run()
}

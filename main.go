package main

import (
	"fmt"
	"framegowork/router"
	"net/http"
)

func echo(w http.ResponseWriter, r *http.Request, params map[string]string) {
	fmt.Fprintf(w, "Hello, "+params["msg"])
}

func middleware(w http.ResponseWriter, r *http.Request, params map[string]string, next router.NextHandler) {
	fmt.Fprintf(w, "Hey!\n")
	next.Exec(w, r, params)
}

func main() {
	router := router.New()
	router.GET("/echo/:msg", echo, middleware)
	router.Run("9999")
}

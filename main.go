package main

import (
	"fmt"
	"framegowork/router"
	"net/http"
)

func echo(w http.ResponseWriter, r *http.Request, params map[string]string) {
	fmt.Fprintf(w, "Hello, "+params["msg"])
}

func main() {
	router := router.New()
	router.GET("/echo/:msg", echo)
	router.Run("9999")
}

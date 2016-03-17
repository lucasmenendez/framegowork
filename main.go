package main

import (
	"fmt"
	"framegowork/server"
	"net/http"
)

func echo(w http.ResponseWriter, r *http.Request, params map[string]string) {
	fmt.Fprintf(w, "Hello, "+params["msg"])
}

func middleware(w http.ResponseWriter, r *http.Request, next server.NextHandler) {
	fmt.Fprintf(w, "Hey!\n")
	next.Exec(w, r)
}

func main() {
	server := server.New()

	server.SetPort("9999")
	server.SetHeader("Content-Type", "text/plain")
	server.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	server.SetHeader("Access-Control-Allow-Origin", "*")
	server.SetHeader("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	server.GET("/echo/:msg", echo, middleware)
	server.Run()
}

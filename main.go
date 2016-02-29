package main

import (
	"fmt"
	"net/http"
	"framegowork/router"
)

func get(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GET METHOD")
}

func post(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "POST METHOD")
}

func put(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PUT METHOD")
}

func del(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "DELETE METHOD")
}

func main () {
	router := router.New()
	router.GET("/method", get)
	router.POST("/method", post)
	router.PUT("/method", put)
	router.DELETE("/method", del)
	router.RunServer("9999")
}

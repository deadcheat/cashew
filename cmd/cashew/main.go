package main

import (
	"log"
	"net/http"

	dh "github.com/deadcheat/cashew/deliver/http"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	d := dh.New(r)
	d.Mount()
	log.Println("start cas server")
	log.Fatal(http.ListenAndServe(":3000", r))
}

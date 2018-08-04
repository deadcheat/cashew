package main

import (
	"log"
	"net/http"

	dh "github.com/deadcheat/cashew/deliver/http"
	"github.com/deadcheat/cashew/foundation"
	"github.com/gorilla/mux"
)

func main() {
	_ = foundation.LoadConfig()
	r := mux.NewRouter()
	d := dh.New(r)
	d.Mount()
	log.Println("start cas server")
	log.Fatal(http.ListenAndServe(":3000", r))
}

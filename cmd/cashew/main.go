package main

import (
	"log"
	"net/http"

	ld "github.com/deadcheat/cashew/deliver/login"
	"github.com/deadcheat/cashew/foundation"
	"github.com/gorilla/mux"
)

func main() {
	// load global setting
	l := foundation.GlobalSettingLoader()
	configFile := "config.yml"
	_, err := l.Load(configFile)
	if err != nil {
		log.Fatalf("failed to load config file %s \n", configFile)
	}

	// create usecase, repository, deliver and mount them
	r := mux.NewRouter()
	login := ld.New(r)
	login.Mount()

	// start cas server
	log.Println("start cas server")
	log.Fatal(http.ListenAndServe(":3000", r))
}

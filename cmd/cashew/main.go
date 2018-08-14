package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/deadcheat/cashew/repository/ticket"
	"github.com/deadcheat/cashew/usecase/login"

	ld "github.com/deadcheat/cashew/deliver/login"
	"github.com/deadcheat/cashew/foundation"
	"github.com/gorilla/mux"
)

func main() {
	var err error
	// load setting
	var configFile string
	flag.StringVar(&configFile, "c", "config.yml", "specify config file path")
	flag.Parse()
	err = foundation.Load(configFile)
	if err != nil {
		log.Fatalf("failed to load config file %s \n", configFile)
	}
	// start database
	err = foundation.StartDatabase()
	if err != nil {
		log.Fatalf("failed to start database %+v \n", err)
	}

	// prepare authenticator
	err = foundation.PrepareAuthenticator()
	if err != nil {
		log.Fatalf("failed to prepare authenticator %+v \n", err)
	}

	// create router
	r := mux.NewRouter()

	// create usecase, repository, deliver and mount them
	ticketRepository := ticket.New(foundation.DB())
	loginUseCase := login.New(ticketRepository)
	login := ld.New(r, loginUseCase)
	login.Mount()

	// start cas server
	bindAddress := fmt.Sprintf("%s:%d", foundation.App().Host, foundation.App().Port)
	log.Println("start cas server on ", bindAddress)
	log.Fatal(http.ListenAndServe(bindAddress, r))
}

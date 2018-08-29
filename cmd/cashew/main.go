package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/deadcheat/cashew/deliver/assets"
	dl "github.com/deadcheat/cashew/deliver/login"
	dv "github.com/deadcheat/cashew/deliver/validate"
	"github.com/deadcheat/cashew/foundation"
	"github.com/deadcheat/cashew/repository/host"
	"github.com/deadcheat/cashew/repository/id"
	"github.com/deadcheat/cashew/repository/ticket"
	"github.com/deadcheat/cashew/usecase/auth"
	"github.com/deadcheat/cashew/usecase/login"
	"github.com/deadcheat/cashew/usecase/logout"
	"github.com/deadcheat/cashew/usecase/validate"

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
	r := mux.NewRouter().PathPrefix(foundation.App().URIPath).Subrouter()

	// create usecase, repository, deliver and mount them
	ticketRepository := ticket.New(foundation.DB())
	idRepository := id.New()
	hostRepository := host.New()
	loginUseCase := login.New(ticketRepository, idRepository, hostRepository)
	logoutUseCase := logout.New(ticketRepository)
	authUseCase := auth.New()
	login := dl.New(r, loginUseCase, logoutUseCase, authUseCase)
	login.Mount()
	validateUseCase := validate.New(ticketRepository)
	v := dv.New(r, validateUseCase)
	v.Mount()

	// mount to static files
	statics := assets.New(r)
	statics.Mount()

	// start cas server
	bindAddress := fmt.Sprintf("%s:%d", foundation.App().Host, foundation.App().Port)
	log.Println("start cas server on ", bindAddress)
	log.Fatal(http.ListenAndServe(bindAddress, r))
}

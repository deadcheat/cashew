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
	"github.com/deadcheat/cashew/repository/proxycallback"
	tr "github.com/deadcheat/cashew/repository/ticket"
	"github.com/deadcheat/cashew/usecase/auth"
	"github.com/deadcheat/cashew/usecase/login"
	"github.com/deadcheat/cashew/usecase/logout"
	tu "github.com/deadcheat/cashew/usecase/ticket"
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
	defer foundation.Authenticator().Close()

	// create router
	r := mux.NewRouter().PathPrefix(foundation.App().URIPath).Subrouter()

	// create usecase, repository, deliver and mount them
	// repositories
	ticketrep := tr.New(foundation.DB())
	idrep := id.New()
	hostrep := host.New()
	pcrep := proxycallback.New()

	// usecases
	ticketUseCase := tu.New(ticketrep, idrep, hostrep, pcrep)
	loginUseCase := login.New(ticketrep)
	logoutUseCase := logout.New(ticketrep)
	authUseCase := auth.New()
	validateUseCase := validate.New(ticketrep)

	// create deliver and mount
	login := dl.New(r, ticketUseCase, loginUseCase, logoutUseCase, authUseCase)
	login.Mount()
	v := dv.New(r, ticketUseCase, validateUseCase)
	v.Mount()

	// mount to static files
	statics := assets.New(r)
	statics.Mount()

	// start cas server
	bindAddress := fmt.Sprintf("%s:%d", foundation.App().Host, foundation.App().Port)
	log.Println("start cas server on ", bindAddress)
	if foundation.App().UseSSL {
		log.Fatal(http.ListenAndServeTLS(bindAddress, foundation.App().SSLCertFile, foundation.App().SSLCertKey, r))
		return
	}
	log.Fatal(http.ListenAndServe(bindAddress, r))

}

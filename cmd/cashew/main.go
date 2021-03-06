package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/deadcheat/cashew/deliver/assets"
	dli "github.com/deadcheat/cashew/deliver/login"
	dlo "github.com/deadcheat/cashew/deliver/logout"
	dp "github.com/deadcheat/cashew/deliver/proxy"
	dv "github.com/deadcheat/cashew/deliver/validate"
	"github.com/deadcheat/cashew/foundation"
	"github.com/deadcheat/cashew/repository/host"
	"github.com/deadcheat/cashew/repository/id"
	"github.com/deadcheat/cashew/repository/proxycallback"
	tr "github.com/deadcheat/cashew/repository/ticket"
	"github.com/deadcheat/cashew/usecase/auth"
	"github.com/deadcheat/cashew/usecase/terminatelogin"
	"github.com/deadcheat/cashew/usecase/terminatelogout"
	tu "github.com/deadcheat/cashew/usecase/ticket"
	"github.com/deadcheat/cashew/usecase/validate"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	var err error
	var configFile string

	flag.StringVar(&configFile, "c", "config.yml", "specify config file path")
	flag.Parse()

	if err = foundation.PrepareApp(configFile); err != nil {
		log.Fatalln(err)
	}
	defer foundation.Authenticator().Close()
	defer foundation.DB().Close()

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
	loginTerminationUseCase := terminatelogin.New(ticketrep)
	logoutTerminationUseCase := terminatelogout.New(ticketrep)
	authUseCase := auth.New()
	validateUseCase := validate.New(ticketrep)

	// create deliver and mount
	// login
	login := dli.New(r, ticketUseCase, validateUseCase, loginTerminationUseCase, authUseCase)
	login.Mount()
	// logout
	logout := dlo.New(r, ticketUseCase, logoutTerminationUseCase)
	logout.Mount()
	// validate
	v := dv.New(r, ticketUseCase, validateUseCase)
	v.Mount()
	// proxy
	p := dp.New(r, ticketUseCase, validateUseCase)
	p.Mount()

	// mount to static files
	statics := assets.New(r)
	statics.Mount()

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	// start cas server
	bindAddress := fmt.Sprintf("%s:%d", foundation.App().Host, foundation.App().Port)
	log.Println("start cas server on ", bindAddress)
	server := &http.Server{
		Addr:    bindAddress,
		Handler: loggedRouter,
	}
	go func() {
		if foundation.App().UseSSL {
			if err := server.ListenAndServeTLS(foundation.App().SSLCertFile, foundation.App().SSLCertKey); err != nil {
				if err != http.ErrServerClosed {
					log.Fatalln("HTTPServer closed with error:", err)
				}
				log.Println("server has been stopped by signal", err)
			}
		} else {
			if err := server.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					log.Fatalln("HTTPServer closed with error:", err)
				}
				log.Println("server has been stopped by signal", err)
			}
		}
	}()
	// シグナルを待つ
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, os.Interrupt)
	<-stop

	// シグナルを受け取ったらShutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}

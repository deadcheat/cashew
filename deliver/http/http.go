package http

import (
	"log"
	"net/http"

	"github.com/deadcheat/cashew"
	"github.com/gorilla/mux"
)

type Deliver struct {
	r *mux.Router
}

func New(r *mux.Router) cashew.Deliver {
	return &Deliver{r}
}

func (d *Deliver) GetLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("GETlogin")
}

func (d *Deliver) PostLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("PostLogin")
}

func (d *Deliver) Mount() {
	d.r.HandleFunc("/login", d.GetLogin).Methods(http.MethodGet)
	d.r.HandleFunc("/login", d.PostLogin).Methods(http.MethodPost)
}

package http

import (
	"log"
	"net/http"
	"time"

	"github.com/deadcheat/cashew/values/consts"

	"github.com/deadcheat/cashew"
	"github.com/gorilla/mux"
)

type Deliver struct {
	r  *mux.Router
	uc cashew.UseCase
}

// New make new Deliver
func New(r *mux.Router) cashew.Deliver {
	return &Deliver{r: r}
}

// GetLogin Login by "GET" Method
func (d *Deliver) GetLogin(w http.ResponseWriter, r *http.Request) {

	// set to be sure that not use cache
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Expires", time.Now().Add(time.Hour*720).Format(consts.RFC2822))

	params := r.URL.Query()
	service := params[consts.ParamKeyService]
	// renew := params[consts.ParamKeyRenew]
	// gateway := params[consts.ParamKeyGateway]

	tgt, _ := r.Cookie(consts.CookieKeyTGT)
	if err := d.uc.ValidateTicket(tgt.Value); err == nil {
		// if ticket is valid, redirect to service
		http.Redirect(w, r, service[0], http.StatusSeeOther)
		return
	}

	log.Println("GETlogin")
}

func (d *Deliver) PostLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("PostLogin")
}

func (d *Deliver) Mount() {
	d.r.HandleFunc("/login", d.GetLogin).Methods(http.MethodGet)
	d.r.HandleFunc("/login", d.PostLogin).Methods(http.MethodPost)
}

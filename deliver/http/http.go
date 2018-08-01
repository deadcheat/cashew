package http

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/deadcheat/cashew/assets"
	"github.com/deadcheat/cashew/helpers/service"
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
	// define error
	var err error

	// set to be sure that not use cache
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Expires", time.Now().Add(time.Hour*720).Format(consts.RFC2822))

	params := r.URL.Query()
	services := params[consts.ParamKeyService]
	svc := ""
	if len(services) > 0 {
		// enable only first url
		svc, err = service.NormalizeURL(services[0])
		if err != nil {
			log.Println(err)
			http.Error(w, "invalid url for query parameter 'service'", http.StatusBadRequest)
			return
		}
	}
	// check renew and if renew, redirect to login page
	renews := params[consts.ParamKeyRenew]
	for _, v := range renews {
		if v == "1" || strings.ToUpper(v) == "TRUE" {
			LoginPage(w)
			return
		}
	}
	gateways := params[consts.ParamKeyGateway]
	gateway := false
	for _, v := range gateways {
		if v == "1" || strings.ToUpper(v) == "TRUE" {
			gateway = true
		}
	}
	var tgt *http.Cookie
	tgt, err = r.Cookie(consts.CookieKeyTGT)
	if err != nil {
		log.Println(err)
		tgt = &http.Cookie{}
		return
	}

	if gateway {

	}

	if err = d.uc.ValidateTicket(tgt.Value); err == nil && svc != "" {
		// if ticket is valid, redirect to service
		http.Redirect(w, r, svc, http.StatusSeeOther)
		return
	}

	log.Println("GETlogin")
}

func LoginPage(w http.ResponseWriter) {
	t := template.New("cas login")
	f, err := assets.Assets.File("/templates/login/index.html")
	if err != nil {
		http.Error(w, "unabled to find template", http.StatusInternalServerError)
		return
	}
	t, _ = t.Parse(string(f.Data))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusFound)
	t.Execute(w, nil)
}

func (d *Deliver) PostLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("PostLogin")
}

func (d *Deliver) Mount() {
	d.r.HandleFunc("/login", d.GetLogin).Methods(http.MethodGet)
	d.r.HandleFunc("/login", d.PostLogin).Methods(http.MethodPost)
}

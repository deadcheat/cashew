package http

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
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
	var svc *url.URL
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
		if b, _ := strconv.ParseBool(v); b {
			loginPage(w)
			return
		}
	}
	gateways := params[consts.ParamKeyGateway]
	for _, v := range gateways {
		if b, _ := strconv.ParseBool(v); b {
			http.Redirect(w, r, svc.String(), http.StatusSeeOther)
			return
		}
	}
	var tgt *http.Cookie
	tgt, err = r.Cookie(consts.CookieKeyTGT)
	if err != nil {
		log.Println(err)
		tgt = &http.Cookie{}
		return
	}

	// redirect service with service ticket when tgt ticket is valid
	if err = d.uc.ValidateTicket(tgt.Value); err == nil && svc != nil {
		st := d.uc.ServiceTicket()
		q := svc.Query()
		q.Add("ticket", st.ID)
		svc.RawQuery = q.Encode()
		// if ticket is valid, redirect to service
		http.Redirect(w, r, svc.String(), http.StatusSeeOther)
		return
	}

	log.Println("GETlogin")
}

func loginPage(w http.ResponseWriter) {
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

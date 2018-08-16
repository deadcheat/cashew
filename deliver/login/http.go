package http

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/assets"
	"github.com/deadcheat/cashew/helpers/service"
	"github.com/deadcheat/cashew/values/consts"

	"github.com/gorilla/mux"
)

// Deliver struct implements cashew.Deliver
type Deliver struct {
	r  *mux.Router
	uc cashew.LoginUseCase
}

// New make new Deliver
func New(r *mux.Router, uc cashew.LoginUseCase) cashew.Deliver {
	return &Deliver{r: r, uc: uc}
}

// GetLogin handle GET request to /login
func (d *Deliver) GetLogin(w http.ResponseWriter, r *http.Request) {
	// define error
	var err error

	setHeaderNoCache(w)

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
	if svc == nil || stringSliceContainsTrue(renews) {
		loginPage(w)
		return
	}
	gateways := params[consts.ParamKeyGateway]
	if stringSliceContainsTrue(gateways) {
		http.Redirect(w, r, svc.String(), http.StatusSeeOther)
		return
	}

	var tgc *http.Cookie
	tgc, err = r.Cookie(consts.CookieKeyTGT)
	if err != nil {
		log.Println("no ticket granting ticket detected ", err)
		tgc = &http.Cookie{}
	}

	// redirect service with service ticket when tgt ticket is valid
	var tgt *cashew.Ticket
	tgt, err = d.uc.ValidateTicket(cashew.TicketTypeTicketGranting, tgc.Value)
	if err == nil {
		st, err := d.uc.ServiceTicket(svc.String(), tgt)
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to issue service ticket", http.StatusBadRequest)
			return
		}
		q := svc.Query()
		q.Add("ticket", st.ID)
		svc.RawQuery = q.Encode()
		// if ticket is valid, redirect to service
		http.Redirect(w, r, svc.String(), http.StatusSeeOther)
		return
	}

	log.Println("GET login")
}

func setHeaderNoCache(w http.ResponseWriter) {
	// set to be sure that not use cache
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Expires", time.Now().Add(time.Hour*720).Format(consts.RFC2822))
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

// stringSliceContainsTrue return true when src []string contains any true-bool-string
func stringSliceContainsTrue(src []string) bool {
	for _, v := range src {
		b, err := strconv.ParseBool(v)
		if err == nil && b {
			return true
		}
	}
	return false
}

// PostLogin handle post method request to /login
func (d *Deliver) PostLogin(w http.ResponseWriter, r *http.Request) {
	log.Println("Post Login")
}

// Mount route with handler
func (d *Deliver) Mount() {
	d.r.HandleFunc("/login", d.GetLogin).Methods(http.MethodGet)
	d.r.HandleFunc("/login", d.PostLogin).Methods(http.MethodPost)
}

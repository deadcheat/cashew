package http

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/deadcheat/cashew/values/errs"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/assets"
	"github.com/deadcheat/cashew/helpers/service"
	"github.com/deadcheat/cashew/values/consts"

	"github.com/gorilla/mux"
)

// Deliver struct implements cashew.Deliver
type Deliver struct {
	r   *mux.Router
	uc  cashew.LoginUseCase
	auc cashew.AuthenticateUseCase
}

// New make new Deliver
func New(r *mux.Router, uc cashew.LoginUseCase, auc cashew.AuthenticateUseCase) cashew.Deliver {
	return &Deliver{r: r, uc: uc, auc: auc}
}

// GetLogin handle GET request to /login
func (d *Deliver) GetLogin(w http.ResponseWriter, r *http.Request) {
	// define error
	var err error

	setHeaderNoCache(w)

	params := r.URL.Query()
	svc, err := serviceURL(params)
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid url for query parameter 'service'", http.StatusBadRequest)
		return
	}

	// check renew and if renew, redirect to login page
	renews := params[consts.ParamKeyRenew]
	if stringSliceContainsTrue(renews) {
		var lt *cashew.Ticket
		lt, err = d.uc.LoginTicket(r)
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to create login ticket", http.StatusInternalServerError)
			return
		}
		loginPage(w, svc, lt.ID)
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
		tgc = new(http.Cookie)
	}
	tgtID := tgc.Value

	// redirect service with service ticket when tgt ticket is valid
	var tgt *cashew.Ticket
	tgt, err = d.uc.ValidateTicket(cashew.TicketTypeTicketGranting, tgtID)
	switch err {
	case nil:
		if svc == nil {
			log.Println("already logged in and no service detected")
			_, _ = fmt.Fprint(w, "you're already logged in and you didn't set an url to be redirected")
			return
		}
		var st *cashew.Ticket
		st, err = d.uc.ServiceTicket(r, svc, tgt)
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
	case errs.ErrNoTicketID, errs.ErrTicketHasBeenExpired, errs.ErrTicketTypeNotMatched:
		log.Println(err, tgtID)
	default:
		log.Println(err)
		http.Error(w, fmt.Sprintf("error occurred when validating ticket: %s", tgtID), http.StatusInternalServerError)
		return
	}

	var lt *cashew.Ticket
	lt, err = d.uc.LoginTicket(r)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to create login ticket", http.StatusInternalServerError)
		return
	}
	loginPage(w, svc, lt.ID)
}

func setHeaderNoCache(w http.ResponseWriter) {
	// set to be sure that not use cache
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Expires", time.Now().Add(time.Hour*720).Format(consts.RFC2822))
}

func serviceURL(v url.Values) (*url.URL, error) {
	services := v[consts.ParamKeyService]
	serviceURL := firstString(services)
	// enable only first url
	// if serviceURL is empty, return nil, nil
	return service.NormalizeURL(serviceURL)
}

func loginPage(w http.ResponseWriter, svc *url.URL, lt string) {
	t := template.New("cas login")
	f, err := assets.Assets.File("/templates/login/index.html")
	if err != nil {
		http.Error(w, "unabled to find template", http.StatusInternalServerError)
		return
	}
	t, _ = t.Parse(string(f.Data))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusFound)
	t.Execute(w, map[string]interface{}{
		"Service":     svc.String(),
		"LoginTicket": lt,
	})
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
	// define error
	var err error

	setHeaderNoCache(w)

	params := r.URL.Query()
	var svc *url.URL
	svc, err = serviceURL(params)
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid url for query parameter 'service'", http.StatusBadRequest)
		return
	}

	// get required parameters
	u := firstString(params["username"])
	p := firstString(params["password"])
	l := firstString(params["lt"])

	u = strings.Trim(u, " ")

	var lt *cashew.Ticket
	lt, err = d.uc.ValidateTicket(cashew.TicketTypeLogin, l)

	if err != nil {
		log.Println(err)
		http.Error(w, "invalid login ticket", http.StatusBadRequest)
		return
	}
	// FIXME use lt (it should be deleted)
	_ = lt

	// authenticate
	if err = d.auc.Authenticate(u, p); err != nil {
		log.Println(err)
		// FIXME to show error message to loginpage
		http.Error(w, "invalid login ticket", http.StatusUnauthorized)
		return
	}
	var tgt *cashew.Ticket
	tgt, err = d.uc.TicketGrantingTicket(r, u, struct{}{})
	if err != nil {
		log.Println(err)
		// FIXME to show error message to loginpage
		http.Error(w, "failed to issue ticket granting ticket", http.StatusBadRequest)
		return
	}
	// set cookie
	http.SetCookie(w, &http.Cookie{
		Name:  consts.CookieKeyTGT,
		Value: tgt.ID,
		Path:  "/",
	})

	var st *cashew.Ticket
	st, err = d.uc.ServiceTicket(r, svc, tgt)
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
}

// Mount route with handler
func (d *Deliver) Mount() {
	d.r.HandleFunc("/login", d.GetLogin).Methods(http.MethodGet)
	d.r.HandleFunc("/login", d.PostLogin).Methods(http.MethodPost)
}

func firstString(s []string) string {
	if len(s) > 0 {
		return s[0]
	}
	return ""
}

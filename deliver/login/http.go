package http

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/deadcheat/goblet"

	"github.com/deadcheat/cashew/provider/message"
	"github.com/deadcheat/cashew/templates"
	"github.com/deadcheat/cashew/values/errs"

	"github.com/deadcheat/cashew"
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

// get handle GET request to /login
func (d *Deliver) get(w http.ResponseWriter, r *http.Request) {
	// define error
	var err error
	mp := message.New()

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
		err = d.showLoginPage(w, r, svc, false, "", "", mp.Info(), mp.Errors(), http.StatusFound)
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to show login page", http.StatusInternalServerError)
			return
		}
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
	tgt, err = d.uc.FindTicket(tgtID)
	if err == nil {
		err = d.uc.ValidateTicket(cashew.TicketTypeTicketGranting, tgt)
		switch err {
		case nil:
			if svc == nil {
				log.Println("already logged in and no service detected")
				mp.AddInfo("you're already logged in and you didn't set an url to be redirected")
				err = d.showLoginPage(w, r, svc, true, "", "", mp.Info(), mp.Errors(), http.StatusOK)
				if err != nil {
					log.Println(err)
					http.Error(w, "failed to show login page", http.StatusInternalServerError)
					return
				}
				return
			}
			var st *cashew.Ticket
			st, err = d.uc.ServiceTicket(r, svc, tgt)
			if err != nil {
				log.Println(err)
				http.Error(w, "failed to issue service ticket", http.StatusInternalServerError)
				return
			}
			q := svc.Query()
			q.Add("ticket", st.ID)
			svc.RawQuery = q.Encode()
			// if ticket is valid, redirect to service
			http.Redirect(w, r, svc.String(), http.StatusSeeOther)
			return
		case errs.ErrTicketHasBeenExpired, errs.ErrTicketTypeNotMatched:
			log.Println(err, tgtID)
		default:
			log.Println(err)
			http.Error(w, fmt.Sprintf("error occurred when validating ticket: %s", tgtID), http.StatusInternalServerError)
			return
		}
	} else {
		log.Println("tgc not found ", err)
	}
	// display login page
	err = d.showLoginPage(w, r, svc, false, "", "", mp.Info(), mp.Errors(), http.StatusFound)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to show login page", http.StatusInternalServerError)
		return
	}
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

func (d Deliver) showLoginPage(w http.ResponseWriter, r *http.Request, svc *url.URL, loggedIn bool, username, password string, messages []string, errors []string, sc int) (err error) {
	service := ""
	if svc != nil {
		service = svc.String()
	}
	ltID := ""
	if !loggedIn {
		var lt *cashew.Ticket
		lt, err = d.uc.LoginTicket(r)
		if err != nil {
			return
		}
		ltID = lt.ID
	}
	t := template.New("cas login")
	var f *goblet.File
	f, err = templates.Assets.File("/login/index.html")
	if err != nil {
		return
	}
	// FIXME parse process should be done when app start
	t, err = t.Parse(string(f.Data))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(sc)
	return t.Execute(w, map[string]interface{}{
		"Service":     service,
		"LoginTicket": ltID,
		"Messages":    messages,
		"Errors":      errors,
		"LoggedIn":    loggedIn,
		"UserName":    username,
		"Password":    password,
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

// post handle post method request to /login
func (d *Deliver) post(w http.ResponseWriter, r *http.Request) {
	// define error
	var err error
	mp := message.New()

	setHeaderNoCache(w)

	if err = r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "failed to parse posted form data", http.StatusInternalServerError)
		return
	}

	params := r.Form

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
	// delete login ticket
	// if it failed, ignore that instantly
	defer func() {
		if lt != nil {
			return
		}
		if internalErr := d.uc.TerminateLoginTicket(lt); internalErr != nil {
			log.Println("login ticket deletion internal error ", internalErr)
		}
	}()
	lt, err = d.uc.FindTicket(l)
	if err != nil {
		// FIXME redirect to /login with service url
		log.Println(err)
		http.Error(w, "failed to find login ticket", http.StatusBadRequest)
		return
	}

	if err = d.uc.ValidateTicket(cashew.TicketTypeLogin, lt); err != nil {
		// FIXME redirect to /login with service url
		log.Println(err)
		http.Error(w, "failed to find login ticket", http.StatusBadRequest)
		return
	}

	// authenticate
	if err = d.auc.Authenticate(u, p); err != nil {
		// FIXME redirect to /login with service url
		log.Println(err)
		mp.AddErr("your authentication is invalid")
		err = d.showLoginPage(w, r, svc, false, u, p, mp.Info(), mp.Errors(), http.StatusUnauthorized)
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to show page", http.StatusInternalServerError)
			return
		}
		return
	}
	// FIXME for now, we don't get any external attributes
	var tgt *cashew.Ticket
	data, err := json.Marshal(struct{}{})
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to convert extra attributes", http.StatusBadRequest)
		return
	}
	tgt, err = d.uc.TicketGrantingTicket(r, u, data)
	if err != nil {
		// FIXME redirect to /login with service url
		log.Println(err)
		http.Error(w, "failed to issue ticket granting ticket", http.StatusInternalServerError)
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
	switch err {
	case nil:
		break
	case errs.ErrNoServiceDetected:
		mp.AddInfo("you're successfully authenticated but no service param was given and we can't redirect anymore ")
		err = d.showLoginPage(w, r, svc, true, "", "", mp.Info(), mp.Errors(), http.StatusOK)
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to show page", http.StatusInternalServerError)
			return
		}
		return
	default:
		log.Println(err)
		// FIXME to show error message to loginpage
		http.Error(w, "failed to issue ticket granting ticket", http.StatusInternalServerError)
		return
	}

	q := svc.Query()
	q.Add("ticket", st.ID)
	svc.RawQuery = q.Encode()
	// if ticket is valid, redirect to service
	http.Redirect(w, r, svc.String(), http.StatusSeeOther)
}

// GetLogOut handle get method request to /logout
func (d *Deliver) GetLogOut(w http.ResponseWriter, r *http.Request) {

}

// Mount route with handler
func (d *Deliver) Mount() {
	d.r.HandleFunc("/login", d.get).Methods(http.MethodGet)
	d.r.HandleFunc("/login", d.post).Methods(http.MethodPost)
}

func firstString(s []string) string {
	if len(s) > 0 {
		return s[0]
	}
	return ""
}

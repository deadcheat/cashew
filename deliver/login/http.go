package login

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/deadcheat/goblet"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/errors"
	"github.com/deadcheat/cashew/foundation"
	"github.com/deadcheat/cashew/helpers/params"
	hs "github.com/deadcheat/cashew/helpers/strings"
	vh "github.com/deadcheat/cashew/helpers/view"
	"github.com/deadcheat/cashew/provider/message"
	"github.com/deadcheat/cashew/templates"
	"github.com/deadcheat/cashew/timer"
	"github.com/deadcheat/cashew/values/consts"

	"github.com/gorilla/mux"
)

// Deliver struct implements cashew.Deliver
type Deliver struct {
	r    *mux.Router
	tuc  cashew.TicketUseCase
	vuc  cashew.ValidateUseCase
	teuc cashew.TerminateUseCase
	auc  cashew.AuthenticateUseCase
}

// New make new Deliver
func New(r *mux.Router, tuc cashew.TicketUseCase, vuc cashew.ValidateUseCase, teuc cashew.TerminateUseCase, auc cashew.AuthenticateUseCase) cashew.Deliver {
	return &Deliver{r: r, tuc: tuc, vuc: vuc, teuc: teuc, auc: auc}
}

// get handle GET request to /login
func (d *Deliver) get(w http.ResponseWriter, r *http.Request) {
	// define error
	var err error
	mp := message.New()

	setHeaderNoCache(w)

	p := r.URL.Query()
	svc, err := params.ServiceURL(p)
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid url for query parameter 'service'", http.StatusBadRequest)
		return
	}

	// check renew and if renew, redirect to login page
	renews := p[consts.ParamKeyRenew]
	if hs.StringSliceContainsTrue(renews) {
		err = d.showLoginPage(w, r, svc, false, "", "", mp.Info(), mp.Errors(), http.StatusOK)
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to show login page", http.StatusInternalServerError)
			return
		}
		return
	}
	gateways := p[consts.ParamKeyGateway]
	if hs.StringSliceContainsTrue(gateways) {
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
	tgt, err = d.tuc.Find(tgtID)
	if err == nil {
		err = d.vuc.ValidateGranting(tgt)
		switch {
		case err == nil:
			if svc == nil {
				log.Println("already logged in and no service detected")
				mp.AddInfo("you didn't set an url to be redirected.")
				mp.AddInfo(fmt.Sprintf("you're already logged in as %s. If this user is not you, please re-login.", tgt.UserName))
				err = d.showLoginPage(w, r, svc, true, "", "", mp.Info(), mp.Errors(), http.StatusOK)
				if err != nil {
					log.Println(err)
					http.Error(w, "failed to show login page", http.StatusInternalServerError)
					return
				}
				return
			}
			var st *cashew.Ticket
			st, err = d.tuc.NewService(r, svc, tgt, false)
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
		case errors.IsAppError(err):
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
	err = d.showLoginPage(w, r, svc, false, "", "", mp.Info(), mp.Errors(), http.StatusOK)
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

func (d Deliver) showLoginPage(w http.ResponseWriter, r *http.Request, svc *url.URL, loggedIn bool, username, password string, messages []string, errors []string, sc int) (err error) {
	service := ""
	if svc != nil {
		service = svc.String()
	}
	var lt *cashew.Ticket
	lt, err = d.tuc.NewLogin(r)
	if err != nil {
		return
	}
	ltID := lt.ID
	t := template.New("cas login").Funcs(vh.FuncMap)
	var f *goblet.File
	f, err = templates.Assets.File("/files/login/login.html")
	if err != nil {
		return
	}
	// FIXME parse process should be done when app start
	t, _ = t.Parse(string(f.Data))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(sc)
	return t.Execute(w, map[string]interface{}{
		"Service":      service,
		"Organization": foundation.App().Organization,
		"LoginTicket":  ltID,
		"Messages":     messages,
		"Errors":       errors,
		"LoggedIn":     loggedIn,
		"UserName":     username,
		"Password":     password,
	})
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

	p := r.Form

	var svc *url.URL
	svc, err = params.ServiceURL(p)
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid url for query parameter 'service'", http.StatusBadRequest)
		return
	}

	// get required parameters
	u := hs.FirstString(p["username"])
	pa := hs.FirstString(p["password"])
	l := hs.FirstString(p["lt"])

	u = strings.Trim(u, " ")

	var lt *cashew.Ticket
	lt, err = d.tuc.Find(l)
	if err != nil {
		// FIXME redirect to /login with service url
		log.Println(err)
		mp.AddErr("failed to find login ticket or login ticket will be expired")
		dispError := d.showLoginPage(w, r, svc, false, u, pa, mp.Info(), mp.Errors(), http.StatusOK)
		if dispError != nil {
			log.Println(dispError)
			http.Error(w, "failed to show page", http.StatusInternalServerError)
			return
		}
		return
	}
	// delete login ticket
	// if it failed, ignore that instantly
	defer func() {
		if lt == nil {
			return
		}
		if internalErr := d.teuc.Terminate(lt); internalErr != nil {
			log.Println("login ticket deletion internal error ", internalErr)
		}
	}()
	if err = d.vuc.ValidateLogin(lt); err != nil {
		// FIXME redirect to /login with service url
		log.Println(err)
		mp.AddErr("failed to find login ticket or login ticket will be expired")
		dispError := d.showLoginPage(w, r, svc, false, u, pa, mp.Info(), mp.Errors(), http.StatusOK)
		if dispError != nil {
			log.Println(dispError)
			http.Error(w, "failed to show page", http.StatusInternalServerError)
			return
		}
		return
	}

	// authenticate
	var attr map[string]interface{}
	attr, err = d.auc.Authenticate(u, pa)
	if err != nil {
		// FIXME redirect to /login with service url
		log.Println(err)
		mp.AddErr("your authentication is invalid: " + err.Error())
		err = d.showLoginPage(w, r, svc, false, u, pa, mp.Info(), mp.Errors(), http.StatusUnauthorized)
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to show page", http.StatusInternalServerError)
			return
		}
		return
	}
	var data interface{}
	if attr != nil {
		var dataBytes []byte
		dataBytes, err = json.Marshal(attr)
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to convert extra attributes", http.StatusBadRequest)
			return
		}
		if len(dataBytes) > 0 {
			data = dataBytes
		}
	}

	var tgt *cashew.Ticket
	tgt, err = d.tuc.NewGranting(r, u, data)
	if err != nil {
		// FIXME redirect to /login with service url
		log.Println(err)
		http.Error(w, "failed to issue ticket granting ticket", http.StatusInternalServerError)
		return
	}
	// set cookie
	tgtCookie := &http.Cookie{
		Name:  consts.CookieKeyTGT,
		Value: tgt.ID,
		Path:  filepath.Join("/", foundation.App().URIPath),
	}

	gde := foundation.App().GrantingDefaultExpire
	if gde > 0 {
		tgtCookie.Expires = timer.Local.Now().Add(time.Duration(gde) * time.Second)
	}
	http.SetCookie(w, tgtCookie)

	var st *cashew.Ticket
	st, err = d.tuc.NewService(r, svc, tgt, true)
	switch err {
	case nil:
		break
	case errors.ErrNoServiceDetected:
		mp.AddInfo("you're successfully authenticated but no service param was given and we can't redirect anymore ")
		err = d.showLoginPage(w, r, svc, true, u, "", mp.Info(), mp.Errors(), http.StatusOK)
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to show page", http.StatusInternalServerError)
			return
		}
		return
	default:
		log.Println(err)
		mp.AddErr("failed to issue ticket granting ticket " + err.Error())
		err = d.showLoginPage(w, r, svc, false, u, pa, mp.Info(), mp.Errors(), http.StatusOK)
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to show page", http.StatusInternalServerError)
			return
		}
		return
	}

	q := svc.Query()
	q.Add("ticket", st.ID)
	svc.RawQuery = q.Encode()
	// if ticket is valid, redirect to service
	http.Redirect(w, r, svc.String(), http.StatusSeeOther)
}

// index handle get method request to /
func (d *Deliver) index(w http.ResponseWriter, r *http.Request) {
	u := r.URL
	u.Path = u.Path + "login"
	cs := r.Cookies()
	for i := range cs {
		http.SetCookie(w, cs[i])
	}
	http.Redirect(w, r, u.String(), http.StatusSeeOther)
}

// Mount route with handler
func (d *Deliver) Mount() {
	d.r.HandleFunc("/", d.index).Methods(http.MethodGet)
	d.r.HandleFunc("/login", d.get).Methods(http.MethodGet)
	d.r.HandleFunc("/login", d.post).Methods(http.MethodPost)
}

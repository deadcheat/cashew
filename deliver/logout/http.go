package logout

import (
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/deadcheat/goblet"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/foundation"
	"github.com/deadcheat/cashew/helpers/params"
	hs "github.com/deadcheat/cashew/helpers/strings"
	vh "github.com/deadcheat/cashew/helpers/view"
	"github.com/deadcheat/cashew/provider/message"
	"github.com/deadcheat/cashew/templates"
	"github.com/deadcheat/cashew/values/consts"

	"github.com/gorilla/mux"
)

// Deliver struct implements cashew.Deliver
type Deliver struct {
	r    *mux.Router
	tuc  cashew.TicketUseCase
	louc cashew.TerminateUseCase
}

// New make new Deliver
func New(r *mux.Router, tuc cashew.TicketUseCase, louc cashew.TerminateUseCase) cashew.Deliver {
	return &Deliver{r: r, tuc: tuc, louc: louc}
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
	ltID := ""
	if !loggedIn {
		var lt *cashew.Ticket
		lt, err = d.tuc.NewLogin(r)
		if err != nil {
			return
		}
		ltID = lt.ID
	}
	t := template.New("cas logout").Funcs(vh.FuncMap)
	var f *goblet.File
	f, err = templates.Assets.File("/files/login/index.html")
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

// logout handle get method request to /logout
func (d *Deliver) logout(w http.ResponseWriter, r *http.Request) {
	var err error
	mp := message.New()

	p := r.URL.Query()
	// FIXME i might check err
	svc, _ := params.ServiceURL(p)
	next, _ := params.ContinueURL(p)
	gateways := p[consts.ParamKeyGateway]

	var tgc *http.Cookie
	tgc, err = r.Cookie(consts.CookieKeyTGT)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to find ticket granting ticket in cookie", http.StatusBadRequest)
		return
	}
	tgtID := tgc.Value
	var tgt *cashew.Ticket
	tgt, err = d.tuc.Find(tgtID)
	if err == nil {
		if err = d.louc.Terminate(tgt); err != nil {
			log.Println(err)
			http.Error(w, "failed to delete ticket granting ticket", http.StatusInternalServerError)
			return
		}
	}

	// delete cookie
	http.SetCookie(w, &http.Cookie{
		Name:    consts.CookieKeyTGT,
		Value:   "",
		Path:    filepath.Join("/", foundation.App().URIPath),
		Expires: time.Unix(0, 0),
	})
	if hs.StringSliceContainsTrue(gateways) && svc != nil && svc.String() != "" {
		http.Redirect(w, r, svc.String(), http.StatusSeeOther)
		return
	}
	if next != nil {
		tmp := template.New("cas logout").Funcs(vh.FuncMap)
		var f *goblet.File
		f, err = templates.Assets.File("/files/login/logout.html")
		if err != nil {
			return
		}
		// FIXME parse process should be done when app start
		tmp, err = tmp.Parse(string(f.Data))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if err = tmp.Execute(w, map[string]interface{}{
			"Next": next.String(),
		}); err != nil {
			log.Println(err)
			http.Error(w, "failed to show logout page", http.StatusInternalServerError)
			return
		}
		return
	}
	// finally render login display
	mp.AddInfo("You're successfully logged out.")
	err = d.showLoginPage(w, r, svc, false, "", "", mp.Info(), mp.Errors(), http.StatusOK)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to show page", http.StatusInternalServerError)
		return
	}
}

// Mount route with handler
func (d *Deliver) Mount() {
	d.r.HandleFunc("/logout", d.logout).Methods(http.MethodGet)
}

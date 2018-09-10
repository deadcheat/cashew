package validate

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/helpers/errors"
	"github.com/deadcheat/cashew/helpers/params"
	"github.com/deadcheat/cashew/helpers/strings"
	"github.com/deadcheat/cashew/templates"
	"github.com/deadcheat/cashew/values/consts"
	"github.com/deadcheat/cashew/values/errs"
	"github.com/deadcheat/goblet"
	"github.com/gorilla/mux"
)

// Deliver struct implements cashew.Deliver
type Deliver struct {
	r   *mux.Router
	tuc cashew.TicketUseCase
	vuc cashew.ValidateUseCase
}

// New make new Deliver
func New(r *mux.Router, tuc cashew.TicketUseCase, vuc cashew.ValidateUseCase) cashew.Deliver {
	return &Deliver{r: r, tuc: tuc, vuc: vuc}
}

func (d *Deliver) validate(w http.ResponseWriter, r *http.Request) {
	isValidated := "no"
	foundUser := ""

	p := r.URL.Query()
	svc, err := params.ServiceURL(p)
	if err != nil || svc == nil {
		log.Printf("error: [%+v] service: [%+v]", err, svc)
		http.Error(w, "invalid url for query parameter 'service'", http.StatusBadRequest)
		return
	}
	ticket := strings.FirstString(p["ticket"])

	renews := p[consts.ParamKeyRenew]

	t, err := d.vuc.Validate(ticket, svc, strings.StringSliceContainsTrue(renews), false)
	if err == nil {
		isValidated = "yes"
		foundUser = t.UserName
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err = fmt.Fprintf(w, "%s\n%s\n", isValidated, foundUser); err != nil {
		http.Error(w, "failed to show response", http.StatusInternalServerError)
	}
}

func (d *Deliver) serviceValidate(w http.ResponseWriter, r *http.Request) {

	// serviceValidate will be rendered as utf-8 xml
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")

	var v view
	p := r.URL.Query()
	// service and ticket are required parameter
	ticket := strings.FirstString(p["ticket"])
	svc, err := params.ServiceURL(p)
	if err != nil || svc == nil || ticket == "" {
		v.e = errors.NewInvalidRequest(errs.ErrRequiredParameterMissed)
		err = d.showServiceValidateXML(w, r, v)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			log.Println(err)
			http.Error(w, "failed to show xml", http.StatusInternalServerError)
			return
		}
		return
	}

	// pgtUrl is optional
	pgtURL, err := params.PgtURL(p)
	if err != nil {
		log.Println("invalid url for query parameter 'pgtUrl'", err)
	}
	renews := p[consts.ParamKeyRenew]
	var st *cashew.Ticket
	st, err = d.vuc.Validate(ticket, svc, strings.StringSliceContainsTrue(renews), false)
	if err != nil {
		// diplay failed xml and finish process
		switch err {
		case errs.ErrTicketNotFound,
			errs.ErrTicketIsProxyTicket,
			errs.ErrServiceTicketIsNoPrimary,
			errs.ErrTicketTypeNotMatched:
			v.e = errors.NewInvalidTicket(ticket, err)
		case errs.ErrServiceURLNotMatched:
			v.e = errors.NewInvalidService(err)
		default:
			v.e = errors.NewInternalError(err)
		}
		err = d.showServiceValidateXML(w, r, v)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			log.Println(err)
			http.Error(w, "failed to show xml", http.StatusInternalServerError)
			return
		}
		return
	}
	var pgt *cashew.Ticket
	pgt, err = d.tuc.ProxyGrantingTicket(r, pgtURL, st)
	if err != nil {
		switch err {
		case errs.ErrProxyCallBackURLMissing:
			// do nothing
		case errs.ErrProxyGrantingURLUnexpectedStatus:
			v.e = errors.NewInvalidProxyCallback(err)
		default:
			v.e = errors.NewInternalError(err)
		}
	}
	v.pgtiou = pgt.IOU
	v.proxied = true
	v.success = true
	v.username = st.UserName
	err = d.showServiceValidateXML(w, r, v)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		log.Println(err)
		http.Error(w, "failed to show xml", http.StatusInternalServerError)
	}
}

type view struct {
	success  bool
	username string
	proxied  bool
	pgt      string
	pgtiou   string
	proxies  []*url.URL
	e        errors.Wrapper
}

func (d *Deliver) showServiceValidateXML(w http.ResponseWriter, r *http.Request, v view) (err error) {
	t := template.New("cas service validate")
	var f *goblet.File
	f, err = templates.Assets.File("/validate/servicevalidate.xml")
	if err != nil {
		return
	}
	t, err = t.Parse(string(f.Data))
	if err != nil {
		return
	}
	var errCode, errBody string
	if v.e != nil {
		errCode = v.e.Code()
		errBody = v.e.Message()
		if v.e.Is(errors.ErrorCodeInternalError) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
	return t.Execute(w, map[string]interface{}{
		"AuthenticationSuccess":  v.success,
		"UserName":               v.username,
		"ProxyGrantingTicketIOU": v.pgtiou,
		"ProxyGranting":          v.proxied,
		"HasProxies":             (len(v.proxies) > 0),
		"Proxies":                v.proxies,
		"ErrorCode":              errCode,
		"ErrorBody":              errBody,
	})
}

// Mount route with handler
func (d *Deliver) Mount() {
	d.r.HandleFunc("/validate", d.validate).Methods(http.MethodGet)
	d.r.HandleFunc("/serviceValidate", d.serviceValidate).Methods(http.MethodGet)
}

package validate

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/errors"
	"github.com/deadcheat/cashew/helpers/params"
	"github.com/deadcheat/cashew/helpers/strings"
	"github.com/deadcheat/cashew/templates"
	"github.com/deadcheat/cashew/values/consts"
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
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	p := r.URL.Query()
	svc, err := params.ServiceURL(p)
	if err != nil || svc == nil {
		log.Printf("error: [%+v] service: [%+v]", err, svc)
		http.Error(w, "invalid url for query parameter 'service'", http.StatusBadRequest)
		return
	}
	ticket := strings.FirstString(p["ticket"])

	renews := p[consts.ParamKeyRenew]
	t, err := d.tuc.Find(ticket)
	if err != nil {
		if _, err = fmt.Fprintf(w, "%s\n%s\n", isValidated, foundUser); err != nil {
			http.Error(w, "failed to show response", http.StatusInternalServerError)
		}
		return
	}
	err = d.vuc.ValidateService(t, svc, strings.StringSliceContainsTrue(renews))
	if err == nil {
		isValidated = "yes"
		foundUser = t.UserName
	}
	if _, err = fmt.Fprintf(w, "%s\n%s\n", isValidated, foundUser); err != nil {
		http.Error(w, "failed to show response", http.StatusInternalServerError)
	}
}

func (d *Deliver) serviceValidate(w http.ResponseWriter, r *http.Request) {
	d.fragmentValidate(w, r, false)
}

func (d *Deliver) proxyValidate(w http.ResponseWriter, r *http.Request) {
	d.fragmentValidate(w, r, true)
}

func (d *Deliver) fragmentValidate(w http.ResponseWriter, r *http.Request, checkAsProxyValidation bool) {
	var v view
	p := r.URL.Query()
	// service and ticket are required parameter
	ticket := strings.FirstString(p["ticket"])
	svc, err := params.ServiceURL(p)
	if err != nil || svc == nil || ticket == "" {
		v.e = errors.ErrRequiredParameterMissed
		d.showServiceValidateXML(w, r, v)
		return
	}

	// pgtUrl is optional
	pgtURL, err := params.PgtURL(p)
	if err != nil {
		log.Println("invalid url for query parameter 'pgtUrl'", err)
	}
	renews := p[consts.ParamKeyRenew]
	var st *cashew.Ticket
	st, err = d.tuc.Find(ticket)
	if err != nil {
		// diplay failed xml and finish process
		var ok bool
		v.e, ok = errors.AsErrorView(err)
		if !ok {
			v.e = errors.NewInternalError(err)
		}
		d.showServiceValidateXML(w, r, v)
		return
	}
	if checkAsProxyValidation {
		err = d.vuc.ValidateProxy(st, svc)
	} else {
		err = d.vuc.ValidateService(st, svc, strings.StringSliceContainsTrue(renews))
	}
	if err != nil {
		// diplay failed xml and finish process
		var ok bool
		v.e, ok = errors.AsErrorView(err)
		if !ok {
			v.e = errors.NewInternalError(err)
		}
		d.showServiceValidateXML(w, r, v)
		return
	}
	var pgt *cashew.Ticket
	pgt, err = d.tuc.NewProxyGranting(r, pgtURL, st)
	if err != nil {
		var ok bool
		v.e, ok = errors.AsErrorView(err)
		if !ok {
			v.e = errors.NewInternalError(err)
		}
	} else {
		v.Success = true
		v.IOU = pgt.IOU
	}
	v.Name = st.UserName
	d.showServiceValidateXML(w, r, v)
}

type view struct {
	Success   bool
	Name      string
	IOU       string
	Proxies   []*url.URL
	e         errors.ErrorView
	ErrorCode string
	ErrorBody string
}

func (d *Deliver) showServiceValidateXML(w http.ResponseWriter, r *http.Request, v view) {
	var err error
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
	// serviceValidate will be rendered as utf-8 xml
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	if v.e != nil {
		v.ErrorCode = v.e.Code()
		v.ErrorBody = v.e.Message()
		if v.e.Is(errors.ErrorCodeInternalError) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
	if err = t.Execute(w, v); err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		log.Println(err)
		http.Error(w, "failed to show xml", http.StatusInternalServerError)
	}
	return
}

// Mount route with handler
func (d *Deliver) Mount() {
	d.r.HandleFunc("/validate", d.validate).Methods(http.MethodGet)
	d.r.HandleFunc("/serviceValidate", d.serviceValidate).Methods(http.MethodGet)
	d.r.HandleFunc("/proxyValidate", d.proxyValidate).Methods(http.MethodGet)
}

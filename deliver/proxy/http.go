package proxy

import (
	"html/template"
	"log"
	"net/http"

	"github.com/deadcheat/cashew/templates"
	"github.com/deadcheat/goblet"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/errors"
	"github.com/deadcheat/cashew/helpers/params"
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

func (d *Deliver) proxy(w http.ResponseWriter, r *http.Request) {
	var err error
	p := r.URL.Query()
	pgt := params.Pgt(p)
	targetService := params.TargetService(p)

	v := new(view)
	if pgt == "" || targetService == "" {
		v.e = errors.ErrRequiredParameterMissed
		d.showServiceValidateXML(w, r, v)
		return
	}
	t, err := d.tuc.Find(pgt)
	if err != nil {
		v.e = errors.ErrTicketNotFound
		d.showServiceValidateXML(w, r, v)
		return
	}
	err = d.vuc.ValidateProxyGranting(t)
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

	var proxy *cashew.Ticket
	proxy, err = d.tuc.NewProxy(r, targetService, t)
	if err != nil {
		v.e = errors.NewInternalError(err)
	} else {
		v.Ticket = proxy.ID
		v.Success = true
	}
	d.showServiceValidateXML(w, r, v)
}

type view struct {
	Success   bool
	Ticket    string
	e         errors.ErrorView
	ErrorCode string
	ErrorBody string
}

func (d *Deliver) showServiceValidateXML(w http.ResponseWriter, r *http.Request, v *view) {
	var err error
	t := template.New("cas service validate")
	var f *goblet.File
	f, err = templates.Assets.File("/proxy/proxy.xml")
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
	d.r.HandleFunc("/proxy", d.proxy).Methods(http.MethodGet)
}

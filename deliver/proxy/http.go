package proxy

import (
	"html/template"
	"net/http"

	"github.com/deadcheat/cashew/templates"
	"github.com/deadcheat/cashew/values/errs"
	"github.com/deadcheat/goblet"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/helpers/errors"
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
	p := r.URL.Query()
	pgt := params.Pgt(p)
	targetService := params.TargetService(p)
	v := new(view)
	if pgt == "" || targetService == "" {
		v.e = errors.NewInvalidRequest(errs.ErrRequiredParameterMissed)
	}
}

type view struct {
	Success   bool
	Ticket    string
	e         errors.Wrapper
	ErrorCode string
	ErrorBody string
}

func (d *Deliver) showServiceValidateXML(w http.ResponseWriter, r *http.Request, v view) (err error) {
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
	return t.Execute(w, v)
}

// Mount route with handler
func (d *Deliver) Mount() {
	d.r.HandleFunc("/proxy", d.proxy).Methods(http.MethodGet)
}

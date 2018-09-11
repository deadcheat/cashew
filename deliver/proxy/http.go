package proxy

import (
	"html/template"
	"log"
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

	if pgt == "" || targetService == "" {
		log.Println(errs.ErrRequiredParameterMissed)
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
	return t.Execute(w, v)
}

// Mount route with handler
func (d *Deliver) Mount() {
	d.r.HandleFunc("/proxy", d.proxy).Methods(http.MethodGet)
}
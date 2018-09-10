package proxy

import (
	"log"
	"net/http"

	"github.com/deadcheat/cashew/values/errs"

	"github.com/deadcheat/cashew"
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

// Mount route with handler
func (d *Deliver) Mount() {
	d.r.HandleFunc("/proxy", d.proxy).Methods(http.MethodGet)
}

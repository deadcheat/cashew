package validate

import (
	"net/http"

	"github.com/deadcheat/cashew"
	"github.com/gorilla/mux"
)

// Deliver struct implements cashew.Deliver
type Deliver struct {
	r  *mux.Router
	uc cashew.ValidateUseCase
}

// New make new Deliver
func New(r *mux.Router, uc cashew.ValidateUseCase) cashew.Deliver {
	return &Deliver{r: r, uc: uc}
}

func (d *Deliver) validate(w http.ResponseWriter, r *http.Request) {

}

// Mount route with handler
func (d *Deliver) Mount() {
	d.r.HandleFunc("/validate", d.validate).Methods(http.MethodGet)
}

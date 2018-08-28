package validate

import (
	"fmt"
	"net/http"
	"net/url"

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
	isValidated := "no"
	foundUser := ""

	ticket := ""
	service := new(url.URL)

	t, err := d.uc.Validate(ticket, service)
	if err == nil {
		isValidated = "yes"
		foundUser = t.UserName
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "%s\n\n%s", isValidated, foundUser)
}

// Mount route with handler
func (d *Deliver) Mount() {
	d.r.HandleFunc("/validate", d.validate).Methods(http.MethodGet)
}

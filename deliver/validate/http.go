package validate

import (
	"fmt"
	"log"
	"net/http"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/helpers/params"
	"github.com/deadcheat/cashew/helpers/strings"
	"github.com/deadcheat/cashew/values/consts"
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

	p := r.URL.Query()
	svc, err := params.ServiceURL(p)
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid url for query parameter 'service'", http.StatusBadRequest)
		return
	}
	ticket := strings.FirstString(p["ticket"])

	renews := p[consts.ParamKeyRenew]

	t, err := d.uc.Validate(ticket, svc, strings.StringSliceContainsTrue(renews))
	if err == nil {
		isValidated = "yes"
		foundUser = t.UserName
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err = fmt.Fprintf(w, "%s\n%s\n", isValidated, foundUser); err != nil {
		http.Error(w, "failed to show response", http.StatusInternalServerError)
	}
}

// Mount route with handler
func (d *Deliver) Mount() {
	d.r.HandleFunc("/validate", d.validate).Methods(http.MethodGet)
}

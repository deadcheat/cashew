package host

import (
	"net/http"

	"github.com/deadcheat/cashew"

	"github.com/tomasen/realip"
)

// Repository implements cashew.ClientHostNameRepository
type Repository struct{}

func New() cashew.ClientHostNameRepository

func (r *Repository) Ensure(r *http.Request) string {
	return realip.FromRequest(r)
}

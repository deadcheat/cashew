package host

import (
	"net/http"

	"github.com/deadcheat/cashew"

	"github.com/tomasen/realip"
)

// Repository implements cashew.ClientHostNameRepository
type Repository struct{}

// New return Repository implemented structure
func New() cashew.ClientHostNameRepository {
	return &Repository{}
}

// Ensure return real ip/host from request
func (re *Repository) Ensure(r *http.Request) string {
	return realip.FromRequest(r)
}

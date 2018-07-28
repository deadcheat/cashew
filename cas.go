package cashew

import "net/http"

// Deliver delivery interface
type Deliver interface {
	Mount()
	GetLogin(w http.ResponseWriter, r *http.Request)
	PostLogin(w http.ResponseWriter, r *http.Request)
}

// UseCase define behaviors for Cas Server
type UseCase interface {
	Login() error
}

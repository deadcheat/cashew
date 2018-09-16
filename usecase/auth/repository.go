package auth

import (
	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/auth/credential"
	"github.com/deadcheat/cashew/foundation"
)

// Usecase authentication Usecase
type Usecase struct {
}

// New return new interface implemented struct
func New() cashew.AuthenticateUseCase {
	return &Usecase{}
}

// Authenticate return error if id/pass is not valid
func (u Usecase) Authenticate(id, pass string) (map[string]interface{}, error) {
	e := credential.Entity{
		Key:    id,
		Secret: pass,
	}
	attr, err := foundation.Authenticator().Authenticate(&e)
	return attr, err
}

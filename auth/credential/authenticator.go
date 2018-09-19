package credential

import (
	"errors"
	"io"
)

var (
	// ErrAuthenticateFailed 認証エラー
	ErrAuthenticateFailed = errors.New("authentication failed")

	// ErrMultipleUserFound defined error when multiple users matched identification
	ErrMultipleUserFound = errors.New("there are many users to match user/password")
)

// Authenticator authentic interface
type Authenticator interface {
	Authenticate(c *Entity) (Attributes, error)
	io.Closer
}

// AuthenticationBuilder will prepare authenticator
type AuthenticationBuilder interface {
	Build() (Authenticator, error)
}

// Entity holds user/password
type Entity struct {
	Key    string
	Secret string
}

// Attributes have extra-attributes
type Attributes map[string]interface{}

// Set is a setter for map
func (a Attributes) Set(k string, v interface{}) {
	a[k] = v
}

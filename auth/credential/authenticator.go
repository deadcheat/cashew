package credential

import (
	"errors"
	"io"
)

var (
	// ErrAuthenticateFailed 認証エラー
	ErrAuthenticateFailed = errors.New("authentication failed")
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

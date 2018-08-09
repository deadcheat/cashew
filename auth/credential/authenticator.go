package credential

import "errors"

var (
	// ErrAuthenticateFailed 認証エラー
	ErrAuthenticateFailed = errors.New("authentication failed")
)

// Authenticator authentic interface
type Authenticator interface {
	Authenticate(c *Entity) error
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

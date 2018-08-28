package dbauth

import (
	"testing"
)

func TestValidate(t *testing.T) {
	password := "password"

	u := user{
		name: "test",
		pass: "6917f0150a9e99a8bdbb506fe310c2d96f3190320732122b7afbfdc3898be12f",
		salt: "",
	}

	a := &Authenticator{}

	err := a.validate(password, &u)
	if err != nil {
		t.Error(err)
	}
}

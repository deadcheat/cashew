package foundation

import (
	"fmt"
	"strings"

	"github.com/deadcheat/cashew/auth/credential/dbauth"

	"github.com/deadcheat/cashew/auth/credential"
	"github.com/deadcheat/cashew/setting"
)

// PrepareAuthenticator prepare authenticator var
func PrepareAuthenticator() (err error) {
	if app == nil {
		return ErrSettingHasNotBeenLoaded
	}
	authenticator, err = assignAuthenticator(app.Authenticator)
	return nil
}

// assignAuthenticator assign new local authenticator
func assignAuthenticator(a *setting.Authenticator) (credential.Authenticator, error) {
	driver := strings.ToLower(a.Driver)
	switch driver {
	case "database":
		dbInfo := a.AuthDatabase
		authDB, err := openDB(dbInfo.Database)
		if err != nil {
			return nil, err
		}
		b := dbauth.NewAuthenticationBuilder(authDB, dbInfo.Table, dbInfo.UserNameKey, dbInfo.PasswordKey)
		return b.Build()
	}
	return nil, fmt.Errorf("unknown driver: %s", driver)
}

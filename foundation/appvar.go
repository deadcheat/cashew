package foundation

import (
	"database/sql"
	"errors"

	"github.com/deadcheat/cashew/auth/credential"

	"github.com/deadcheat/cashew/setting"
)

var (
	// loader setting loader accessable globally
	loader *settingLoader

	app *setting.App

	// db global database connection
	db *sql.DB

	// authenticator hold auth implements
	authenticator credential.Authenticator
)

var (
	// ErrSettingHasNotBeenLoaded setting has not been loaded
	ErrSettingHasNotBeenLoaded = errors.New("setting has not been loaded")
)

// SettingLoader access local setting loader from global
func SettingLoader() setting.Loader {
	return loader
}

// DB holds db connection globally
func DB() *sql.DB {
	return db
}

// App return setting held first
func App() *setting.App {
	return app
}

// Authenticator return predeclared authenticator
func Authenticator() credential.Authenticator {
	return authenticator
}

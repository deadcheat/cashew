package foundation

import (
	"database/sql"

	"github.com/deadcheat/cashew/setting"
)

var (
	// loader setting loader accessable globally
	loader *settingLoader

	app *setting.App

	// db global database connection
	db *sql.DB
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

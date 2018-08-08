package foundation

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/deadcheat/cashew/database"
	"github.com/deadcheat/cashew/database/mysql"
	"github.com/deadcheat/cashew/setting"
)

// StartDatabase connect database and hold connection globally
func StartDatabase() (err error) {
	if app == nil {
		return errors.New("setting has not been loaded")
	}
	db, err = openDB(app.Database)
	return
}

// openDB open database
func openDB(d *setting.Database) (*sql.DB, error) {
	driver := strings.ToLower(d.Driver)
	var c database.Connector
	switch driver {
	case "mysql":
		c = mysql.New(
			d.Name,
			d.User,
			d.Pass,
			d.Host,
			d.Port,
			map[string]string{
				"parseTime": "true",
			},
		)
	default:
		return nil, fmt.Errorf("unknown driver: %s", driver)
	}
	return c.Open()
}

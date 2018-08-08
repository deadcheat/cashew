package foundation

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/deadcheat/cashew/database/mysql"
	"github.com/deadcheat/cashew/setting"
)

func OpenDB(d setting.Database) (*sql.DB, error) {
	driver := strings.ToLower(d.Driver)
	var c Connector
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

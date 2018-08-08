package mysql

import (
	"database/sql"
	"fmt"
	"net/url"

	// do blank import on this line
	_ "github.com/go-sql-driver/mysql"
)

// Connector interface implements
type Connector struct {
	Name   string
	User   string
	Pass   string
	Host   string
	Port   int
	Params map[string]string
}

// Open connection
func (c *Connector) Open() (*sql.DB, error) {
	params := url.Values{}
	params.Add("parseTime", "true")
	for k, v := range c.Params {
		params.Add(k, v)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", c.User, c.Pass, c.Host, c.Port, c.Name, params.Encode())
	return sql.Open("mysql", dsn)
}

package mysql

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/deadcheat/cashew/database"

	// do blank import on this line
	_ "github.com/go-sql-driver/mysql"
)

// New return new connector
func New(name, user, pass, socket, host string, port int, params map[string]string) database.Connector {
	return &Connector{
		Name:   name,
		User:   user,
		Pass:   pass,
		Host:   host,
		Port:   port,
		Socket: socket,
		Params: params,
	}
}

// Connector interface implements
type Connector struct {
	Name   string
	User   string
	Pass   string
	Host   string
	Port   int
	Socket string
	Params map[string]string
}

// Open connection
func (c *Connector) Open() (*sql.DB, error) {
	params := url.Values{}
	for k, v := range c.Params {
		params.Add(k, v)
	}
	destination := c.Socket
	if destination == "" {
		destination = fmt.Sprintf("%s:%d", c.Host, c.Port)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", c.User, c.Pass, destination, c.Name, params.Encode())
	return sql.Open("mysql", dsn)
}

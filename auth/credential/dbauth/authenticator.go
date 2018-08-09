package dbauth

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/deadcheat/cashew/auth/credential"
)

var stmt *sql.Stmt

// Authenticator implement auth.Authenticator
type Authenticator struct{}

// AuthenticationBuilder env for authenticate
type AuthenticationBuilder struct {
	db             *sql.DB
	table          string
	userNameColumn string
	passWordColumn string
}

// NewAuthenticationBuilder create new builder
func NewAuthenticationBuilder(db *sql.DB, table, userNameColumn, passWordColumn string) credential.AuthenticationBuilder {
	return &AuthenticationBuilder{
		db:             db,
		table:          table,
		userNameColumn: userNameColumn,
		passWordColumn: passWordColumn,
	}
}

// Build prepare authenticate environment and prepare statement
func (a *AuthenticationBuilder) Build() (credential.Authenticator, error) {
	var err error
	stmt, err = a.db.Prepare(fmt.Sprintf(queryBase, a.table, a.userNameColumn, a.passWordColumn))
	if err != nil {
		return nil, err
	}
	return new(Authenticator), nil
}

var queryBase = `SELECT 1
FROM 
	%s target
WHERE
	target.%s = ?
AND
	target.%s = ?
`

var (
	// ErrMultipleUserFound defined error when multiple users matched identification
	ErrMultipleUserFound = errors.New("there are many users to match user/password")
)

// Authenticate implement method for auth.Authenticator
func (a *Authenticator) Authenticate(c *credential.Entity) (err error) {
	var r *sql.Rows
	r, err = stmt.Query(c.Key, c.Secret)
	if err != nil {
		return
	}
	defer r.Close()
	count := 0
	for r.Next() {
		c := new(interface{})
		if err = r.Scan(c); err != nil {
			return err
		}
		count++
	}
	if err = r.Err(); err != nil {
		return err
	}
	if count > 1 {
		return ErrMultipleUserFound
	}
	return nil
}

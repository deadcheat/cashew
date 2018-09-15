package dbauth

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"strings"

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
	saltColumn     string
	attributes     map[string]string
}

// NewAuthenticationBuilder create new builder
func NewAuthenticationBuilder(db *sql.DB, table, userNameColumn, passWordColumn, saltColumn string, attributes map[string]string) credential.AuthenticationBuilder {

	return &AuthenticationBuilder{
		db:             db,
		table:          table,
		userNameColumn: userNameColumn,
		passWordColumn: passWordColumn,
		saltColumn:     saltColumn,
		attributes:     attributes,
	}
}

// Build prepare authenticate environment and prepare statement
func (a *AuthenticationBuilder) Build() (credential.Authenticator, error) {
	var err error
	stmt, err = a.db.Prepare(a.createSelectStatement())
	if err != nil {
		return nil, err
	}
	return new(Authenticator), nil
}

func (a *AuthenticationBuilder) createSelectStatement() string {
	saltPhrase := "''"
	if a.saltColumn != "" {
		saltPhrase = fmt.Sprintf("target.%s", a.saltColumn)
	}
	attributeSlice := make([]string, len(a.attributes))
	i := 0
	for k, v := range a.attributes {
		attributeSlice[i] = fmt.Sprintf(attributeQueryFormat, k, v)
		i++
	}
	attributePhrase := strings.Join(attributeSlice, "\n")
	queryFormat := fmt.Sprintf("%s%s%s", selectBaseFormat, attributePhrase, fromFormat)
	return fmt.Sprintf(queryFormat, a.userNameColumn, a.passWordColumn, saltPhrase, a.table, a.userNameColumn)
}

var (
	selectBaseFormat = `SELECT
  target.%s as user,
  target.%s as password,
  %s as salt
`

	fromFormat = `
FROM 
  %s target
WHERE
  target.%s = ?
`
	attributeQueryFormat = "  , %s as %s"
)

type user struct {
	name string
	pass string
	salt string
}

// Close implemented Closer to close inner stmt
func (a *Authenticator) Close() error {
	return stmt.Close()
}

// Authenticate implement method for auth.Authenticator
func (a *Authenticator) Authenticate(c *credential.Entity) (attr credential.Attributes, err error) {
	var r *sql.Rows
	r, err = stmt.Query(c.Key)
	if err != nil {
		return
	}
	defer r.Close()
	count := 0
	var u *user
	for r.Next() {
		u = new(user)
		if err = r.Scan(&u.name, &u.pass, &u.salt); err != nil {
			return
		}
		count++
	}
	if err = r.Err(); err != nil {
		return
	}
	if count > 1 {
		return nil, credential.ErrMultipleUserFound
	}

	// validate found user
	if err = a.validate(c.Secret, u); err != nil {
		return
	}

	return nil, nil
}

func (a *Authenticator) validate(secret string, user *user) error {
	if user == nil {
		return credential.ErrAuthenticateFailed
	}
	base := fmt.Sprintf("%s::%s", user.salt, secret)
	crypt := fmt.Sprintf("%x", sha256.Sum256([]byte(base)))
	if user.pass != crypt {
		return credential.ErrAuthenticateFailed
	}
	return nil
}

package dbauth

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/deadcheat/cashew/auth/credential"
)

// define singleton elements
var (
	stmt             *sql.Stmt
	attributeSpec    map[string]string
	attributeNames   []string
	attributeColumns []string
)

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
	attributeSpec = a.attributes
	attributeColumns = make([]string, len(attributeSpec))
	attributeNames = make([]string, len(attributeSpec))
	i := 0
	for k, v := range attributeSpec {
		attributeColumns[i] = k
		attributeNames[i] = v
		i++
	}
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
	for i, v := range attributeColumns {
		attributeSlice[i] = fmt.Sprintf(attributeQueryFormat, v, attributeNames[i])
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
	defer func() {
		if closeErr := r.Close(); closeErr != nil {
			log.Println("[WARN]failed to close rows: ", closeErr)
			err = closeErr
		}
	}()
	count := 0
	var nullableSalt sql.NullString
	u := new(user)
	attrValues := make([]interface{}, len(attributeNames))
	attrPointers := make([]interface{}, len(attrValues)+3)
	attrPointers[0] = &u.name
	attrPointers[1] = &u.pass
	attrPointers[2] = &nullableSalt
	for i := range attrValues {
		attrPointers[i+3] = &attrValues[i]
	}
	for r.Next() {
		if err = r.Scan(attrPointers...); err != nil {
			return nil, err
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
	if nullableSalt.Valid {
		u.salt = nullableSalt.String
	}
	attr = make(credential.Attributes)
	for i := range attrValues {
		switch t := attrValues[i].(type) {
		case []uint8:
			attr.Set(attributeNames[i], string(t))
		default:
			attr.Set(attributeNames[i], attrValues[i])
		}
	}

	return
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

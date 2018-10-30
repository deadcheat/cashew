package dbauth

import (
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestValidate(t *testing.T) {
	password := "password"

	u := user{
		name: "test",
		pass: "6917f0150a9e99a8bdbb506fe310c2d96f3190320732122b7afbfdc3898be12f",
		salt: "",
	}

	a := &Authenticator{}

	err := a.validate(password, &u)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateSelectStatement(t *testing.T) {
	db, _, _ := sqlmock.New()
	testAttributes := map[string]string{
		"testcol1": "testcolname",
		"testcol2": "testcolother",
	}
	testee := &AuthenticationBuilder{
		db:             db,
		table:          "test_table",
		userNameColumn: "test_user_name",
		passWordColumn: "test_password",
		saltColumn:     "test_salt",
		attributes:     testAttributes,
	}
	testee.Build()
	expected := `SELECT
  target.test_user_name as user,
  target.test_password as password,
  target.test_salt as salt
  , testcol1 as testcolname
  , testcol2 as testcolother
FROM 
  test_table target
WHERE
  target.test_user_name = ?
`
	actual := testee.createSelectStatement()

	if expected != actual {
		t.Errorf(`queries are not same
expected: [%s]

actual: [%s]`, expected, actual)
	}
}

package file

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/deadcheat/cashew/setting"
)

func TestNew(t *testing.T) {
	if New() == nil {
		t.Error("method New() didn't return Loader ")
	}
}

var (
	fileName      = "config.yml"
	errorFileName = "config.error.yml"
)

func TestLoad(t *testing.T) {
	expected := &setting.App{
		UseSSL:      false,
		SSLCertFile: "",
		SSLCertKey:  "",
		Host:        "localhost",
		Port:        3000,
		URIPath:     "/cas",
		Database: &setting.Database{
			Driver: "mysql",
			Name:   "casdb",
			User:   "casuser",
			Pass:   "password",
			Host:   "localhost",
			Port:   3306,
		},
		Authenticator: &setting.Authenticator{
			Driver: "database",
			AuthDatabase: &setting.AuthDatabase{
				Database: &setting.Database{
					Driver: "mysql",
					Name:   "roledb",
					User:   "roleuser",
					Pass:   "rolepass",
					Host:   "localhost",
					Port:   3306,
				},
				Table:       "auth",
				UserNameKey: "admin",
				PasswordKey: "password",
			},
			LDAP: nil,
		},
		Logging: &setting.Logging{
			Driver:   "file",
			FileName: "casserver.log",
			LogLevel: "debug",
		},
	}

	dir, _ := os.Getwd()
	path := filepath.Join(dir, fileName)
	l := new(Loader)
	s, err := l.Load(path)
	if err != nil {
		t.Errorf("test failed by error %#+v", err)
	}

	if !reflect.DeepEqual(expected, s) {
		t.Errorf(`returned value check failed \n
			expected: %#+v \n
			actual  : %#+v \n
		`, expected, s)
	}
}

func TestLoadReturnErrorWhenFileRead(t *testing.T) {
	l := new(Loader)
	_, err := l.Load("/tmp/doesnotexists/none.yml")
	if err == nil {
		t.Errorf("it should return any errors")
	}
}

func TestLoadReturnErrorWhenInvalidYaml(t *testing.T) {
	l := new(Loader)
	dir, _ := os.Getwd()
	path := filepath.Join(dir, errorFileName)
	_, err := l.Load(path)
	if err == nil {
		t.Errorf("it should return any errors")
	}

}

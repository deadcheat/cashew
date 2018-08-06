package file

import (
	"io/ioutil"
	"os"

	"github.com/deadcheat/cashew/setting"
	yaml "gopkg.in/yaml.v2"
)

// Loader struct implemented Loader interface for file
type Loader struct{}

// New return new loader
func New() setting.Loader {
	return new(Loader)
}

// Load load from file identified by filepath
func (l *Loader) Load(id string) (a *setting.App, err error) {
	file, err := os.Open(id)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// ioutil.ReadAll may not return error when file exists
	b, _ := ioutil.ReadAll(file)
	sc := setting.DefaultSetting
	a = &sc
	if err := yaml.Unmarshal(b, a); err != nil {
		return nil, err
	}
	return a, nil
}

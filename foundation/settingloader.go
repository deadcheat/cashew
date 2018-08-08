package foundation

import (
	"github.com/deadcheat/cashew/setting"
)

// Load load config and hold app-setting
func Load(configFile string) (err error) {
	app, err = loader.Load(configFile)
	return
}

// settingLoader implant setting.Loader
type settingLoader struct {
	setting.Loader
}

func newLoader(l setting.Loader) *settingLoader {
	return &settingLoader{l}
}

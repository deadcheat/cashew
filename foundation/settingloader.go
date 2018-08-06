package foundation

import (
	"github.com/deadcheat/cashew/setting"
)

var (
	// loader setting loader accessable globally
	loader *settingLoader
)

// GlobalSettingLoader access local setting loader from global
func GlobalSettingLoader() setting.Loader {
	return loader
}

// settingLoader implant setting.Loader
type settingLoader struct {
	setting.Loader
}

func newLoader(l setting.Loader) *settingLoader {
	return &settingLoader{l}
}

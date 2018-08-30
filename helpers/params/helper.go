package params

import (
	"net/url"

	"github.com/deadcheat/cashew/helpers/service"
	"github.com/deadcheat/cashew/helpers/strings"
	"github.com/deadcheat/cashew/values/consts"
)

func getURL(v url.Values, key string) (*url.URL, error) {
	values := v[key]
	u := strings.FirstString(values)
	// enable only first url
	// if serviceURL is empty, return nil, nil
	return service.NormalizeURL(u)
}

// ServiceURL get service param from url params
func ServiceURL(v url.Values) (*url.URL, error) {
	return getURL(v, consts.ParamKeyService)
}

// ContinueURL get "url" param from url params
func ContinueURL(v url.Values) (*url.URL, error) {
	return getURL(v, consts.ParamKeyURL)
}

// PgtURL get "url" param from url params
func PgtURL(v url.Values) (*url.URL, error) {
	return getURL(v, consts.ParamKeyPgtURL)
}

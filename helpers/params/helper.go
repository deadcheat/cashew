package params

import (
	"net/url"

	"github.com/deadcheat/cashew/helpers/service"
	"github.com/deadcheat/cashew/helpers/strings"
	"github.com/deadcheat/cashew/values/consts"
)

// ServiceURL get service param from url params
func ServiceURL(v url.Values) (*url.URL, error) {
	services := v[consts.ParamKeyService]
	u := strings.FirstString(services)
	// enable only first url
	// if serviceURL is empty, return nil, nil
	return service.NormalizeURL(u)
}

// ContinueURL get "url" param from url params
func ContinueURL(v url.Values) (*url.URL, error) {
	urls := v[consts.ParamKeyURL]
	u := strings.FirstString(urls)
	// enable only first url
	// if serviceURL is empty, return nil, nil
	return service.NormalizeURL(u)
}

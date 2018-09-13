package params

import (
	"net/url"

	"github.com/deadcheat/cashew/helpers/service"
	"github.com/deadcheat/cashew/helpers/strings"
	"github.com/deadcheat/cashew/values/consts"
)

func firstParam(v url.Values, key string) string {
	values := v[key]
	return strings.FirstString(values)
}

// URL get well formed url from url params
func URL(v url.Values, key string) (*url.URL, error) {
	// enable only first url
	// if serviceURL is empty, return nil, nil
	return service.NormalizeURL(firstParam(v, key))
}

// ServiceURL get service param from url params
func ServiceURL(v url.Values) (*url.URL, error) {
	return URL(v, consts.ParamKeyService)
}

// ContinueURL get "url" param from url params
func ContinueURL(v url.Values) (*url.URL, error) {
	return URL(v, consts.ParamKeyURL)
}

// PgtURL get "url" param from url params
func PgtURL(v url.Values) (*url.URL, error) {
	return URL(v, consts.ParamKeyPgtURL)
}

// Pgt get "pgt" param from url params
func Pgt(v url.Values) string {
	return firstParam(v, consts.ParamKeyPgt)
}

// TargetService get "targetService" param from url params
func TargetService(v url.Values) string {
	return firstParam(v, consts.ParamKeyTargetService)
}

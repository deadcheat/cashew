package service

import (
	"net/url"

	"github.com/PuerkitoBio/purell"
)

// NormalizeURL normalize url and unescape
func NormalizeURL(u string) (string, error) {
	urlVal, err := url.Parse(u)
	if err != nil {
		return u, err
	}
	urlVal.RawQuery = urlVal.Query().Encode()

	urlStr, err := url.QueryUnescape(urlVal.String())
	if err != nil {
		return u, err
	}
	urlStr, err = purell.NormalizeURLString(urlStr, purell.FlagsAllGreedy)
	if err != nil {
		return u, err
	}
	return urlStr, nil
}

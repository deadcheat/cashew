package service

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/purell"
)

// NormalizeURL normalize and unescape url
func NormalizeURL(u string) (*url.URL, error) {
	if len(strings.Replace(u, " ", "", -1)) == 0 {
		return nil, nil
	}
	urlStr, err := url.QueryUnescape(u)
	if err != nil {
		return nil, err
	}
	urlVal, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	urlVal.RawQuery = urlVal.Query().Encode()
	return url.Parse(purell.NormalizeURL(urlVal, purell.FlagsAllGreedy))
}

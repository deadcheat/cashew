package service

import (
	"net/url"

	"github.com/PuerkitoBio/purell"
)

// NormalizeURL normalize url and unescape
func NormalizeURL(u string) (string, error) {
	urlStr, err := url.QueryUnescape(u)
	if err != nil {
		return "", err
	}
	// fmt.Println(urlStr)
	urlVal, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	// fmt.Printf("%#v", urlVal)
	urlVal.RawQuery = urlVal.Query().Encode()

	return purell.NormalizeURL(urlVal, purell.FlagsAllGreedy), nil
}

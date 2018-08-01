package service

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/purell"
)

// NormalizeURL normalize and unescape url
func NormalizeURL(u string) (string, error) {
	if len(strings.Replace(u, " ", "", -1)) == 0 {
		return "", nil
	}
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

package service

import (
	"net/url"

	"github.com/PuerkitoBio/purell"
)

func NormalizeURL(u string) (string, error) {
	urlStr, err := url.QueryUnescape(u)
	if err != nil {
		return u, err
	}
	urlStr, err = purell.NormalizeURLString(urlStr, purell.FlagRemoveTrailingSlash|purell.FlagRemoveDuplicateSlashes|purell.FlagRemoveEmptyQuerySeparator|purell.FlagRemoveEmptyPortSeparator|purell.FlagRemoveDirectoryIndex|purell.FlagRemoveDotSegments|purell.FlagRemoveFragment)
	urlVal, err := url.Parse(urlStr)
	if err != nil {
		return u, err
	}
	return urlVal.String(), nil
}

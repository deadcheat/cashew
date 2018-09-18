package view

import (
	"html/template"
	"path/filepath"

	"github.com/deadcheat/cashew/foundation"
)

var (
	uriParser = func(path string) string {
		return filepath.Join("/", foundation.App().URIPath, path)
	}

	safe = func(u string) template.HTMLAttr {
		return template.HTMLAttr(u)
	}

	FuncMap = template.FuncMap{
		"parseURI": uriParser,
		"safe":     safe,
	}
)

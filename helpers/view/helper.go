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

	FuncMap = template.FuncMap{
		"parseURI": uriParser,
	}
)

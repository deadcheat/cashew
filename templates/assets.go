package templates

import (
	"time"

	"github.com/deadcheat/goblet"
)

//go:generate goblet -o ./assets.go -p templates files

// Assets a generated file system
var Assets = goblet.NewFS(
	map[string][]string{
		"/files": []string{
			"login", "proxy", "validate",
		},
		"/files/login": []string{
			"index.html", "logout.html",
		},
		"/files/proxy": []string{
			"proxy.xml",
		},
		"/files/validate": []string{
			"servicevalidate.xml",
		},
	},
	map[string]*goblet.File{
		"/files":                              goblet.NewFile("/files", nil, 0x800001ed, time.Unix(1537188609, 1537188609682389055)),
		"/files/login":                        goblet.NewFile("/files/login", nil, 0x800001ed, time.Unix(1535549520, 1535549520668606396)),
		"/files/login/index.html":             goblet.NewFile("/files/login/index.html", []byte(_Assetsd29f119cb8a3017c01e0ffabccfaac8408934e6f), 0x1a4, time.Unix(1537186990, 1537186990873005120)),
		"/files/login/logout.html":            goblet.NewFile("/files/login/logout.html", []byte(_Assetsdd409359cf11595b87d49b50aff81e478fecca74), 0x1a4, time.Unix(1537187025, 1537187025021967006)),
		"/files/proxy":                        goblet.NewFile("/files/proxy", nil, 0x800001ed, time.Unix(1537005745, 1537005745719603764)),
		"/files/proxy/proxy.xml":              goblet.NewFile("/files/proxy/proxy.xml", []byte(_Assets39007ceb50a14baba29c4d2d7f13cfcc597931ca), 0x1a4, time.Unix(1537005745, 1537005745719809185)),
		"/files/validate":                     goblet.NewFile("/files/validate", nil, 0x800001ed, time.Unix(1537005745, 1537005745726209597)),
		"/files/validate/servicevalidate.xml": goblet.NewFile("/files/validate/servicevalidate.xml", []byte(_Assetsdd5750ddeceedefd4ff68902c2496f670da48609), 0x1a4, time.Unix(1537005745, 1537005745726430518)),
	},
)

// binary data
var (
	_Assetsd29f119cb8a3017c01e0ffabccfaac8408934e6f = "<!DOCTYPE html>\n<html>\n\n<head>\n    <meta charset=\"utf-8\">\n    <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n    <title>login - istyle Central Login</title>\n    <link rel=\"stylesheet\" href=\"{{ parseURI \"/assets/css/font-awesome.min.css\" }}\">\n    <!-- Bulma Version 0.7.1-->\n    <link rel=\"stylesheet\" href=\"{{ parseURI \"/assets/css/bulma.min.css\" }}\" />\n    <link rel=\" stylesheet\" href=\"{{ parseURI \"/assets/css/base.css\" }}\" />\n</head>\n\n<body onload=\" if\n        (document.getElementById('username')) document.getElementById('username').focus()\">\n    <section class=\"hero is-success is-fullheight\">\n        {{ range $m := .Messages }}\n        <article class=\"message is-primary\">\n            <div class=\"message-header\">\n                <p>Message</p>\n                <button class=\"delete\" aria-label=\"delete\"></button>\n            </div>\n            <div class=\"message-body\">\n                {{ $m }}\n            </div>\n        </article>\n        <br /> {{ end }} {{ range $m := .Errors }}\n        <article class=\"message is-danger\">\n            <div class=\"message-header\">\n                <p>Error</p>\n                <button class=\"delete\" aria-label=\"delete\"></button>\n            </div>\n            <div class=\"message-body\">\n                {{ $m }}\n            </div>\n        </article><br /> {{ end }}\n        <div class=\"hero-body\">\n            <div class=\"container has-text-centered\">\n                {{ if not .LoggedIn }}\n                <div class=\"column is-4 is-offset-4\">\n                    <h3 class=\"title has-text-grey\">istyle Central Login</h3>\n                    <p class=\"subtitle has-text-grey\">Please login to proceed.</p>\n                    <div class=\"box\">\n                        <form method=\"post\" action=\"{{ parseURI \"/login\" }}\">\n                            <div class=\"field\">\n                                <div class=\"control\">\n                                    <input class=\"input is-large\" name=\"username\" type=\"text\" placeholder=\"Your User Name\"\n                                        autofocus=\"\" tabindex=\"1\" accesskey=\"u\" value=\"{{ .UserName }}\">\n                                </div>\n                            </div>\n\n                            <div class=\"field\">\n                                <div class=\"control\">\n                                    <input class=\"input is-large\" type=\"password\" name=\"password\" placeholder=\"Your Password\"\n                                        tabindex=\"2\" accesskey=\"p\" autocomplete=\"off\" value=\"{{ .Password }}\">\n                                </div>\n                            </div>\n                            <div class=\"field\">\n                                <label class=\"checkbox\">\n                                    <input type=\"checkbox\"> Remember me\n                                </label>\n                            </div>\n                            <input type=\"hidden\" name=\"lt\" value=\"{{ .LoginTicket }}\">\n                            <input type=\"hidden\" name=\"service\" value=\"{{ .Service }}\">\n                            <button class=\"button is-block is-info is-large is-fullwidth\">Login</button>\n                        </form>\n                    </div>\n                </div>\n                {{ else }}\n                <img src=\"{{ parseURI \"/assets/images/success.jpg\" }}\" alt=\"\"> {{ end }}\n            </div>\n        </div>\n    </section>\n</body>\n\n</html>"
	_Assetsdd409359cf11595b87d49b50aff81e478fecca74 = "<!DOCTYPE html>\n<html>\n\n<head>\n    <meta charset=\"utf-8\">\n    <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n    <title>logout - istyle Central Login</title>\n    <link rel=\"stylesheet\" href=\"{{ parseURI \"/assets/css/font-awesome.min.css\" }}\">\n    <!-- Bulma Version 0.7.1-->\n    <link rel=\"stylesheet\" href=\"{{ parseURI \"/assets/css/bulma.min.css\" }}\" />\n    <link rel=\" stylesheet\" href=\"{{ parseURI \"/assets/css/base.css\" }}\" />\n</head>\n\n<body>\n    <section class=\"hero is-success is-fullheight\">\n        <div class=\"hero-body\">\n            <div class=\"column is-4 is-offset-4\">\n                <h3 class=\"title has-text-grey\">istyle Central Login</h3>\n                <div class=\"container has-text-centered\">\n                    <article class=\"message is-info\">\n                        <div class=\"message-body\">\n                            You have successfully logged out. Please click on the following link to continue:\n                        </div>\n                    </article>\n                    <a class=\"is-size-3 has-text-grey\" href=\"{{ .Next }}\">{{ .Next }}</a>\n                </div>\n            </div>\n        </div>\n    </section>\n</body>\n\n</html>"

	_Assets39007ceb50a14baba29c4d2d7f13cfcc597931ca = "<cas:serviceResponse xmlns:cas='http://www.yale.edu/tp/cas'>{{ if .Success }}\n    <cas:proxySuccess>\n        <cas:proxyTicket>{{ .Ticket }}</cas:proxyTicket>\n    </cas:proxySuccess>{{ else }}\n    <cas:proxyFailure code=\"{{ .ErrorCode }}\">\n        {{ .ErrorBody }}\n    </cas:proxyFailure>{{ end }}\n</cas:serviceResponse>\n"

	_Assetsdd5750ddeceedefd4ff68902c2496f670da48609 = "<cas:serviceResponse xmlns:cas='http://www.yale.edu/tp/cas'>{{ if .Success }}\n    <cas:authenticationSuccess>\n        <cas:user>{{ .Name }}</cas:user>{{ if ne .IOU \"\" }}\n        <cas:proxyGrantingTicket>{{ .IOU }}</cas:proxyGrantingTicket>{{ end }}\n        {{ if gt (len .Proxies) 0 }}\n        <cas:proxies>\n        {{ range $p := .Proxies }}\n        <cas:proxy>{{ $p }}</cas:proxy>\n        {{ end }}\n        </cas:proxies>{{ end }}\n    </cas:authenticationSuccess>{{ else }}\n    <cas:authenticationFailure code=\"{{ .ErrorCode }}\">\n        {{ .ErrorBody }}\n    </cas:authenticationFailure>{{ end }}\n</cas:serviceResponse>\n"
)


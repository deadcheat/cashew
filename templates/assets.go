package templates

import (
	"time"

	"github.com/deadcheat/goblet"
)

//go:generate goblet -g -o assets.go -p templates files

// Assets a generated file system
var Assets = goblet.NewFS(
	map[string][]string{
		"/files": []string{
			"login", "proxy", "validate",
		},
		"/files/login": []string{
			"login.html", "logout.html",
		},
		"/files/proxy": []string{
			"proxy.xml",
		},
		"/files/validate": []string{
			"servicevalidate.xml",
		},
	},
	map[string]*goblet.File{
		"/files":                              goblet.NewFile("/files", nil, 0x800001ed, time.Unix(1537798571, 1537798571007687686)),
		"/files/login":                        goblet.NewFile("/files/login", nil, 0x800001ed, time.Unix(1547861190, 1547861190419318184)),
		"/files/login/login.html":             goblet.NewFile("/files/login/login.html", []byte(_Assetsde4243da79c0324a6719ff4a9316d4a31ecacdf9), 0x1a4, time.Unix(1547866214, 1547866214450370518)),
		"/files/login/logout.html":            goblet.NewFile("/files/login/logout.html", []byte(_Assetsdd409359cf11595b87d49b50aff81e478fecca74), 0x1a4, time.Unix(1547689508, 1547689508098883131)),
		"/files/proxy":                        goblet.NewFile("/files/proxy", nil, 0x800001ed, time.Unix(1537798571, 1537798571007340026)),
		"/files/proxy/proxy.xml":              goblet.NewFile("/files/proxy/proxy.xml", []byte(_Assets39007ceb50a14baba29c4d2d7f13cfcc597931ca), 0x1a4, time.Unix(1537798571, 1537798571007442126)),
		"/files/validate":                     goblet.NewFile("/files/validate", nil, 0x800001ed, time.Unix(1537798571, 1537798571007995817)),
		"/files/validate/servicevalidate.xml": goblet.NewFile("/files/validate/servicevalidate.xml", []byte(_Assetsdd5750ddeceedefd4ff68902c2496f670da48609), 0x1a4, time.Unix(1537798571, 1537798571009295560)),
	},
)

// binary data
var (
	_Assetsde4243da79c0324a6719ff4a9316d4a31ecacdf9 = "<!DOCTYPE html>\n<html>\n\n<head>\n    <meta charset=\"utf-8\">\n    <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n    <title>Login - {{.Organization}} Central Authentication Service</title>\n    <script src=\"{{ parseURI \"/assets/js/cashew.js\" }}\"></script>\n    <link rel=\"stylesheet\" href=\"{{ parseURI \"/assets/css/font-awesome.min.css\" }}\">\n    <!-- Bulma Version 0.7.1-->\n    <link rel=\"stylesheet\" href=\"{{ parseURI \"/assets/css/bulma.min.css\" }}\" />\n    <link rel=\" stylesheet\" href=\"{{ parseURI \"/assets/css/base.css\" }}\" />\n</head>\n\n<body onload=\"if\n        (document.getElementById('username')) document.getElementById('username').focus()\">\n    <section class=\"hero is-success is-fullheight\">\n        {{ if .Messages }}\n        <article id=\"messageArea\" class=\"message is-primary\">\n            <div class=\"message-header\">\n                <p>Message</p>\n                <button class=\"delete\" aria-label=\"delete\" onclick=\"hide('messageArea');\"></button>\n            </div>\n            {{ range $m := .Messages }}\n            <div class=\"message-body\">\n                {{ $m }}\n            </div>{{ end }}\n        </article>\n        <br />{{ end }} {{ range $m := .Errors }}\n        <article id=\"errorArea\" class=\"message is-danger\">\n            <div class=\"message-header\">\n                <p>Error</p>\n                <button class=\"delete\" aria-label=\"delete\" onclick=\"hide('errorArea');\"></button>\n            </div>\n            <div class=\"message-body\">\n                {{ $m }}\n            </div>\n        </article><br /> {{ end }}\n        <div class=\"hero-body\">\n            <div class=\"container has-text-centered\">\n                <div class=\"column is-4 is-offset-4\">\n                    <h3 class=\"title has-text-grey\">{{.Organization}} Central Authentication Service Login</h3>\n                    <p class=\"subtitle has-text-grey\">Please login to proceed.</p>\n                    <div class=\"box\">\n                        <form method=\"post\" action=\"{{ parseURI \"/login\" }}\">\n                            <div class=\"field\">\n                                <div class=\"control\">\n                                    <input class=\"input is-large\" name=\"username\" type=\"text\" placeholder=\"Your User Name\"\n                                        autofocus=\"\" tabindex=\"1\" accesskey=\"u\" value=\"{{ .UserName }}\">\n                                </div>\n                            </div>\n\n                            <div class=\"field\">\n                                <div class=\"control\">\n                                    <input class=\"input is-large\" type=\"password\" name=\"password\" placeholder=\"Your Password\"\n                                        tabindex=\"2\" accesskey=\"p\" autocomplete=\"off\" value=\"{{ .Password }}\">\n                                </div>\n                            </div>\n                            <!-- this checkbox does not work yet\n                            <div class=\"field\">\n                                <label class=\"checkbox\">\n                                    <input type=\"checkbox\"> Remember me\n                                </label>\n                            </div>\n                            -->\n                            <input type=\"hidden\" name=\"lt\" value=\"{{ .LoginTicket }}\">\n                            <input type=\"hidden\" name=\"service\" value=\"{{ .Service }}\">\n                            <button class=\"button is-block is-info is-large is-fullwidth\">Login</button>\n                        </form>\n                    </div>\n                </div>\n            </div>\n        </div>\n    </section>\n</body>\n\n</html>"
	_Assetsdd409359cf11595b87d49b50aff81e478fecca74 = "<!DOCTYPE html>\n<html>\n\n<head>\n    <meta charset=\"utf-8\">\n    <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n    <title>logout - {{.Organization}} Central Authentication Service</title>\n    <link rel=\"stylesheet\" href=\"{{ parseURI \"/assets/css/font-awesome.min.css\" }}\">\n    <!-- Bulma Version 0.7.1-->\n    <link rel=\"stylesheet\" href=\"{{ parseURI \"/assets/css/bulma.min.css\" }}\" />\n    <link rel=\" stylesheet\" href=\"{{ parseURI \"/assets/css/base.css\" }}\" />\n</head>\n\n<body>\n    <section class=\"hero is-success is-fullheight\">\n        <div class=\"hero-body\">\n            <div class=\"column is-4 is-offset-4\">\n                <h3 class=\"title has-text-grey\">{{.Organization}} Central Authentication Service Login</h3>\n                <div class=\"container has-text-centered\">\n                    <article class=\"message is-info\">\n                        <div class=\"message-body\">\n                            You have successfully logged out. Please click on the following link to continue:\n                        </div>\n                    </article>\n                    <a class=\"is-size-3 has-text-grey\" href=\"{{ .Next }}\">{{ .Next }}</a>\n                </div>\n            </div>\n        </div>\n    </section>\n</body>\n\n</html>"

	_Assets39007ceb50a14baba29c4d2d7f13cfcc597931ca = "<cas:serviceResponse xmlns:cas='http://www.yale.edu/tp/cas'>{{ if .Success }}\n    <cas:proxySuccess>\n        <cas:proxyTicket>{{ .Ticket }}</cas:proxyTicket>\n    </cas:proxySuccess>{{ else }}\n    <cas:proxyFailure code=\"{{ .ErrorCode }}\">\n        {{ .ErrorBody }}\n    </cas:proxyFailure>{{ end }}\n</cas:serviceResponse>\n"

	_Assetsdd5750ddeceedefd4ff68902c2496f670da48609 = "<cas:serviceResponse xmlns:cas='http://www.yale.edu/tp/cas'>{{ if .Success }}\n    <cas:authenticationSuccess>\n        <cas:user>{{ .Name }}</cas:user>{{ if ne .IOU \"\" }}\n        <cas:proxyGrantingTicket>{{ .IOU }}</cas:proxyGrantingTicket>{{ end }}\n        {{ if gt (len .Proxies) 0 }}\n        <cas:proxies>\n        {{ range $p := .Proxies }}\n        <cas:proxy>{{ $p }}</cas:proxy>\n        {{ end }}\n        </cas:proxies>{{ end }}\n        {{ if gt (len .ExtraAttributes) 0 }}\n        <cas:attributes>\n        {{ range $index, $element := .ExtraAttributes }}\n        <cas:{{ $index |safe }}>{{ $element }}</cas:{{ $index |safe }}>\n        {{ end }}\n        </cas:attributes>{{ end }}\n    </cas:authenticationSuccess>{{ else }}\n    <cas:authenticationFailure code=\"{{ .ErrorCode }}\">\n        {{ .ErrorBody }}\n    </cas:authenticationFailure>{{ end }}\n</cas:serviceResponse>\n"
)


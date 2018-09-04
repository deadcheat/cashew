package proxycallback

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/foundation"
	"github.com/deadcheat/cashew/values/errs"
)

// Repository implements cashew.ProxyCallBackRepository
type Repository struct{}

// New return new ProxyCallBackRepository
func New() cashew.ProxyCallBackRepository {
	return new(Repository)
}

// Dial request proxy url and handle with status code
func (r *Repository) Dial(u *url.URL, pgt, iou string) (err error) {
	q := u.Query()
	q.Set("pgtId", pgt)
	q.Set("pgtIou", iou)
	u.RawQuery = q.Encode()

	c := new(http.Client)
	if strings.ToLower(u.Scheme) == "https" {
		tlsConfig := new(tls.Config)
		if foundation.App().UseSSL {
			var cer tls.Certificate
			if cer, err = tls.LoadX509KeyPair(foundation.App().SSLCertFile, foundation.App().SSLCertKey); err != nil {
				return
			}
			tlsConfig.Certificates = []tls.Certificate{cer}
		} else {
			tlsConfig.InsecureSkipVerify = true
		}
		c.Transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	var req *http.Request
	if req, err = http.NewRequest(http.MethodGet, u.String(), nil); err != nil {
		return
	}

	var resp *http.Response
	if resp, err = c.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	log.Printf("URL %s returned %d(%s) status", u.String(), resp.StatusCode, resp.Status)
	switch resp.StatusCode {
	case http.StatusOK,
		http.StatusAccepted,
		http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusNotModified:
		return
	default:
		return errs.ErrProxyGrantingURLUnexpectedStatus
	}
}

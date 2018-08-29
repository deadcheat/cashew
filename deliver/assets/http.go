package assets

import (
	"net/http"
	"path/filepath"

	"github.com/deadcheat/goblet"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/assets"
	"github.com/deadcheat/cashew/foundation"

	"github.com/gorilla/mux"
)

// Deliver struct implements cashew.Deliver
type Deliver struct {
	r  *mux.Router
	fs *goblet.FileSystem
}

// New return new deliver interface
func New(r *mux.Router) cashew.Deliver {
	fs := assets.Assets.WithPrefix("/assets/").WithIgnoredPrefix("/files/")
	return &Deliver{r, fs}
}

// Mount mounts assetfiles
func (d *Deliver) Mount() {
	d.r.PathPrefix("/assets/").Handler(http.StripPrefix(filepath.Join(foundation.App().URIPath, "/assets/"), http.FileServer(d.fs)))
}

package main

import (
	"context"
	"net/http"

	"github.com/traefik/yaegi/interp"
)

// BackendIP is a struct that holds the configuration for the plugin.
type BackendIP struct {
	HeaderName string `json:"headerName,omitempty"`
}

// New creates a new instance of the plugin with the given configuration.
func New(_ context.Context, next http.Handler, config *BackendIP, _ string) (http.Handler, error) {
	if config.HeaderName == "" {
		config.HeaderName = "X-Backend-Server"
	}

	return &backendIPHandler{
		next:       next,
		headerName: config.HeaderName,
	}, nil
}

type backendIPHandler struct {
	next       http.Handler
	headerName string
}

func (h *backendIPHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	backendIP := req.RemoteAddr
	rw.Header().Set(h.headerName, backendIP)
	h.next.ServeHTTP(rw, req)
}

// No vendor, no version, no plugin module path.
const pluginReference = `
package main

func init() {
	// plugin is the name of the exported symbol in the Go plugin.
	// It has to be an object implementing the interface "github.com/traefik/traefik/v2/pkg/config/pluginconfig/plugin.Plugin".
	// Plugin is an alias for that interface.
	// The symbol has to be named "plugin" (all lowercase).
	plugin = New
}
`

func init() {
	err := interp.AddModule(interp.Module{Name: "github.com/0xgeert/traefeik_reponseip_plugin", Content: pluginReference})
	if err != nil {
		panic(err)
	}
}

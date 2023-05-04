package plugins

import (
	"net/http"
)

// simple plugin interface
type PluginInOut interface {
	API() // URL https://auth.estuary.tech/

	// structure of the input must be defined
	Input(interface{}) error

	// structure of the output must be defined
	Output() interface{} // can be an error
}

type Plugin interface {
	Name() string
	Initialize() error
	Middleware() func(http.Handler) http.Handler
}

type PluginRegistry struct {
	plugins map[string]Plugin
}

// flow
// 1. load all the plugins in this directory
// 2. override pre-existing plugin with the plugins from this directory
// 3. run the plugins

package plugins

import (
	"fmt"
	"net/http"
	"plugin"
)

func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		plugins: make(map[string]Plugin),
	}
}

func (r *PluginRegistry) Register(plugin Plugin) {
	r.plugins[plugin.Name()] = plugin
}

func (r *PluginRegistry) Middleware(name string) (func(http.Handler) http.Handler, bool) {
	plugin, ok := r.plugins[name]
	if !ok {
		return nil, false
	}
	return plugin.Middleware(), true
}

func LoadPlugin(path string, registry *PluginRegistry) error {
	plug, err := plugin.Open(path)
	if err != nil {
		return err
	}

	symPlugin, err := plug.Lookup("PluginInstance")
	if err != nil {
		return err
	}

	pluginInstance, ok := symPlugin.(Plugin)
	if !ok {
		return fmt.Errorf("unexpected plugin type: %T", symPlugin)
	}

	registry.Register(pluginInstance)
	return pluginInstance.Initialize()
}

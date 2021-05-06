package drpc

type Plugin interface {
	Name() string
}

type PluginContainer struct {
}

func newPluginContainer() *PluginContainer {
	p := &PluginContainer{}
	return p
}

func (that *PluginContainer) cloneAndAppendMiddle(plugins ...Plugin) *PluginContainer {
	return that
}

func warnInvalidHandlerHooks(plugin []Plugin) {

}

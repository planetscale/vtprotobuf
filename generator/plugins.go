package generator

import (
	"sort"

	"google.golang.org/protobuf/compiler/protogen"
)

var defaultPlugins []PluginFactory

func pluginsForFile(file *GeneratedFile) []Plugin {
	plugins := make([]Plugin, 0, len(defaultPlugins))

	for _, pg := range defaultPlugins {
		plugins = append(plugins, pg(file))
	}

	// Sort the slice stably so the contents of the generated files don't
	// "jump around" when enabling or disabling more plugins
	sort.SliceStable(plugins, func(i, j int) bool {
		return plugins[i].Name() < plugins[j].Name()
	})

	return plugins
}

func RegisterPlugin(plugin PluginFactory) {
	defaultPlugins = append(defaultPlugins, plugin)
}

type PluginFactory func(gen *GeneratedFile) Plugin

type Plugin interface {
	Name() string
	GenerateFile(file *protogen.File) bool
	GenerateHelpers()
}

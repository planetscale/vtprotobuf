package generator

import (
	"fmt"
	"sort"

	"google.golang.org/protobuf/compiler/protogen"
)

var defaultPlugins = make(map[string]PluginFactory)

func findPlugins(features []string) ([]PluginFactory, error) {
	sort.Strings(features)

	plugins := make([]PluginFactory, 0, len(defaultPlugins))
	for _, feature := range features {
		if feature == "all" {
			plugins = plugins[:0]
			for _, pg := range defaultPlugins {
				plugins = append(plugins, pg)
			}
			break
		}

		pg, ok := defaultPlugins[feature]
		if !ok {
			return nil, fmt.Errorf("unknown feature: %q", feature)
		}
		plugins = append(plugins, pg)
	}
	return plugins, nil
}

func RegisterPlugin(name string, plugin PluginFactory) {
	defaultPlugins[name] = plugin
}

type PluginFactory func(gen *GeneratedFile) Plugin

type Plugin interface {
	GenerateFile(file *protogen.File) bool
	GenerateHelpers()
}

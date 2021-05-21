package generator

import (
	"fmt"
	"sort"

	"google.golang.org/protobuf/compiler/protogen"
)

var defaultPlugins = make(map[string]PluginFactory)

type sortedPlugin struct {
	f PluginFactory
	n string
}

func findPlugins(features []string) ([]PluginFactory, error) {
	var sorted []sortedPlugin
	for _, name := range features {
		if name == "all" {
			sorted = sorted[:0]
			for n, pg := range defaultPlugins {
				sorted = append(sorted, sortedPlugin{f: pg, n: n})
			}
			break
		}

		pg, ok := defaultPlugins[name]
		if !ok {
			return nil, fmt.Errorf("unknown feature: %q", name)
		}
		sorted = append(sorted, sortedPlugin{f: pg, n: name})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].n < sorted[j].n
	})

	var plugins []PluginFactory
	for _, sp := range sorted {
		plugins = append(plugins, sp.f)
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

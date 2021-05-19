package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/planetscale/vtprotobuf/generator"
	_ "github.com/planetscale/vtprotobuf/plugins/marshal"
	_ "github.com/planetscale/vtprotobuf/plugins/pool"
	_ "github.com/planetscale/vtprotobuf/plugins/size"
	_ "github.com/planetscale/vtprotobuf/plugins/unmarshal"

	"google.golang.org/protobuf/compiler/protogen"
)

type ObjectSet map[protogen.GoIdent]bool

func (o ObjectSet) String() string {
	return fmt.Sprintf("%#v", o)
}

func (o ObjectSet) Set(s string) error {
	idx := strings.LastIndexByte(s, '.')
	if idx < 0 {
		return fmt.Errorf("invalid object name: %q", s)
	}

	ident := protogen.GoIdent{
		GoImportPath: protogen.GoImportPath(s[0:idx]),
		GoName:       s[idx+1:],
	}
	o[ident] = true
	return nil
}

func main() {
	poolable := make(ObjectSet)
	var f flag.FlagSet
	f.Var(poolable, "P", "use memory pooling for this object")
	protogen.Options{ParamFunc: f.Set}.Run(func(plugin *protogen.Plugin) error {
		generateAllFiles(plugin, poolable)
		return nil
	})
}

func generateAllFiles(plugin *protogen.Plugin, poolable ObjectSet) {
	ext := &generator.Extensions{Poolable: poolable}
	gen := generator.NewGenerator(ext)

	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}

		gf := plugin.NewGeneratedFile(file.GeneratedFilenamePrefix+"_vtproto.pb.go", file.GoImportPath)
		if !gen.GenerateFile(gf, file) {
			gf.Skip()
		}
	}
}

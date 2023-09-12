package main

import (
	"flag"
	"strings"

	_ "github.com/planetscale/vtprotobuf/features/clone"
	_ "github.com/planetscale/vtprotobuf/features/equal"
	_ "github.com/planetscale/vtprotobuf/features/grpc"
	_ "github.com/planetscale/vtprotobuf/features/marshal"
	_ "github.com/planetscale/vtprotobuf/features/pool"
	_ "github.com/planetscale/vtprotobuf/features/size"
	_ "github.com/planetscale/vtprotobuf/features/unmarshal"
	"github.com/planetscale/vtprotobuf/generator"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	// os.WriteFile("/tmp/go_dbg_pid", []byte(strconv.Itoa(os.Getpid())), 0666)
	// sig := make(chan os.Signal, 1)
	// signal.Notify(sig, syscall.SIGCONT)
	// <-sig

	var features string
	ext := &generator.Extensions{
		Poolable: make(generator.ObjectSet),
	}

	var f flag.FlagSet
	f.StringVar(&features, "features", "all", "list of features to generate (separated by '+')")
	f.BoolVar(&ext.Foreign, "freestanding", false, "")
	f.BoolVar(&ext.AllowEmpty, "allow-empty", false, "allow generation of empty files")
	f.Var(ext.Poolable, "pool", "use memory pooling for this object")

	protogen.Options{ParamFunc: f.Set}.Run(func(plugin *protogen.Plugin) error {
		return generateAllFiles(plugin, strings.Split(features, "+"), ext)
	})
}

var SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

func generateAllFiles(plugin *protogen.Plugin, featureNames []string, ext *generator.Extensions) error {
	gen, err := generator.NewGenerator(plugin.Files, featureNames, ext)
	if err != nil {
		return err
	}

	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}

		var importPath protogen.GoImportPath
		if !ext.Foreign {
			importPath = file.GoImportPath
		}

		gf := plugin.NewGeneratedFile(file.GeneratedFilenamePrefix+"_vtproto.pb.go", importPath)
		if !gen.GenerateFile(gf, file) && !ext.AllowEmpty {
			gf.Skip()
		}
	}

	plugin.SupportedFeatures = SupportedFeatures
	return nil
}

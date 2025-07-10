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
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	var cfg generator.Config
	var features string
	var f flag.FlagSet

	f.BoolVar(&cfg.AllowEmpty, "allow-empty", false, "allow generation of empty files")
	cfg.Poolable = generator.NewObjectSet()
	cfg.PoolableExclude = generator.NewObjectSet()
	cfg.IgnoreUnknownFields = generator.NewObjectSet()
	f.Var(&cfg.Poolable, "pool", "use memory pooling for this object")
	f.Var(&cfg.PoolableExclude, "pool-exclude", "do not use memory pooling for this object")
	f.Var(&cfg.IgnoreUnknownFields, "ignoreUnknownFields", "ignore unknown fields instead of saving them")
	f.BoolVar(&cfg.Wrap, "wrap", false, "generate wrapper types")
	f.StringVar(&features, "features", "all", "list of features to generate (separated by '+')")
	f.StringVar(&cfg.BuildTag, "buildTag", "", "the go:build tag to set on generated files")

	protogen.Options{ParamFunc: f.Set}.Run(func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)
		plugin.SupportedEditionsMinimum = descriptorpb.Edition_EDITION_PROTO3
		plugin.SupportedEditionsMaximum = descriptorpb.Edition_EDITION_MAX

		gen, err := generator.NewGenerator(plugin, strings.Split(features, "+"), &cfg)
		if err != nil {
			return err
		}
		gen.Generate()
		return nil
	})
}

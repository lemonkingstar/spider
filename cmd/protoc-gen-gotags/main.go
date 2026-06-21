package main

import (
	"github.com/lemonkingstar/spider/cmd/protoc-gen-gotags/handler"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	protogen.Options{}.Run(
		func(plugin *protogen.Plugin) error {
			plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
			return handler.NewTagHandler(plugin).Execute()
		},
	)
}

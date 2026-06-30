package handler

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"

	"github.com/fatih/structtag"
	"github.com/lemonkingstar/spider/cmd/protoc-gen-gotags/tagger"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

type tagHandler struct {
	plugin *protogen.Plugin
}

func NewTagHandler(plugin *protogen.Plugin) Handler {
	return &tagHandler{plugin: plugin}
}

func (*tagHandler) Name() string {
	return "go-tags"
}

func (h *tagHandler) Execute() error {
	for _, file := range h.plugin.Files {
		if !file.Generate {
			continue
		}
		filename := file.GeneratedFilenamePrefix + ".pb.go"
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			_, _ = fmt.Fprintf(os.Stderr, "WARN: generated file %q not found on disk\n", filename)
			continue
		}

		fs := token.NewFileSet()
		f, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parse %s error: %w", filename, err)
		}

		fieldIndex := buildFieldIndex(file.Messages)
		ast.Inspect(f, func(n ast.Node) bool {
			typeSpec, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				return true
			}
			msgInfo, exists := fieldIndex[typeSpec.Name.Name]
			if !exists {
				return true
			}

			for _, astField := range structType.Fields.List {
				if len(astField.Names) == 0 || astField.Tag == nil {
					continue
				}

				goFieldName := astField.Names[0].Name
				if protoField, found := msgInfo.fields[goFieldName]; found {
					if err = replaceTag(astField, protoField, nil); err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "WARN: replace tag failed for %s.%s: %v\n",
							typeSpec.Name.Name, goFieldName, err)
					}
				} else if protoOneof, found := msgInfo.oneofs[goFieldName]; found {
					if err = replaceTag(astField, nil, protoOneof); err != nil {
						_, _ = fmt.Fprintf(os.Stderr, "WARN: replace tag failed for %s.%s: %v\n",
							typeSpec.Name.Name, goFieldName, err)
					}
				}
			}
			return true
		})

		gf := h.plugin.NewGeneratedFile(filename, file.GoImportPath)
		var buf bytes.Buffer
		err = printer.Fprint(&buf, fs, f)
		if err != nil {
			return err
		}
		gf.P(buf.String())
	}
	return nil
}

func replaceTag(astField *ast.Field, protoField *protogen.Field, protoOneof *protogen.Oneof) error {
	if astField.Tag == nil {
		return nil
	}

	rawTag := strings.Trim(astField.Tag.Value, "`")
	tags, err := structtag.Parse(rawTag)
	if err != nil {
		return err
	}

	modified := false
	var strTags []string
	if protoField != nil {
		strTags = proto.GetExtension(protoField.Desc.Options(), tagger.E_Tags).([]string)
		disableOmitempty, _ := proto.GetExtension(protoField.Parent.Desc.Options(), tagger.E_DisableOmitempty).(bool)
		if disableOmitempty {
			tags.DeleteOptions("json", "omitempty")
			modified = true
		}
	} else if protoOneof != nil {
		strTags = proto.GetExtension(protoOneof.Desc.Options(), tagger.E_OneofTags).([]string)
		disableOmitempty, _ := proto.GetExtension(protoOneof.Parent.Desc.Options(), tagger.E_DisableOmitempty).(bool)
		if disableOmitempty {
			tags.DeleteOptions("json", "omitempty")
			modified = true
		}
	}
	for _, t := range strTags {
		pair := strings.SplitN(t, ":", 2)
		if len(pair) != 2 {
			return fmt.Errorf("invalid tag format %s", t)
		}
		key := pair[0]
		value := strings.Trim(pair[1], `"`)
		if err = tags.Set(&structtag.Tag{Key: key, Name: value}); err != nil {
			return err
		}
		modified = true
	}

	if modified {
		astField.Tag.Value = "`" + tags.String() + "`"
	}
	return nil
}

package main

import (
	"regexp"
	"strings"

	"github.com/mattrx/protoc-gen-catalog/model"
	"github.com/mattrx/protoc-gen-catalog/renderer"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
)

var filterRegex *regexp.Regexp = regexp.MustCompile(".*")

func main() {
	protogen.Options{
		ParamFunc: func(name, value string) error {
			switch name {
			case "filter":
				r, err := regexp.Compile(value)
				if err != nil {
					return err
				}
				filterRegex = r
			}

			return nil
		},
	}.Run(func(gen *protogen.Plugin) error {
		entities := map[string]model.Entity{}

		for _, file := range gen.Files {
			for _, message := range file.Messages {
				collect(entities, message)
			}
		}

		for _, e := range entities {
			if e, ok := e.(*model.Message); ok {
				for _, f := range e.GetFields() {
					if e1, ok := entities[f.GetKind()]; ok {
						f.SetEntity(e1.(model.Entity))
						e1.AddUsage(f)
					}
				}
			}
		}

		return renderer.Do(gen, entities, filterRegex)
	})
}

func collect(entityMap map[string]model.Entity, m *protogen.Message) {

	for _, message := range m.Messages {
		collect(entityMap, message)
	}

	for i := 0; i < m.Desc.Enums().Len(); i++ {
		enum := model.NewEnum(
			string(m.Desc.Enums().Get(i).FullName()),
			strings.Trim(strings.ReplaceAll(m.Comments.Leading.String(), "//", "<br />"), "<br />"),
		)

		if opts, ok := m.Desc.Options().(*descriptorpb.MessageOptions); ok {
			enum.SetDeprecated(opts.GetDeprecated())
		}

		for j := 0; j < m.Desc.Enums().Get(i).Values().Len(); j++ {
			enum.AddValue(string(m.Desc.Enums().Get(i).Values().Get(j).Name()))
		}

		entityMap[enum.GetName()] = enum
	}

	if len(m.Fields) > 0 {
		message := model.NewMessage(
			string(m.Desc.FullName()),
			strings.Trim(strings.ReplaceAll(m.Comments.Leading.String(), "//", "<br />"), "<br />"),
		)

		if opts, ok := m.Desc.Options().(*descriptorpb.MessageOptions); ok {
			message.SetDeprecated(opts.GetDeprecated())
		}

		for _, field := range m.Fields {

			f := model.NewField(
				string(field.Desc.Name()),
				strings.Trim(strings.ReplaceAll(field.Comments.Leading.String(), "//", "<br />"), "<br />"),
			)

			f.SetKind(field.Desc.Kind().String())

			if f.GetKind() == "message" {
				f.SetKind(string(field.Desc.Message().FullName()))
			}

			if field.Enum != nil {
				f.SetKind(string(field.Enum.Desc.FullName()))
			}

			if opts, ok := field.Desc.Options().(*descriptorpb.FieldOptions); ok {
				f.SetDeprecated(opts.GetDeprecated())
			}

			if field.Desc.IsList() {
				f.AddLabel("list")
			}

			if field.Desc.IsMap() {
				f.AddLabel("map")
			}

			message.AddField(f)

		}

		entityMap[message.GetName()] = message
	}
}

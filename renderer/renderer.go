package renderer

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/mattrx/protoc-gen-catalog/model"
	"google.golang.org/protobuf/compiler/protogen"
)

func Do(gen *protogen.Plugin, entityMap map[string]model.Entity, filter *regexp.Regexp) error {

	f := gen.NewGeneratedFile("index.html", "")

	entities := []model.Entity{}

	for _, e := range entityMap {
		entities = append(entities, e)
	}

	sort.SliceStable(entities, func(i, j int) bool {
		iRunes := []rune(entities[i].GetName())
		jRunes := []rune(entities[j].GetName())

		max := len(iRunes)
		if max > len(jRunes) {
			max = len(jRunes)
		}

		for idx := 0; idx < max; idx++ {
			ir := iRunes[idx]
			jr := jRunes[idx]

			lir := unicode.ToLower(ir)
			ljr := unicode.ToLower(jr)

			if lir != ljr {
				return lir < ljr
			}

			if ir != jr {
				return ir < jr
			}
		}

		return len(iRunes) < len(jRunes)
	})

	navigation := &strings.Builder{}
	main := &strings.Builder{}
	for _, entity := range entities {
		if !entity.IsUsed() {
			if filter.MatchString(entity.GetName()) {
				if entity.IsDeprecated() {
					navigation.WriteString("<li><s><a href='#" + makeLink([]string{entity.GetName()}) + "'>" + entity.GetName() + "</a></s></li>")
				} else {
					navigation.WriteString("<li><a href='#" + makeLink([]string{entity.GetName()}) + "'>" + entity.GetName() + "</a></li>")
				}
				main.WriteString(renderEntity(f, entity, []string{}).String())
			}
		}
	}
	f.P("<html>")
	f.P("<head>")
	f.P("<style>" + stylesheet + "</style>")
	f.P("</head>")
	f.P("<body>")

	f.P("<nav id='toc'><ul>")
	f.P(navigation.String())
	f.P("</ul></nav>")

	f.P("<section id='main'>")
	f.P(main.String())
	f.P("</section>")

	f.P("</body>")
	f.P("</html>")

	return nil

}

func renderEntity(f *protogen.GeneratedFile, e model.Entity, parents []string) *strings.Builder {

	sb := &strings.Builder{}

	for _, p := range parents {
		if e.GetName() == p {
			return sb
		}
	}

	sb.WriteString("<details class='entity level-" + strconv.Itoa(len(parents)) + "' id='" + makeLink(append(parents, e.GetName())) + "'>")

	sb.WriteString("<summary class='title'>")
	sb.WriteString(e.GetName())
	if _, ok := e.(*model.Enum); ok {
		sb.WriteString("<span class='label'>enum</span>")
	}
	if e.IsDeprecated() {
		sb.WriteString("<span class='label'>deprecated</span>")
	}
	sb.WriteString("</summary>")

	if e.GetDescription() != "" {
		sb.WriteString("<div class='description'>")
		sb.WriteString(e.GetDescription())
		sb.WriteString("</div>")
	}

	if e, ok := e.(*model.Enum); ok {
		sb.WriteString("<ul>")
		for _, value := range e.GetValues() {
			sb.WriteString("<li>")
			sb.WriteString(value)
			sb.WriteString("</li>")
		}
		sb.WriteString("</ul>")
	}

	if e, ok := e.(*model.Message); ok {
		sb.WriteString("<ul class='fields'>")

		fieldTypes := map[string]model.Entity{}
		fieldTypeList := []model.Entity{}

		for _, field := range e.GetFields() {
			sb.WriteString("<li>")
			sb.WriteString("<strong>")
			sb.WriteString(field.GetName() + ": ")
			sb.WriteString("</strong>")
			if field.GetEntity() != nil {
				if _, ok := fieldTypes[field.GetEntity().GetName()]; !ok {
					fieldTypes[field.GetEntity().GetName()] = field.GetEntity()
					fieldTypeList = append(fieldTypeList, field.GetEntity())
				}
				sb.WriteString("<a href='#" + makeLink(append(parents, e.GetName(), field.GetKind())) + "'>" + field.GetKind() + "</a>")
			} else {
				sb.WriteString(field.GetKind())
			}

			for _, label := range field.GetLabels() {
				sb.WriteString("<span class='label'>" + label + "</span>")
			}

			if field.IsDeprecated() {
				sb.WriteString("<span class='label'>deprecated</span>")
			}

			if field.GetDescription() != "" {
				sb.WriteString("<div>")
				sb.WriteString(field.GetDescription())
				sb.WriteString("</div>")
			}

			sb.WriteString("</li>")
		}
		sb.WriteString("</ul>")

		for _, fieldType := range fieldTypeList {
			sb.WriteString(renderEntity(f, fieldType, append(parents, e.GetName())).String())
		}
	}

	sb.WriteString("</details>")

	return sb
}

func makeLink(path []string) string {
	return strings.ToLower(strings.ReplaceAll(strings.Join(path, "-"), ".", ""))
}

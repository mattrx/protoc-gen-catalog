package html

import (
	"sort"
	"strings"
	"unicode"

	"github.com/mattrx/protoc-gen-library/model"
	"google.golang.org/protobuf/compiler/protogen"
)

const stylesheet = `

#toc {
  position: fixed;
  height: 98%;
  width: 400px;
  overflow-y: auto;
}
#toc ul {
  list-style: none;
  padding-left: 10px;
  padding-right: 10px;
}
li {
  font-family: "Verdana", Verdana, Arial, sans-serif;

  font-size: 12px;
  min-height: 20px;
  margin: 4px;
}
a {
  text-decoration: none;
}
#main {
  margin-left: 410px;
}
fieldset {
  margin-top: 5px;
  margin-bottom: 40px;
}
legend {
  font-weight: bold;
  font-size: larger;
  margin-bottom: 20px;
}
td.desc {
  font-style: italic;
  padding-bottom: 25px;
}
div.desc {
  font-style: italic;
  padding-bottom: 12px;
  margin-bottom: 12px;
  border-bottom: 1px solid black;
}
span.label {
  font-size: 13px;
  background: #e4b9c0;
  padding: 4px;
  border-radius: 6px;
  margin: 0 10px;
  text-transform: uppercase;
}
.deprecated {
  color: gray;
}
`

func Render(gen *protogen.Plugin, groups []*model.Group) error {

	sort.SliceStable(groups, func(i, j int) bool {
		iRunes := []rune(groups[i].GetName())
		jRunes := []rune(groups[j].GetName())

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

	f := gen.NewGeneratedFile("library.html", "")

	f.P("<html>")
	f.P("<head>")
	f.P("<style>" + stylesheet + "</style>")
	f.P("</head>")
	f.P("<body>")

	f.P("<nav id='toc'>")
	for _, group := range groups {
		entities := group.GetEntities()
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

		f.P("<p>" + group.GetName() + "</p>")
		f.P("<ul>")
		for _, entity := range entities {
			t := "<li><a href='#" + strings.ToLower(strings.ReplaceAll(entity.GetName(), ".", "")) + "'>" + strings.ReplaceAll(entity.GetName(), group.GetName()+".", "") + "</a></li>"
			if entity.IsDeprecated() {
				t = "<s>" + t + "</s>"
			}
			f.P(t)
		}
		f.P("</ul>")
	}
	f.P("</nav>")

	f.P("<section id='main'>")

	for _, group := range groups {
		entities := group.GetEntities()
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

		f.P("<h1>" + group.GetName() + "</h1>")
		for _, entity := range entities {
			if e, ok := entity.(*model.Enum); ok {
				renderEnum(f, e)
			}
			if e, ok := entity.(*model.Message); ok {
				renderMessage(f, e)
			}
		}
	}

	f.P("</section>")
	f.P("</body>")
	f.P("</html>")

	return nil

}

func renderEnum(f *protogen.GeneratedFile, e *model.Enum) {
	f.P("<fieldset id='" + makeLink(e.GetName()) + "'>")
	f.P("<legend>")
	f.P(e.GetName())
	if e.IsDeprecated() {
		f.P("<span class='label'>deprecated</span>")
	}
	f.P("</legend>")
	if e.GetDescription() != "" {
		f.P("<div class='desc'>" + e.GetDescription() + "</div>")
	}
	f.P("enum values:")
	f.P("<ul>")
	for _, value := range e.GetValues() {
		f.P("<li>" + value + "</li>")
	}
	f.P("</ul>")
	f.P("</fieldset>")
}

func renderMessage(f *protogen.GeneratedFile, e *model.Message) {
	f.P("<fieldset id='" + makeLink(e.GetName()) + "'>")
	f.P("<legend>")
	f.P(e.GetName())
	if e.IsDeprecated() {
		f.P("<span class='label'>DEPRECATED</span>")
	}
	f.P("</legend>")
	if e.GetDescription() != "" {
		f.P("<div class='desc'>" + e.GetDescription() + "</div>")
	}
	f.P("<table>")
	for _, field := range e.GetFields() {

		fieldType := field.GetKind()
		if strings.Contains(fieldType, ".") {
			fieldType = "<a href='#" + strings.ToLower(strings.ReplaceAll(field.GetKind(), ".", "")) + "'>" + field.GetKind() + "</a>"
		}

		labels := ""

		if field.IsDeprecated() {
			labels += "<span class='label'>deprecated</span>"
		}

		for _, label := range field.GetLabels() {
			labels += "<span class='label'>" + label + "</span>"
		}

		if field.IsDeprecated() {
			f.P("<tr class='deprecated'>")
		} else {
			f.P("<tr>")
		}
		f.P("<td><strong>" + field.GetName() + "</strong></td>")
		f.P("<td>" + fieldType + labels + "</td>")
		f.P("</tr>")
		if field.IsDeprecated() {
			f.P("<tr class='deprecated'>")
		} else {
			f.P("<tr>")
		}
		f.P("<td class='desc' colspan='2'>" + field.GetDescription() + "</td>")
		f.P("</tr>")
	}
	f.P("</table>")
	f.P("</fieldset>")
}

func makeLink(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, ".", ""))
}

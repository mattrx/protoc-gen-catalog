package model

type Field struct {
	deprecated  bool
	description string
	labels      []string
	name        string
	kind        string
}

func NewField(name string, description string) *Field {
	return &Field{
		name:        name,
		description: description,
	}
}

func (e Field) GetName() string {
	return e.name
}

func (e Field) GetDescription() string {
	return e.description
}

func (e *Field) SetDeprecated(v bool) {
	e.deprecated = v
}

func (e Field) IsDeprecated() bool {
	return e.deprecated
}

func (e *Field) SetKind(v string) {
	e.kind = v
}

func (e Field) GetKind() string {
	return e.kind
}

func (e *Field) AddLabel(v string) {
	e.labels = append(e.labels, v)
}

func (e Field) GetLabels() []string {
	return e.labels
}

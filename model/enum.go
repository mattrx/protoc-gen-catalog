package model

type Enum struct {
	deprecated  bool
	description string
	name        string
	values      []string
}

func NewEnum(name string, description string) *Enum {
	return &Enum{
		name:        name,
		description: description,
	}
}

func (e Enum) GetName() string {
	return e.name
}

func (e Enum) GetDescription() string {
	return e.description
}

func (e *Enum) SetDeprecated(v bool) {
	e.deprecated = v
}

func (e Enum) IsDeprecated() bool {
	return e.deprecated
}

func (e *Enum) AddValue(v string) {
	e.values = append(e.values, v)
}

func (e Enum) GetValues() []string {
	return e.values
}

package model

type Message struct {
	deprecated  bool
	description string
	fields      []*Field
	name        string
}

func NewMessage(name string, description string) *Message {
	return &Message{
		name:        name,
		description: description,
	}
}

func (e Message) GetName() string {
	return e.name
}

func (e Message) GetDescription() string {
	return e.description
}

func (e *Message) SetDeprecated(v bool) {
	e.deprecated = v
}

func (e Message) IsDeprecated() bool {
	return e.deprecated
}

func (e *Message) AddField(v *Field) {
	e.fields = append(e.fields, v)
}

func (e Message) GetFields() []*Field {
	return e.fields
}

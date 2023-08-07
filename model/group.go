package model

type Entity interface {
	GetName() string
	IsDeprecated() bool
}

type Group struct {
	name     string
	entities []Entity
}

func NewGroup(name string) *Group {
	return &Group{
		name: name,
	}
}

func (e Group) GetName() string {
	return e.name
}

func (e *Group) AddEntity(v Entity) {
	e.entities = append(e.entities, v)
}

func (e Group) GetEntities() []Entity {
	return e.entities
}

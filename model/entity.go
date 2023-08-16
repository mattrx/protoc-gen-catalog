package model

type Entity interface {
	AddUsage(*Field)
	GetDescription() string
	GetName() string
	IsDeprecated() bool
	IsUsed() bool
}

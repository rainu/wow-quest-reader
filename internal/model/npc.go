package model

type NonPlayerCharacter struct {
	Id   int64
	Name string
	Type string

	Male bool

	Locations []Coordinate
}

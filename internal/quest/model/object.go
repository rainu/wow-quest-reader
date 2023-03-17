package model

type Object struct {
	Id   int64
	Name string

	Locations []Coordinate
}

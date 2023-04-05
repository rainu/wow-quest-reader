package model

import "github.com/rainu/wow-quest-reader/internal/locale"

type Quest struct {
	Id       int64
	Obsolete bool

	Title       string
	Description string
	Progress    string
	Completion  string
	Locale      locale.Locale

	StartNPC    *NonPlayerCharacter
	StartObject *Object
	StartItem   *Item

	EndNPC    *NonPlayerCharacter
	EndObject *Object
	EndItem   *Item
}

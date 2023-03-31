package model

import (
	"github.com/rainu/wow-quest-client/internal/locale"
	"strings"
)

type QuestInformation struct {
	Id          string `json:"i"`
	Description string `json:"d"`
	Progress    string `json:"p"`
	Completion  string `json:"c"`
	L           string `json:"l"`
}

func (q QuestInformation) Locale() locale.Locale {
	if strings.HasPrefix(q.L, "de") {
		return locale.German
	}
	if strings.HasPrefix(q.L, "en") {
		return locale.English
	}

	return locale.Locale(q.L)
}

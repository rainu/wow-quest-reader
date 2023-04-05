package model

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/rainu/wow-quest-reader/internal/locale"
	"strings"
)

type Info struct {
	Quest struct {
		Id          int64  `json:"id"`
		Description string `json:"description"`
		Progress    string `json:"progress"`
		Completion  string `json:"completion"`
	} `json:"quest"`
	Gossip string `json:"gossip"`
	Npc    struct {
		Id   Guid   `json:"id"`
		Name string `json:"name"`
		Sex  Sex    `json:"Sex"`
	} `json:"npc"`
	Player struct {
		Name  string `json:"name"`
		Realm string `json:"realm"`
		Sex   Sex    `json:"Sex"`
		Race  string `json:"race"`
		Class string `json:"class"`
	} `json:"player"`
	Shortcut string `json:"shortcut"`
	L        string `json:"locale"`
}

func (i Info) GossipId() string {
	h := sha1.New()
	h.Write([]byte(i.Gossip))
	return hex.EncodeToString(h.Sum(nil))
}

func (i Info) Locale() locale.Locale {
	if strings.HasPrefix(i.L, "de") {
		return locale.German
	}
	if strings.HasPrefix(i.L, "en") {
		return locale.English
	}

	return locale.Locale(i.L)
}

package wowhead

import (
	"fmt"
	common "github.com/rainu/wow-quest-reader/internal/crawling/crawler"
	"github.com/rainu/wow-quest-reader/internal/locale"
	"regexp"
	"strconv"
)

var obsoleteTextMapping = map[locale.Locale]string{
	locale.English: "quest was marked obsolete",
	locale.German:  "quest wurde von blizzard als nicht genutzt markiert",
}
var notAvailableTextMapping = map[locale.Locale]string{
	locale.English: "quest is no longer available",
	locale.German:  "quest ist nicht mehr im spiel verfügbar",
}
var descriptionSectionMapping = map[locale.Locale]string{
	locale.English: "description",
	locale.German:  "beschreibung",
}
var completionSectionMapping = map[locale.Locale]string{
	locale.English: "completion",
	locale.German:  "vervollständigung",
}
var progressSectionMapping = map[locale.Locale]string{
	locale.English: "progress",
	locale.German:  "fortschritt",
}
var startSectionMapping = map[locale.Locale]string{
	locale.English: "start",
	locale.German:  "anfang",
}
var endSectionMapping = map[locale.Locale]string{
	locale.English: "end",
	locale.German:  "ende",
}

var gmapperDataRegex = regexp.MustCompile(`g_mapperData[ =]*({.*});`)

type gmapperDataZone struct {
	Coords [][]float32 `json:"coords"`
}
type gmapperData map[string][]gmapperDataZone

type crawler struct {
	locale locale.Locale
	client *client
}

func New(l locale.Locale) common.Crawler {
	switch l {
	case locale.English:
	case locale.German:
	default:
		panic(fmt.Sprintf("unsupported locale: %s", l))
	}

	return &crawler{
		locale: l,
		client: newClient(l),
	}
}

const (
	sectionStart = "__START__"
	sectionEnd   = "__END__"
)

func i64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return -1
	}
	return i
}

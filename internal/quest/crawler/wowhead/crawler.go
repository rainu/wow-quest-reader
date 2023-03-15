package wowhead

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/rainu/wow-quest-client/internal/locale"
	common "github.com/rainu/wow-quest-client/internal/quest/crawler"
	"github.com/rainu/wow-quest-client/internal/quest/model"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
	"time"
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

var suffixRegex = regexp.MustCompile(sectionStart + `.*$`)
var factsRegex = regexp.MustCompile(`printHtml\(.*`)
var npcIdRegex = regexp.MustCompile(`npc=([0-9]*)`)
var objectIdRegex = regexp.MustCompile(`object=([0-9]*)`)
var itemIdRegex = regexp.MustCompile(`item=([0-9]*)`)

type questSection struct {
	Index   int
	Title   string
	Content string
}

func (c *crawler) GetQuest(ctx context.Context, id int64) (*model.Quest, error) {
	doc, err := c.client.GetQuest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get content: %w", err)
	}
	startTime := time.Now()
	defer func() {
		logrus.
			WithField("duration", time.Now().Sub(startTime)).
			WithField("locale", c.locale).
			WithField("quest_id", id).
			Debug("Parse quest page.")
	}()

	if doc.Find(".database-detail-page-not-found-message").Length() > 0 {
		// quest not found
		return nil, nil
	}

	var sections []questSection
	doc.Find(`.text .heading-size-3`).Each(func(i int, selection *goquery.Selection) {
		sections = append(sections, questSection{
			Index: i,
			Title: selection.Text(),
		})
	})

	doc.Find(`.text .heading-size-3`).PrependHtml(sectionStart).AppendHtml(sectionEnd)
	doc.Find(".text br").PrependHtml("\n")

	rawContent := strings.ReplaceAll(doc.Find(`.text`).Text(), "\n", "\\n")

	split := strings.Split(rawContent, sectionEnd)
	for i := range sections {
		split[i+1] = suffixRegex.ReplaceAllString(split[i+1], "")
		split[i+1] = strings.ReplaceAll(split[i+1], "\\n", "\n")
		sections[i].Content = strings.TrimSpace(split[i+1])
	}

	boltText := strings.ToLower(doc.Find(".text b").Text())
	result := model.Quest{
		Id:       id,
		Obsolete: strings.Contains(boltText, obsoleteTextMapping[c.locale]) || strings.Contains(boltText, notAvailableTextMapping[c.locale]),
		Locale:   c.locale,
		Title:    doc.Find(".text .heading-size-1").First().Text(),
	}

	if i := c.findSection(sections, descriptionSectionMapping); i >= 0 {
		result.Description = sections[i].Content
	}
	if i := c.findSection(sections, completionSectionMapping); i >= 0 {
		result.Completion = sections[i].Content
	}
	if i := c.findSection(sections, progressSectionMapping); i >= 0 {
		result.Progress = sections[i].Content
	}

	listItems := strings.Split(factsRegex.FindString(doc.Find("script").Text()), "[li]")
	for _, item := range listItems {
		if strings.Contains(strings.ToLower(item), startSectionMapping[c.locale]) {
			if npcIdRegex.MatchString(item) {
				result.StartNPC = &model.NonPlayerCharacter{
					Id: i64(npcIdRegex.FindStringSubmatch(item)[1]),
				}
			} else if objectIdRegex.MatchString(item) {
				result.StartObject = &model.Object{
					Id: i64(objectIdRegex.FindStringSubmatch(item)[1]),
				}
			} else if itemIdRegex.MatchString(item) {
				result.StartItem = &model.Item{
					Id: i64(itemIdRegex.FindStringSubmatch(item)[1]),
				}
			}
		}
		if strings.Contains(strings.ToLower(item), endSectionMapping[c.locale]) {
			if npcIdRegex.MatchString(item) {
				result.EndNPC = &model.NonPlayerCharacter{
					Id: i64(npcIdRegex.FindStringSubmatch(item)[1]),
				}
			} else if objectIdRegex.MatchString(item) {
				result.EndObject = &model.Object{
					Id: i64(objectIdRegex.FindStringSubmatch(item)[1]),
				}
			} else if itemIdRegex.MatchString(item) {
				result.EndItem = &model.Item{
					Id: i64(itemIdRegex.FindStringSubmatch(item)[1]),
				}
			}
		}
	}

	return &result, err
}

func (c *crawler) findSection(sections []questSection, mapping map[locale.Locale]string) int {
	for i, section := range sections {
		if strings.ToLower(section.Title) == mapping[c.locale] {
			return i
		}
	}
	return -1
}

func i64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return -1
	}
	return i
}

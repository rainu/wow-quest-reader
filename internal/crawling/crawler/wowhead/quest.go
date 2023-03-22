package wowhead

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/rainu/wow-quest-client/internal/locale"
	"github.com/rainu/wow-quest-client/internal/model"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
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

	docPart := doc.Find(`.text`).First()
	docPart.Find("br").PrependHtml("\n")

	boltText := strings.ToLower(docPart.Find("b").Text())
	result := model.Quest{
		Id:       id,
		Obsolete: strings.Contains(boltText, obsoleteTextMapping[c.locale]) || strings.Contains(boltText, notAvailableTextMapping[c.locale]),
		Locale:   c.locale,
		Title:    docPart.Find(".heading-size-1").First().Text(),
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

	var sections []questSection
	docPart.Find(`.heading-size-3`).Each(func(i int, selection *goquery.Selection) {
		sections = append(sections, questSection{
			Index: i,
			Title: selection.Text(),
		})
	})

	// remove script sections - otherwise the content will be included in .Text() later
	docPart.Find("script").Remove()

	// remove the quest-check-script
	docPart.Find(`div>pre:contains("/run")`).Parent().Remove()

	// remove the start of other-section
	docPart.Find(`.heading-size-2.clear`).Remove()

	docPart.Find(`.heading-size-3`).PrependHtml(sectionStart).AppendHtml(sectionEnd)

	rawContent := strings.ReplaceAll(docPart.Text(), "\n", "\\n")
	rawContent = strings.ReplaceAll(rawContent, "\u00a0", "")
	rawContent = strings.ReplaceAll(rawContent, "  ", " ")

	split := strings.Split(rawContent, sectionEnd)
	for i := range sections {
		split[i+1] = suffixRegex.ReplaceAllString(split[i+1], "")
		split[i+1] = strings.ReplaceAll(split[i+1], "\\n", "\n")
		sections[i].Content = strings.TrimSpace(split[i+1])
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

package wowhead

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/rainu/wow-quest-client/internal/locale"
	"github.com/rainu/wow-quest-client/internal/model"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

var typeSectionMapping = map[locale.Locale]string{
	locale.English: "type",
	locale.German:  "art",
}
var typeRegex = regexp.MustCompile(`: ([^\[]*)`)

func (c *crawler) GetNpc(ctx context.Context, id int64) (*model.NonPlayerCharacter, error) {
	doc, err := c.client.GetNpc(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get content: %w", err)
	}
	startTime := time.Now()
	defer func() {
		logrus.
			WithField("duration", time.Now().Sub(startTime)).
			WithField("locale", c.locale).
			WithField("npc_id", id).
			Debug("Parse npc page.")
	}()

	result := model.NonPlayerCharacter{
		Id:   id,
		Name: doc.Find(".heading-size-1").First().Text(),
	}
	result.Name = strings.ReplaceAll(result.Name, "\u00a0", " ")

	var locationScript *goquery.Selection
	doc.Find("script").Each(func(i int, selection *goquery.Selection) {
		if strings.Contains(selection.Text(), "g_mapperData") {
			locationScript = selection
		}
	})

	if locationScript != nil && gmapperDataRegex.MatchString(locationScript.Text()) {
		var data gmapperData

		jData := gmapperDataRegex.FindStringSubmatch(locationScript.Text())[1]
		if err := json.Unmarshal([]byte(jData), &data); err == nil {
			for zoneId, locations := range data {
				for _, location := range locations {
					for _, coord := range location.Coords {
						result.Locations = append(result.Locations, model.Coordinate{
							Zone: model.Zone{
								Id: i64(zoneId),
							},
							X: coord[0],
							Y: coord[1],
						})
					}
				}
			}
		}
	}

	listItems := strings.Split(factsRegex.FindString(doc.Find("script").Text()), "[li]")
	for _, item := range listItems {
		if strings.Contains(strings.ToLower(item), typeSectionMapping[c.locale]) {
			if typeRegex.MatchString(item) {
				result.Type = typeRegex.FindStringSubmatch(item)[1]
			}
		}
	}

	return &result, nil
}

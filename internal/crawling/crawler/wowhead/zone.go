package wowhead

import (
	"context"
	"fmt"
	"github.com/rainu/wow-quest-client/internal/model"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func (c *crawler) GetZone(ctx context.Context, id int64) (*model.Zone, error) {
	doc, err := c.client.GetZone(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get content: %w", err)
	}
	startTime := time.Now()
	defer func() {
		logrus.
			WithField("duration", time.Now().Sub(startTime)).
			WithField("locale", c.locale).
			WithField("zone_id", id).
			Debug("Parse zone page.")
	}()

	result := model.Zone{
		Id:   id,
		Name: doc.Find(".heading-size-1").First().Text(),
	}
	result.Name = strings.ReplaceAll(result.Name, "\u00a0", " ")

	return &result, nil
}

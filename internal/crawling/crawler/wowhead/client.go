package wowhead

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/rainu/wow-quest-reader/internal/locale"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type client struct {
	baseUrl    string
	httpClient *http.Client
}

func newClient(l locale.Locale) *client {
	result := &client{
		baseUrl:    "https://www.wowhead.com",
		httpClient: &http.Client{},
	}

	switch l {
	case locale.English:
	case locale.German:
		result.baseUrl += "/de"
	default:
		panic(fmt.Sprintf("unsupported locale: %s", l))
	}

	return result
}

func (c *client) GetQuest(ctx context.Context, id int64) (*goquery.Document, error) {
	return c.get(ctx, fmt.Sprintf("%s/quest=%d", c.baseUrl, id))
}

func (c *client) GetNpc(ctx context.Context, id int64) (*goquery.Document, error) {
	return c.get(ctx, fmt.Sprintf("%s/npc=%d", c.baseUrl, id))
}

func (c *client) GetObject(ctx context.Context, id int64) (*goquery.Document, error) {
	return c.get(ctx, fmt.Sprintf("%s/object=%d", c.baseUrl, id))
}

func (c *client) GetItem(ctx context.Context, id int64) (*goquery.Document, error) {
	return c.get(ctx, fmt.Sprintf("%s/item=%d", c.baseUrl, id))
}

func (c *client) GetZone(ctx context.Context, id int64) (*goquery.Document, error) {
	return c.get(ctx, fmt.Sprintf("%s/zone=%d", c.baseUrl, id))
}

func (c *client) get(ctx context.Context, url string) (*goquery.Document, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to prepare request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/110.0")

	startTime := time.Now()
	resp, err := c.httpClient.Do(req)
	logrus.
		WithField("duration", time.Now().Sub(startTime)).
		WithField("req", fmt.Sprintf("%s - %s", req.Method, req.URL.String())).
		Debug("Do http call.")

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	return goquery.NewDocumentFromReader(resp.Body)
}

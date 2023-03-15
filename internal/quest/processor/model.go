package processor

import (
	"github.com/rainu/wow-quest-client/internal/quest/crawler"
	"github.com/rainu/wow-quest-client/internal/quest/model"
)

type job struct {
	Crawler    crawler.Crawler
	QuestId    int64
	ResultChan chan model.Quest
}

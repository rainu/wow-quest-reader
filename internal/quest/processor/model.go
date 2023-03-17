package processor

import (
	"github.com/rainu/wow-quest-client/internal/quest/crawler"
	"github.com/rainu/wow-quest-client/internal/quest/model"
)

type job struct {
	Crawler crawler.Crawler

	QuestId  *int64
	NpcId    *int64
	ItemId   *int64
	ObjectId *int64

	ResultQuestChan  chan model.Quest
	ResultNpcChan    chan model.NonPlayerCharacter
	ResultItemChan   chan model.Item
	ResultObjectChan chan model.Object
}

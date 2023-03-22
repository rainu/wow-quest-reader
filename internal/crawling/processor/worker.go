package processor

import (
	"context"
	"github.com/rainu/wow-quest-client/internal/crawling/crawler"
	"github.com/rainu/wow-quest-client/internal/model"
	"github.com/sirupsen/logrus"
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

type worker struct {
	jobChan chan job
}

func (w *worker) run(ctx context.Context) {
	logrus.Debug("Start worker.")
	defer func() {
		logrus.Debug("Stop worker.")
	}()

	for {
		var curJob job
		select {
		case <-ctx.Done():
			//context closed -> application is shutting down
			return
		case j, ok := <-w.jobChan:
			if !ok {
				return
			}
			curJob = j
		}

		if curJob.QuestId != nil {
			w.doQuestJob(ctx, curJob)
		} else if curJob.NpcId != nil {
			w.doNpcJob(ctx, curJob)
		} else if curJob.ObjectId != nil {
			w.doObjectJob(ctx, curJob)
		} else if curJob.ItemId != nil {
			w.doItemJob(ctx, curJob)
		}
	}
}

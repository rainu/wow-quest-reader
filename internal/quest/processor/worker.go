package processor

import (
	"context"
	"github.com/sirupsen/logrus"
)

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
			q, err := curJob.Crawler.GetQuest(ctx, *curJob.QuestId)
			if err != nil {
				logrus.WithField("quest_id", *curJob.QuestId).WithError(err).Error("Error while crawling quest.")
			}

			if q != nil {
				curJob.ResultQuestChan <- *q
			}
		}

		if curJob.NpcId != nil {
			npc, err := curJob.Crawler.GetNpc(ctx, *curJob.NpcId)
			if err != nil {
				logrus.WithField("npc_id", *curJob.NpcId).WithError(err).Error("Error while crawling npc.")
			}

			if npc != nil {
				curJob.ResultNpcChan <- *npc
			}
		}

		if curJob.ObjectId != nil {
			object, err := curJob.Crawler.GetObject(ctx, *curJob.ObjectId)
			if err != nil {
				logrus.WithField("object_id", *curJob.ObjectId).WithError(err).Error("Error while crawling object.")
			}

			if object != nil {
				curJob.ResultObjectChan <- *object
			}
		}

		if curJob.ItemId != nil {
			item, err := curJob.Crawler.GetItem(ctx, *curJob.ItemId)
			if err != nil {
				logrus.WithField("item_id", *curJob.ItemId).WithError(err).Error("Error while crawling item.")
			}

			if item != nil {
				curJob.ResultItemChan <- *item
			}
		}
	}
}

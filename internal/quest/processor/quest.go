package processor

import (
	"context"
	"github.com/rainu/wow-quest-client/internal/quest/model"
	"github.com/sirupsen/logrus"
)

func (p *processor) generateQuestJobs(ctx context.Context, jobChan chan job, questChan chan model.Quest) {
	knownIds, err := p.store.GetQuestIds(ctx)
	if err != nil {
		logrus.WithError(err).Error("Unable to get known quest ids!")
		return
	}
	idIter := newQuestIter(knownIds)

	for ctx.Err() == nil {
		nextId := idIter.Next()
		if nextId == -1 {
			//end reached
			break
		}

		// for each crawler (language)
		for _, c := range p.crawler {
			jobChan <- job{
				Crawler:         c,
				QuestId:         &nextId,
				ResultQuestChan: questChan,
			}
		}
	}
}

func (w *worker) doQuestJob(ctx context.Context, curJob job) {
	q, err := curJob.Crawler.GetQuest(ctx, *curJob.QuestId)
	if err != nil {
		logrus.WithField("quest_id", *curJob.QuestId).WithError(err).Error("Error while crawling quest.")
	}

	if q != nil {
		curJob.ResultQuestChan <- *q
	}
}

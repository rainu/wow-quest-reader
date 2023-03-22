package processor

import (
	"context"
	"github.com/rainu/wow-quest-client/internal/model"
	"github.com/sirupsen/logrus"
)

func (p *processor) generateItemJobs(ctx context.Context, jobChan chan job, itemChan chan model.Item) {
	ids, err := p.store.GetUnfinishedItemIds(ctx)
	if err != nil {
		logrus.WithError(err).Error("Unable to get item ids!")
		return
	}
	for i := 0; i < len(ids) && ctx.Err() == nil; i++ {
		jobChan <- job{
			Crawler:        p.nqCrawler,
			ItemId:         &ids[i],
			ResultItemChan: itemChan,
		}
	}
}

func (w *worker) doItemJob(ctx context.Context, curJob job) {
	item, err := curJob.Crawler.GetItem(ctx, *curJob.ItemId)
	if err != nil {
		logrus.WithField("item_id", *curJob.ItemId).WithError(err).Error("Error while crawling item.")
	}

	if item != nil {
		curJob.ResultItemChan <- *item
	}
}

package processor

import (
	"context"
	"github.com/rainu/wow-quest-reader/internal/model"
	"github.com/sirupsen/logrus"
)

func (p *processor) generateObjectJobs(ctx context.Context, jobChan chan job, objectChan chan model.Object) {
	ids, err := p.store.GetUnfinishedObjectIds(ctx)
	if err != nil {
		logrus.WithError(err).Error("Unable to get object ids!")
		return
	}
	for i := 0; i < len(ids) && ctx.Err() == nil; i++ {
		jobChan <- job{
			Crawler:          p.nqCrawler,
			ObjectId:         &ids[i],
			ResultObjectChan: objectChan,
		}
	}
}

func (w *worker) doObjectJob(ctx context.Context, curJob job) {
	object, err := curJob.Crawler.GetObject(ctx, *curJob.ObjectId)
	if err != nil {
		logrus.WithField("object_id", *curJob.ObjectId).WithError(err).Error("Error while crawling object.")
	}

	if object != nil {
		curJob.ResultObjectChan <- *object
	}
}

package processor

import (
	"context"
	"github.com/rainu/wow-quest-client/internal/quest/model"
	"github.com/sirupsen/logrus"
)

type worker struct {
	jobChan    chan job
	resultChan chan model.Quest
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

		q, err := curJob.Crawler.GetQuest(ctx, curJob.QuestId)
		if err != nil {
			logrus.WithError(err).Error("Error while crawling quest.")
		}

		if q != nil {
			w.resultChan <- *q
		}
	}
}

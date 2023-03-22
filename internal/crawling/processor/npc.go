package processor

import (
	"context"
	"github.com/rainu/wow-quest-client/internal/model"
	"github.com/sirupsen/logrus"
)

func (p *processor) generateNpcJobs(ctx context.Context, jobChan chan job, npcChan chan model.NonPlayerCharacter) {
	ids, err := p.store.GetUnfinishedNpcIds(ctx)
	if err != nil {
		logrus.WithError(err).Error("Unable to get npc ids!")
		return
	}
	for i := 0; i < len(ids) && ctx.Err() == nil; i++ {
		jobChan <- job{
			Crawler:       p.nqCrawler,
			NpcId:         &ids[i],
			ResultNpcChan: npcChan,
		}
	}
}

func (w *worker) doNpcJob(ctx context.Context, curJob job) {
	npc, err := curJob.Crawler.GetNpc(ctx, *curJob.NpcId)
	if err != nil {
		logrus.WithField("npc_id", *curJob.NpcId).WithError(err).Error("Error while crawling npc.")
	}

	if npc != nil {
		curJob.ResultNpcChan <- *npc
	}
}

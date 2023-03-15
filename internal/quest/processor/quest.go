package processor

import (
	"context"
	"github.com/rainu/wow-quest-client/internal/locale"
	"github.com/rainu/wow-quest-client/internal/quest/crawler"
	"github.com/rainu/wow-quest-client/internal/quest/crawler/wowhead"
	"github.com/rainu/wow-quest-client/internal/quest/model"
	"github.com/rainu/wow-quest-client/internal/quest/store"
	"github.com/sirupsen/logrus"
	"sync"
)

type processor struct {
	store   store.Store
	crawler map[locale.Locale]crawler.Crawler
}

func NewQuest(store store.Store, locales ...locale.Locale) *processor {
	result := &processor{
		store:   store,
		crawler: map[locale.Locale]crawler.Crawler{},
	}

	for _, l := range locales {
		result.crawler[l] = wowhead.New(l)
	}

	return result
}

func (p *processor) Run(ctx context.Context, workerCount int) {
	workers := make([]worker, workerCount)
	jobChan := make(chan job)
	resultChan := make(chan model.Quest)
	workerWg := sync.WaitGroup{}
	persistenceWg := sync.WaitGroup{}

	// initialise worker
	persistenceWg.Add(1)
	go func() {
		defer persistenceWg.Done()

		p.runPersistence(ctx, resultChan)
	}()

	for i := 0; i < workerCount; i++ {
		workers[i] = worker{
			jobChan:    jobChan,
			resultChan: resultChan,
		}

		workerWg.Add(1)
		go func(w worker) {
			defer workerWg.Done()

			w.run(ctx)
		}(workers[i])
	}

	// generate jobs...
	p.generateJobs(ctx, jobChan, resultChan)

	//wait until jobs are finished
	workerWg.Wait()

	//no one will write in resultChan anymore -> close them
	close(resultChan)

	//wait for persistence
	persistenceWg.Wait()
}

func (p *processor) generateJobs(ctx context.Context, jobChan chan job, resultChan chan model.Quest) {
	defer close(jobChan)

	knownIds, err := p.store.GetQuestIds(ctx)
	if err != nil {
		logrus.WithError(err).Error("Unable to get known quest ids!")
		return
	}
	idIter := newIter(knownIds)

	for ctx.Err() == nil {
		nextId := idIter.Next()
		if nextId == -1 {
			//end reached
			return
		}

		// for each crawler (language)
		for _, c := range p.crawler {
			jobChan <- job{
				Crawler:    c,
				QuestId:    nextId,
				ResultChan: resultChan,
			}
		}
	}
}

func (p *processor) runPersistence(ctx context.Context, resultChan chan model.Quest) {
	logrus.Debug("Start persistence worker.")
	defer func() {
		logrus.Debug("Stop persistence worker.")
	}()

	for {
		select {
		case <-ctx.Done():
			//context closed
			return
		case q, ok := <-resultChan:
			if !ok {
				return
			}

			log := logrus.WithField("quest_id", q.Id).WithField("locale", q.Locale)

			if vErr := q.IsValid(); vErr != nil {
				log.
					WithError(vErr).
					Warning("Quest is invalid.")
				continue
			}

			err := p.store.SaveQuest(ctx, q)
			if err != nil {
				log.WithError(err).Error("Error while persisting quest!")
			}
		}
	}
}

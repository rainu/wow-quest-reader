package processor

import (
	"context"
	"github.com/rainu/wow-quest-client/internal/locale"
	"github.com/rainu/wow-quest-client/internal/quest/crawler"
	"github.com/rainu/wow-quest-client/internal/quest/crawler/wowhead"
	"github.com/rainu/wow-quest-client/internal/quest/model"
	"github.com/rainu/wow-quest-client/internal/quest/store"
	"sync"
)

type processor struct {
	store   store.Store
	crawler map[locale.Locale]crawler.Crawler

	nqCrawler crawler.Crawler
}

func New(store store.Store, locales ...locale.Locale) *processor {
	result := &processor{
		store:     store,
		crawler:   map[locale.Locale]crawler.Crawler{},
		nqCrawler: wowhead.New(locale.English),
	}

	for _, l := range locales {
		result.crawler[l] = wowhead.New(l)
	}

	return result
}

func (p *processor) Run(ctx context.Context, workerCount int) {
	workers := make([]worker, workerCount)
	jobChan := make(chan job)
	questResultChan := make(chan model.Quest)
	npcResultChan := make(chan model.NonPlayerCharacter)
	itemResultChan := make(chan model.Item)
	objectResultChan := make(chan model.Object)
	workerWg := sync.WaitGroup{}
	persistenceWg := sync.WaitGroup{}

	// initialise worker
	persistenceWg.Add(1)
	go func() {
		defer persistenceWg.Done()

		p.runPersistence(ctx, questResultChan, npcResultChan, itemResultChan, objectResultChan)
	}()

	for i := 0; i < workerCount; i++ {
		workers[i] = worker{
			jobChan: jobChan,
		}

		workerWg.Add(1)
		go func(w worker) {
			defer workerWg.Done()

			w.run(ctx)
		}(workers[i])
	}

	// generate jobs...
	//p.generateQuestJobs(ctx, jobChan, questResultChan)
	p.generateNpcJobs(ctx, jobChan, npcResultChan)
	p.generateItemJobs(ctx, jobChan, itemResultChan)
	p.generateObjectJobs(ctx, jobChan, objectResultChan)
	close(jobChan)

	//wait until jobs are finished
	workerWg.Wait()

	//no one will write in *ResultChan anymore -> close them
	close(questResultChan)
	close(npcResultChan)
	close(itemResultChan)
	close(objectResultChan)

	//wait for persistence
	persistenceWg.Wait()
}

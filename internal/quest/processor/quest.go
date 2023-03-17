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

	nqCrawler crawler.Crawler
}

func NewQuest(store store.Store, locales ...locale.Locale) *processor {
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
	p.generateJobs(ctx, jobChan, questResultChan, npcResultChan, itemResultChan, objectResultChan)

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

func (p *processor) generateJobs(ctx context.Context, jobChan chan job,
	questChan chan model.Quest,
	npcChan chan model.NonPlayerCharacter,
	itemChan chan model.Item,
	objectChan chan model.Object,
) {
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

	knownIds, err = p.store.GetUnfinishedNpcIds(ctx)
	if err != nil {
		logrus.WithError(err).Error("Unable to get known npc ids!")
		return
	}
	for i := 0; i < len(knownIds) && ctx.Err() == nil; i++ {
		jobChan <- job{
			Crawler:       p.nqCrawler,
			NpcId:         &knownIds[i],
			ResultNpcChan: npcChan,
		}
	}

	knownIds, err = p.store.GetUnfinishedObjectIds(ctx)
	if err != nil {
		logrus.WithError(err).Error("Unable to get known object ids!")
		return
	}
	for i := 0; i < len(knownIds) && ctx.Err() == nil; i++ {
		jobChan <- job{
			Crawler:          p.nqCrawler,
			ObjectId:         &knownIds[i],
			ResultObjectChan: objectChan,
		}
	}

	knownIds, err = p.store.GetUnfinishedItemIds(ctx)
	if err != nil {
		logrus.WithError(err).Error("Unable to get known item ids!")
		return
	}
	for i := 0; i < len(knownIds) && ctx.Err() == nil; i++ {
		jobChan <- job{
			Crawler:        p.nqCrawler,
			ItemId:         &knownIds[i],
			ResultItemChan: itemChan,
		}
	}
}

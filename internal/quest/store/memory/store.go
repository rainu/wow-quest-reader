package memory

import (
	"context"
	"github.com/rainu/wow-quest-client/internal/quest/model"
	common "github.com/rainu/wow-quest-client/internal/quest/store"
	"sync"
)

type store struct {
	quests     map[int64]model.Quest
	questIds   []int64
	questMutex *sync.RWMutex
}

type iter struct {
	store    *store
	curIndex int
}

func New() common.Store {
	return &store{
		quests:     map[int64]model.Quest{},
		questIds:   []int64{},
		questMutex: &sync.RWMutex{},
	}
}

func (s *store) GetQuestIds(ctx context.Context) ([]int64, error) {
	s.questMutex.RLock()
	defer s.questMutex.RUnlock()

	return s.questIds, nil
}

func (s *store) SaveQuest(ctx context.Context, quest model.Quest) error {
	s.questMutex.Lock()
	defer s.questMutex.Unlock()

	s.questIds = append(s.questIds, quest.Id)
	s.quests[quest.Id] = quest

	return nil
}

func (s *store) Iterator() common.Iterator {
	return &iter{
		store:    s,
		curIndex: 0,
	}
}

func (i *iter) Next(ctx context.Context) *model.Quest {
	i.store.questMutex.RLock()
	defer i.store.questMutex.RUnlock()

	if ctx.Err() != nil {
		//context closed -> interrupt iteration
		return nil
	}

	if i.curIndex < len(i.store.questIds) {
		q := i.store.quests[i.store.questIds[i.curIndex]]
		i.curIndex++

		return &q
	}

	return nil
}

func (s *store) Close() error {
	return nil
}

package store

import (
	"context"
	"github.com/rainu/wow-quest-client/internal/quest/model"
)

type Store interface {
	GetQuestIds(ctx context.Context) ([]int64, error)
	SaveQuest(ctx context.Context, quest model.Quest) error
	Iterator() Iterator
	Close() error
}

type Iterator interface {
	Next(ctx context.Context) *model.Quest
}

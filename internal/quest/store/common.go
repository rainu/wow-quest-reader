package store

import (
	"context"
	"github.com/rainu/wow-quest-client/internal/locale"
	"github.com/rainu/wow-quest-client/internal/quest/model"
)

type Store interface {
	GetQuestIds(ctx context.Context) ([]int64, error)
	SaveQuest(ctx context.Context, quest model.Quest) error
	GetQuest(ctx context.Context, id int64, l locale.Locale) (*model.Quest, error)
	QuestIterator() Iterator

	GetUnfinishedNpcIds(ctx context.Context) ([]int64, error)
	SaveNpc(ctx context.Context, npc model.NonPlayerCharacter) error
	GetNpc(ctx context.Context, id int64) (*model.NonPlayerCharacter, error)

	GetUnfinishedObjectIds(ctx context.Context) ([]int64, error)
	SaveObject(ctx context.Context, object model.Object) error
	GetObject(ctx context.Context, id int64) (*model.Object, error)

	GetUnfinishedItemIds(ctx context.Context) ([]int64, error)
	SaveItem(ctx context.Context, item model.Item) error
	GetItem(ctx context.Context, id int64) (*model.Item, error)

	Close() error
}

type Iterator interface {
	Next(ctx context.Context) *model.Quest
}

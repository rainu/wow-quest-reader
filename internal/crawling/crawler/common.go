package crawler

import (
	"context"
	"github.com/rainu/wow-quest-client/internal/model"
)

type Crawler interface {
	GetQuest(ctx context.Context, id int64) (*model.Quest, error)
	GetNpc(ctx context.Context, id int64) (*model.NonPlayerCharacter, error)
	GetItem(ctx context.Context, id int64) (*model.Item, error)
	GetObject(ctx context.Context, id int64) (*model.Object, error)
	GetZone(ctx context.Context, id int64) (*model.Zone, error)
}

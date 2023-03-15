package crawler

import (
	"context"
	"github.com/rainu/wow-quest-client/internal/quest/model"
)

type Crawler interface {
	GetQuest(ctx context.Context, id int64) (*model.Quest, error)
}

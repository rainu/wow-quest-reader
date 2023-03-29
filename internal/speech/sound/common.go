package sound

import (
	"context"
	"github.com/rainu/wow-quest-client/internal/model"
	"io"
)

type SpeechGenerator interface {
	SpeechAsNpc(ctx context.Context, text string, npc model.NonPlayerCharacter) (io.ReadCloser, error)
	SpeechAsObject(ctx context.Context, text string, object model.Object) (io.ReadCloser, error)
	SpeechAsItem(ctx context.Context, text string, item model.Item) (io.ReadCloser, error)
}

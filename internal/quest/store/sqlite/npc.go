package sqlite

import (
	"context"
	"github.com/rainu/wow-quest-client/internal/quest/model"
	"github.com/sirupsen/logrus"
	"time"
)

func (s *store) GetUnfinishedNpcIds(ctx context.Context) ([]int64, error) {
	r, err := s.db.QueryContext(ctx, `SELECT `+fieldId+` FROM `+tableNpc+` WHERE `+fieldName+` IS NULL`)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var result []int64
	for r.Next() {
		var v int64
		if err := r.Scan(&v); err != nil {
			return nil, err
		}

		result = append(result, v)
	}

	return result, nil
}

func (s *store) SaveNpc(ctx context.Context, npc model.NonPlayerCharacter) error {
	startTime := time.Now()
	defer func() {
		logrus.
			WithField("duration", time.Now().Sub(startTime)).
			WithField("npc_id", npc.Id).
			Debug("Persist npc data.")
	}()

	_, err := s.insertNpcStmt.ExecContext(
		ctx,
		npc.Id,
		npc.Name,
		npc.Type,
	)

	return err
}

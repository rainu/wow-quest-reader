package sqlite

import (
	"context"
	"github.com/rainu/wow-quest-reader/internal/model"
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

func (s *store) GetNpc(ctx context.Context, id int64) (*model.NonPlayerCharacter, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT `+fieldId+`, 
		`+fieldName+`,
		`+fieldType+
		` FROM `+tableNpc+` WHERE `+fieldId+` = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	npc := model.NonPlayerCharacter{}

	err = rows.Scan(&npc.Id, &npc.Name, &npc.Type)
	if err != nil {
		logrus.WithError(err).Error("Unable to scan npc!")
		return nil, err
	}

	return &npc, nil
}

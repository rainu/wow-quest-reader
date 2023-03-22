package sqlite

import (
	"context"
	"github.com/rainu/wow-quest-client/internal/model"
	"github.com/sirupsen/logrus"
	"time"
)

func (s *store) GetUnfinishedItemIds(ctx context.Context) ([]int64, error) {
	r, err := s.db.QueryContext(ctx, `SELECT `+fieldId+` FROM `+tableItem+` WHERE `+fieldName+` IS NULL`)
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

func (s *store) SaveItem(ctx context.Context, item model.Item) error {
	startTime := time.Now()
	defer func() {
		logrus.
			WithField("duration", time.Now().Sub(startTime)).
			WithField("item_id", item.Id).
			Debug("Persist item data.")
	}()

	_, err := s.insertItemStmt.ExecContext(
		ctx,
		item.Id,
		item.Name,
	)

	return err
}

func (s *store) GetItem(ctx context.Context, id int64) (*model.Item, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT `+fieldId+`, 
		`+fieldName+
		` FROM `+tableItem+` WHERE `+fieldId+` = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	i := model.Item{}

	err = rows.Scan(&i.Id, &i.Name)
	if err != nil {
		logrus.WithError(err).Error("Unable to scan item!")
		return nil, err
	}

	return &i, nil
}

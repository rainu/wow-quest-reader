package sqlite

import (
	"context"
	"github.com/rainu/wow-quest-client/internal/model"
	"github.com/sirupsen/logrus"
	"time"
)

func (s *store) GetUnfinishedObjectIds(ctx context.Context) ([]int64, error) {
	r, err := s.db.QueryContext(ctx, `SELECT `+fieldId+` FROM `+tableObject+` WHERE `+fieldName+` IS NULL`)
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

func (s *store) SaveObject(ctx context.Context, object model.Object) error {
	startTime := time.Now()
	defer func() {
		logrus.
			WithField("duration", time.Now().Sub(startTime)).
			WithField("object_id", object.Id).
			Debug("Persist object data.")
	}()

	_, err := s.insertObjectStmt.ExecContext(
		ctx,
		object.Id,
		object.Name,
	)

	return err
}

func (s *store) GetObject(ctx context.Context, id int64) (*model.Object, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT `+fieldId+`, 
		`+fieldName+
		` FROM `+tableObject+` WHERE `+fieldId+` = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	o := model.Object{}

	err = rows.Scan(&o.Id, &o.Name)
	if err != nil {
		logrus.WithError(err).Error("Unable to scan object!")
		return nil, err
	}

	return &o, nil
}

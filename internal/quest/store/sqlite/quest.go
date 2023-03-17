package sqlite

import (
	"context"
	"database/sql"
	"github.com/rainu/wow-quest-client/internal/locale"
	"github.com/rainu/wow-quest-client/internal/quest/model"
	common "github.com/rainu/wow-quest-client/internal/quest/store"
	"github.com/sirupsen/logrus"
	"time"
)

func (s *store) GetQuestIds(ctx context.Context) ([]int64, error) {
	r, err := s.db.QueryContext(ctx, `SELECT DISTINCT `+fieldId+` FROM `+tableQuest)
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

func (s *store) SaveQuest(ctx context.Context, quest model.Quest) error {
	startTime := time.Now()
	defer func() {
		logrus.
			WithField("duration", time.Now().Sub(startTime)).
			WithField("locale", quest.Locale).
			WithField("quest_id", quest.Id).
			Debug("Persist quest data.")
	}()

	if quest.StartNPC != nil || quest.EndNPC != nil {
		if quest.StartNPC != nil {
			_, err := s.insertNpcStmt.ExecContext(ctx, quest.StartNPC.Id, nil, nil)
			if err != nil {
				return err
			}
		}
		if quest.EndNPC != nil {
			_, err := s.insertNpcStmt.ExecContext(ctx, quest.EndNPC.Id, nil, nil)
			if err != nil {
				return err
			}
		}
	}
	if quest.StartObject != nil || quest.EndObject != nil {
		if quest.StartObject != nil {
			_, err := s.insertObjectStmt.ExecContext(ctx, quest.StartObject.Id, nil)
			if err != nil {
				return err
			}
		}
		if quest.EndObject != nil {
			_, err := s.insertObjectStmt.ExecContext(ctx, quest.EndObject.Id, nil)
			if err != nil {
				return err
			}
		}
	}
	if quest.StartItem != nil || quest.EndItem != nil {
		if quest.StartItem != nil {
			_, err := s.insertItemStmt.ExecContext(ctx, quest.StartItem.Id, nil)
			if err != nil {
				return err
			}
		}
		if quest.EndItem != nil {
			_, err := s.insertItemStmt.ExecContext(ctx, quest.EndItem.Id, nil)
			if err != nil {
				return err
			}
		}
	}

	_, err := s.insertQuestStmt.ExecContext(
		ctx,
		quest.Id,
		quest.Locale,
		quest.Obsolete,
		quest.Title,
		quest.Description,
		quest.Progress,
		quest.Completion,
		npcIdOrNil(quest.StartNPC),
		objectIdOrNil(quest.StartObject),
		itemIdOrNil(quest.StartItem),
		npcIdOrNil(quest.EndNPC),
		objectIdOrNil(quest.EndObject),
		itemIdOrNil(quest.EndItem),
	)

	return err
}

func (s *store) GetQuest(ctx context.Context, id int64, l locale.Locale) (*model.Quest, error) {
	rows, err := s.db.Query(`SELECT `+fieldId+`, 
		`+fieldLocale+`, 
		`+fieldObsolete+`, 
		`+fieldTitle+`, 
		`+fieldDescription+`, 
		`+fieldProgress+`, 
		`+fieldCompletion+`, 
		`+fieldStartNpc+`, 
		`+fieldStartObject+`, 
		`+fieldStartItem+`, 
		`+fieldEndNpc+`, 
		`+fieldEndObject+`, 
		`+fieldEndItem+
		` FROM `+tableQuest+` WHERE `+fieldId+` = ? AND `+fieldLocale+` = ?`, id, l)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.extractQuest(ctx, rows), nil
}

func (s *store) extractQuest(ctx context.Context, rows *sql.Rows) *model.Quest {
	if rows == nil {
		return nil
	}

	if ctx.Err() != nil {
		rows.Close()
		return nil
	}

	if !rows.Next() {
		rows.Close()
		return nil
	}

	q := model.Quest{}

	var sNpc, sObject, sItem, eNpc, eObject, eItem *int64

	err := rows.Scan(&q.Id, &q.Locale, &q.Obsolete,
		&q.Title, &q.Description, &q.Progress, &q.Completion,
		&sNpc, &sObject, &sItem,
		&eNpc, &eObject, &eItem,
	)
	if err != nil {
		rows.Close()
		logrus.WithError(err).Error("Unable to scan quest!")
		return nil
	}

	if sNpc != nil {
		q.StartNPC, err = s.GetNpc(ctx, *sNpc)
		if err != nil {
			logrus.WithError(err).Error("Unable to get npc!")
		}
	}
	if sObject != nil {
		q.StartObject, err = s.GetObject(ctx, *sObject)
		if err != nil {
			logrus.WithError(err).Error("Unable to get object!")
		}
	}
	if sItem != nil {
		q.StartItem, err = s.GetItem(ctx, *sItem)
		if err != nil {
			logrus.WithError(err).Error("Unable to get item!")
		}
	}
	if eNpc != nil {
		q.EndNPC, err = s.GetNpc(ctx, *eNpc)
		if err != nil {
			logrus.WithError(err).Error("Unable to get npc!")
		}
	}
	if eObject != nil {
		q.EndObject, err = s.GetObject(ctx, *eObject)
		if err != nil {
			logrus.WithError(err).Error("Unable to get object!")
		}
	}
	if eItem != nil {
		q.EndItem, err = s.GetItem(ctx, *eItem)
		if err != nil {
			logrus.WithError(err).Error("Unable to get item!")
		}
	}

	return &q
}

func (s *store) QuestIterator() common.Iterator {
	rows, err := s.db.Query(`SELECT ` + fieldId + `, 
		` + fieldLocale + `, 
		` + fieldObsolete + `, 
		` + fieldTitle + `, 
		` + fieldDescription + `, 
		` + fieldProgress + `, 
		` + fieldCompletion + `, 
		` + fieldStartNpc + `, 
		` + fieldStartObject + `, 
		` + fieldStartItem + `, 
		` + fieldEndNpc + `, 
		` + fieldEndObject + `, 
		` + fieldEndItem +
		` FROM ` + tableQuest)

	if err != nil {
		logrus.WithError(err).Error("Unable to select all quests from database!")
		return &iter{}
	}

	return &iter{
		rows:  rows,
		store: s,
	}
}

func (i *iter) Next(ctx context.Context) *model.Quest {
	return i.store.extractQuest(ctx, i.rows)
}

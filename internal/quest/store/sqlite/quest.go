package sqlite

import (
	"context"
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
		rows: rows,
	}
}

func (i *iter) Next(ctx context.Context) *model.Quest {
	if i.rows == nil {
		return nil
	}

	if ctx.Err() != nil {
		i.rows.Close()
		return nil
	}

	if !i.rows.Next() {
		i.rows.Close()
		return nil
	}

	q := model.Quest{}

	var sNpc, sObject, sItem, eNpc, eObject, eItem *int64

	err := i.rows.Scan(&q.Id, &q.Locale, &q.Obsolete,
		&q.Title, &q.Description, &q.Progress, &q.Completion,
		&sNpc, &sObject, &sItem,
		&eNpc, &eObject, &eItem,
	)
	if err != nil {
		i.rows.Close()
		logrus.WithError(err).Error("Unable to scan quest!")
		return nil
	}

	if sNpc != nil {
		q.StartNPC = &model.NonPlayerCharacter{Id: *sNpc}
	}
	if sObject != nil {
		q.StartObject = &model.Object{Id: *sObject}
	}
	if sItem != nil {
		q.StartItem = &model.Item{Id: *sItem}
	}
	if eNpc != nil {
		q.EndNPC = &model.NonPlayerCharacter{Id: *eNpc}
	}
	if eObject != nil {
		q.EndObject = &model.Object{Id: *eObject}
	}
	if eItem != nil {
		q.EndItem = &model.Item{Id: *eItem}
	}

	return &q
}

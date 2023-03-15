package sqlite

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rainu/wow-quest-client/internal/quest/model"
	common "github.com/rainu/wow-quest-client/internal/quest/store"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	tableQuest  = "quest"
	tableNpc    = "npc"
	tableObject = "object"
	tableItem   = "item"

	fieldId          = "id"
	fieldObsolete    = "obsolete"
	fieldTitle       = "title"
	fieldDescription = "description"
	fieldProgress    = "progress"
	fieldCompletion  = "completion"
	fieldLocale      = "locale"
	fieldStartNpc    = "start_npc"
	fieldStartObject = "start_object"
	fieldStartItem   = "start_item"
	fieldEndNpc      = "end_npc"
	fieldEndObject   = "end_object"
	fieldEndItem     = "end_item"

	fieldName = "name"
	fieldType = "type"
)

type store struct {
	db *sql.DB

	insertQuestStmt  *sql.Stmt
	insertNpcStmt    *sql.Stmt
	insertObjectStmt *sql.Stmt
	insertItemStmt   *sql.Stmt
}

type iter struct {
	rows *sql.Rows
}

func New(path string) (common.Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ` + tableNpc + `( 
		` + fieldId + ` INTEGER NOT NULL PRIMARY KEY, 
		` + fieldName + ` TEXT, 
		` + fieldType + ` TEXT
	);`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ` + tableObject + `( 
		` + fieldId + ` INTEGER NOT NULL PRIMARY KEY, 
		` + fieldName + ` TEXT 
	);`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ` + tableItem + `( 
		` + fieldId + ` INTEGER NOT NULL PRIMARY KEY, 
		` + fieldName + ` TEXT 
	);`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ` + tableQuest + `( 
		` + fieldId + ` INTEGER NOT NULL, 
		` + fieldLocale + ` TEXT NOT NULL, 
		` + fieldObsolete + ` INTEGER NOT NULL, 
		` + fieldTitle + ` TEXT, 
		` + fieldDescription + ` TEXT, 
		` + fieldProgress + ` TEXT, 
		` + fieldCompletion + ` TEXT, 
		` + fieldStartNpc + ` INTEGER, 
		` + fieldStartObject + ` INTEGER, 
		` + fieldStartItem + ` INTEGER, 
		` + fieldEndNpc + ` INTEGER, 
		` + fieldEndObject + ` INTEGER, 
		` + fieldEndItem + ` INTEGER,

		FOREIGN KEY(` + fieldStartNpc + `) REFERENCES ` + tableNpc + `(` + fieldId + `),
		FOREIGN KEY(` + fieldEndNpc + `) REFERENCES ` + tableNpc + `(` + fieldId + `),
		FOREIGN KEY(` + fieldStartObject + `) REFERENCES ` + tableObject + `(` + fieldId + `),
		FOREIGN KEY(` + fieldEndObject + `) REFERENCES ` + tableObject + `(` + fieldId + `),
		FOREIGN KEY(` + fieldStartItem + `) REFERENCES ` + tableItem + `(` + fieldId + `),
		FOREIGN KEY(` + fieldEndItem + `) REFERENCES ` + tableItem + `(` + fieldId + `),

		PRIMARY KEY ( ` + fieldId + `,` + fieldLocale + ` )
	);`)
	if err != nil {
		return nil, err
	}

	result := &store{
		db: db,
	}

	result.insertQuestStmt, err = db.Prepare(`INSERT INTO ` + tableQuest + `(
		` + fieldId + `, 
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
		` + fieldEndItem + `
	) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		return nil, err
	}
	result.insertNpcStmt, err = db.Prepare(`INSERT OR IGNORE INTO ` + tableNpc + `(
		` + fieldId + `, 
		` + fieldName + `, 
		` + fieldType + `
	) VALUES(?,?,?)`)
	if err != nil {
		return nil, err
	}
	result.insertObjectStmt, err = db.Prepare(`INSERT OR IGNORE INTO ` + tableObject + `(
		` + fieldId + `, 
		` + fieldName + `
	) VALUES(?,?)`)
	if err != nil {
		return nil, err
	}
	result.insertItemStmt, err = db.Prepare(`INSERT OR IGNORE INTO ` + tableItem + `(
		` + fieldId + `, 
		` + fieldName + `
	) VALUES(?,?)`)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *store) GetQuestIds(ctx context.Context) ([]int64, error) {
	r, err := s.db.Query(`SELECT DISTINCT ` + fieldId + ` FROM ` + tableQuest)
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

func npcIdOrNil(npc *model.NonPlayerCharacter) *int64 {
	if npc == nil {
		return nil
	}
	return &npc.Id
}

func itemIdOrNil(item *model.Item) *int64 {
	if item == nil {
		return nil
	}
	return &item.Id
}

func objectIdOrNil(object *model.Object) *int64 {
	if object == nil {
		return nil
	}
	return &object.Id
}

func (s *store) Iterator() common.Iterator {
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

func (s *store) Close() error {
	return s.db.Close()
}

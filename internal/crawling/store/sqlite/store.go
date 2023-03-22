package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	common "github.com/rainu/wow-quest-client/internal/crawling/store"
	"github.com/rainu/wow-quest-client/internal/model"
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
	rows  *sql.Rows
	store *store
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
	result.insertNpcStmt, err = db.Prepare(`INSERT OR REPLACE INTO ` + tableNpc + `(
		` + fieldId + `, 
		` + fieldName + `, 
		` + fieldType + `
	) VALUES(?,?,?)`)
	if err != nil {
		return nil, err
	}
	result.insertObjectStmt, err = db.Prepare(`INSERT OR REPLACE INTO ` + tableObject + `(
		` + fieldId + `, 
		` + fieldName + `
	) VALUES(?,?)`)
	if err != nil {
		return nil, err
	}
	result.insertItemStmt, err = db.Prepare(`INSERT OR REPLACE INTO ` + tableItem + `(
		` + fieldId + `, 
		` + fieldName + `
	) VALUES(?,?)`)
	if err != nil {
		return nil, err
	}

	return result, nil
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

func (s *store) Close() error {
	return s.db.Close()
}

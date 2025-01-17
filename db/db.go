package db

import (
	"database/sql"
)

type Database struct {
	*sql.DB
}

func (db *Database) Bootstrap() (err error) {
	stmt := `CREATE TABLE IF NOT EXISTS  keys (key VARCHAR NOT NULL PRIMARY KEY, value TEXT NOT NULL);`

	_, err = db.Exec(stmt)

	return
}

func (db *Database) GetValue(key string) (value string, err error) {
	stmt := `SELECT value FROM keys WHERE key=$1 LIMIT 1;`

	err = db.QueryRow(stmt, key).Scan(&value)

	return
}

func (db *Database) SetKey(key, value string) (storedValue string, err error) {
	stmt := `INSERT INTO keys(key,value) values($1,$2) RETURNING value;`

	err = db.QueryRow(stmt, key, value).Scan(&storedValue)

	return
}

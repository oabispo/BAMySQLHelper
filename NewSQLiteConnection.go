package bamysqlhelper

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func NewSQLiteConnection(fileName string) *sql.DB {
	db, err := sql.Open("sqlite3", fileName)

	if err != nil {
		panic(err.Error())
	}

	return db
}

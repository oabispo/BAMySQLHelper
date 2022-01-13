package bamysqlhelper

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLConnection(hostname string, database string, username string, password string) *sql.DB {
	var connectionString = username + ":" + password + "@tcp(" + hostname + ")/" + database + "?parseTime=true"
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		panic(err.Error())
	}

	return db
}

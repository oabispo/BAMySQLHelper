package bamysqlhelper

import (
	"database/sql"

	"github.com/oabispo/bamysqlhelper/internal/bamyhelper"
)

//"github.com/oabispo/bamysqlhelper/internal/bamyhelper"

type BASQL struct {
	db *sql.DB
}

type BAHelper interface {
	MapFields(fetch *sql.Rows) error
}

func NewBASQL(db *sql.DB) *BASQL {
	return &BASQL{db: db}
}

func (bas *BASQL) GetDB() *sql.DB {
	return bas.db
}

func (bas *BASQL) FetchOne(newCallback func() interface{}, stmt string, params ...interface{}) (interface{}, error) {
	return bamyhelper.FetchOne(bas.db, newCallback, stmt, params...)
}

func (bas *BASQL) FetchMany(newCallback func() interface{}, stmt string, params ...interface{}) ([]interface{}, error) {
	return bamyhelper.FetchMany(bas.db, newCallback, stmt, params...)
}

func (bas *BASQL) Insert(stmt string, params ...interface{}) (int64, error) {
	return bamyhelper.Insert(bas.db, stmt, params...)
}

func (bas *BASQL) Update(stmt string, params ...interface{}) (int64, error) {
	return bamyhelper.Update(bas.db, stmt, params...)
}

func (bas *BASQL) Delete(stmt string, params ...interface{}) (int64, error) {
	return bamyhelper.Delete(bas.db, stmt, params...)
}

func (bas *BASQL) RunSQL(stmt string, params ...interface{}) (interface{}, error) {
	return bamyhelper.RunSQL(bas.db, stmt, params...)
}

func (bas *BASQL) FetchPage(rowsPerPage int, currentPage int, newCallback func() interface{}, stmt string, params ...interface{}) ([]interface{}, error) {
	return bamyhelper.FetchPage(bas.db, rowsPerPage, currentPage, newCallback, stmt, params...)
}

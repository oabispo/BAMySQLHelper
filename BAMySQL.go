package bamysqlhelper

import (
	"database/sql"
	"strconv"
	"strings"
)

type BAMySQL struct {
	db *sql.DB
}

type BAHelper interface {
	MapFields(fetch *sql.Rows) error
}

func NewBAHelper(db *sql.DB) *BAMySQL {
	return &BAMySQL{db: db}
}

func (dbm *BAMySQL) fetchOne(fetch *sql.Rows, newCallback func() interface{}) (interface{}, error) {
	data := newCallback()
	var helper BAHelper = data.(BAHelper)
	err := helper.MapFields(fetch)
	return data, err
}

func (dbm *BAMySQL) FetchOne(newCallback func() interface{}, stmt string, params ...interface{}) (interface{}, error) {
	if newCallback != nil {
		fetch, err := dbm.makeFetch(stmt, params)

		if err != nil {
			return nil, err
		}

		if fetch.Next() {
			data, err := dbm.fetchOne(fetch, newCallback)
			fetch.Close()

			if err != nil {
				return nil, err
			} else {
				return data, nil
			}
		} else {
			fetch.Close()
			return nil, nil
		}
	} else {
		panic("callback não definida!")
	}
}

func (dbm *BAMySQL) makeFetch(stmt string, params ...interface{}) (*sql.Rows, error) {
	var fetch *sql.Rows
	var err error

	if params != nil {
		fetch, err = dbm.db.Query(stmt, params...)
	} else {
		fetch, err = dbm.db.Query(stmt)
	}

	return fetch, err
}

func (dbm *BAMySQL) makeSQLResult(stmt string, params ...interface{}) (sql.Result, error) {
	var err error
	var result sql.Result

	if params != nil {
		result, err = dbm.db.Exec(stmt, params...)
	} else {
		result, err = dbm.db.Exec(stmt)
	}

	return result, err
}

func (dbm *BAMySQL) FetchMany(newCallback func() interface{}, stmt string, params ...interface{}) ([]interface{}, error) {
	if newCallback != nil {
		fetch, err := dbm.makeFetch(stmt, params)

		if err != nil {
			return nil, err
		} else {
			var result []interface{} = make([]interface{}, 0, 10)
			for fetch.Next() {
				data, err := dbm.fetchOne(fetch, newCallback)
				if err != nil {
					panic(err.Error)
				} else {
					result = append(result, data)
				}
			}
			fetch.Close()

			return result, nil
		}
	} else {
		panic("callback não definida!")
	}
}

func (dbm *BAMySQL) FetchPage(rowsPerPage int, currentPage int, newCallback func() interface{}, stmt string, params ...interface{}) ([]interface{}, error) {
	// Paginador super simples. O suficiente para que eu possa deixar o retorno de dados um pouco melhor. O sonho seria construir um interpretador de SQL.
	var sb strings.Builder
	sb.WriteString(stmt)
	sb.WriteString(" limit ")
	sb.WriteString(strconv.Itoa(rowsPerPage * (currentPage - 1)))
	sb.WriteString(", ")
	sb.WriteString(strconv.Itoa(rowsPerPage))

	var newStmt = sb.String()
	data, err := dbm.FetchMany(newCallback, newStmt, params...)
	return data, err
}

func (dbm *BAMySQL) Insert(stmt string, params ...interface{}) (interface{}, error) {
	result, err := dbm.makeSQLResult(stmt, params)

	if err == nil {
		var id int64
		if id, err = result.LastInsertId(); err == nil {
			return id, err
		}
		return nil, err
	} else {
		return nil, err
	}
}

func (dbm *BAMySQL) Update(stmt string, params ...interface{}) (int64, error) {
	result, err := dbm.makeSQLResult(stmt, params)

	if err == nil {
		var total int64 = 0
		if total, err = result.RowsAffected(); err == nil {
			return total, err
		}
		return -1, err
	} else {
		return -1, err
	}
}

func (dbm *BAMySQL) Delete(stmt string, params ...interface{}) (int64, error) {
	result, err := dbm.makeSQLResult(stmt, params)

	if err == nil {
		var total int64 = 0
		if total, err = result.RowsAffected(); err == nil {
			return total, err
		}
		return -1, err
	} else {
		return -1, err
	}
}

func (dbm *BAMySQL) RunSQL(stmt string, params ...interface{}) (interface{}, error) {
	processResult := func(result interface{}, err error) (interface{}, error) {
		if err == nil {
			switch result.(type) {
			case *sql.Rows:
				{
					return result, err
				}
			default:
				{
					if total, e := result.(sql.Result).RowsAffected(); err == nil {
						return total, e
					}
				}
			}
		}
		return nil, err
	}

	if params != nil {
		result, err := dbm.db.Exec(stmt, params...)
		return processResult(result, err)
	} else {
		result, err := dbm.db.Exec(stmt)
		return processResult(result, err)
	}
}

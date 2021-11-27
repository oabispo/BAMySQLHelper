package bamyhelper

import (
	"database/sql"
	"strconv"
	"strings"
)

type BASQLMapper interface {
	MapFields(fetch *sql.Rows) error
}

func makeFetch(db *sql.DB, stmt string, params ...interface{}) (*sql.Rows, error) {
	var fetch *sql.Rows
	var err error

	if params != nil {
		fetch, err = db.Query(stmt, params...)
	} else {
		fetch, err = db.Query(stmt)
	}

	return fetch, err
}

func makeSQLResult(db *sql.DB, stmt string, params ...interface{}) (sql.Result, error) {
	var err error
	var result sql.Result

	if params != nil {
		result, err = db.Exec(stmt, params...)
	} else {
		result, err = db.Exec(stmt)
	}

	return result, err
}

func fetchOne(db *sql.DB, fetch *sql.Rows, newCallback func() interface{}) (interface{}, error) {
	data := newCallback()
	var helper BASQLMapper = data.(BASQLMapper)
	err := helper.MapFields(fetch)
	return data, err
}

func FetchOne(db *sql.DB, newCallback func() interface{}, stmt string, params ...interface{}) (interface{}, error) {
	if newCallback != nil {
		fetch, err := makeFetch(db, stmt, params...)

		if err != nil {
			return nil, err
		}

		if fetch.Next() {
			data, err := fetchOne(db, fetch, newCallback)
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

func FetchMany(db *sql.DB, newCallback func() interface{}, stmt string, params ...interface{}) ([]interface{}, error) {
	if newCallback != nil {
		fetch, err := makeFetch(db, stmt, params...)

		if err != nil {
			return nil, err
		} else {
			var result []interface{} = make([]interface{}, 0, 10)
			for fetch.Next() {
				data, err := fetchOne(db, fetch, newCallback)
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

func Insert(db *sql.DB, stmt string, params ...interface{}) (int64, error) {
	result, err := makeSQLResult(db, stmt, params...)

	if err == nil {
		if id, er := result.LastInsertId(); er == nil {
			return id, er
		} else {
			return 0, err
		}
	} else {
		return 0, err
	}
}

func Update(db *sql.DB, stmt string, params ...interface{}) (int64, error) {
	result, err := makeSQLResult(db, stmt, params...)

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

func Delete(db *sql.DB, stmt string, params ...interface{}) (int64, error) {
	result, err := makeSQLResult(db, stmt, params...)

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

func RunSQL(db *sql.DB, stmt string, params ...interface{}) (interface{}, error) {
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
		result, err := db.Exec(stmt, params...)
		return processResult(result, err)
	} else {
		result, err := db.Exec(stmt)
		return processResult(result, err)
	}
}

func FetchPage(db *sql.DB, rowsPerPage int, currentPage int, newCallback func() interface{}, stmt string, params ...interface{}) ([]interface{}, error) {
	// Paginador super simples. O suficiente para que eu possa deixar o retorno de dados um pouco melhor. O sonho seria construir um interpretador de SQL.
	var sb strings.Builder
	sb.WriteString(stmt)
	sb.WriteString(" limit ")
	sb.WriteString(strconv.Itoa(rowsPerPage * (currentPage - 1)))
	sb.WriteString(", ")
	sb.WriteString(strconv.Itoa(rowsPerPage))

	var newStmt = sb.String()
	data, err := FetchMany(db, newCallback, newStmt, params...)
	return data, err
}

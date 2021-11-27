package bamysqlhelper

import "database/sql"

type DbIntFieldHelper struct {
	Value int
}

func (u *DbIntFieldHelper) MapFields(fetch *sql.Rows) error {
	err := fetch.Scan(&u.Value)
	if err != nil {
		return err
	}

	return nil
}

func GetIntValue(db *sql.DB, stmt string, params ...interface{}) (int, error) {
	dbh := NewBAMySQL(db)
	data, err := dbh.FetchOne(func() interface{} { return &DbIntFieldHelper{} }, stmt, params...)

	if err != nil {
		return -1, err
	} else {
		var i *DbIntFieldHelper = data.(*DbIntFieldHelper)
		return i.Value, err
	}
}

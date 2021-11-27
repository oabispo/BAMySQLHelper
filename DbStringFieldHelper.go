package bamysqlhelper

import "database/sql"

type DbStringFieldHelper struct {
	Value string
}

func (u *DbStringFieldHelper) MapFields(fetch *sql.Rows) error {
	err := fetch.Scan(&u.Value)
	if err != nil {
		return err
	}

	return nil
}

func GetStringValue(db *sql.DB, stmt string, params ...interface{}) (string, error) {
	dbh := NewBAMySQL(db)
	data, err := dbh.FetchOne(func() interface{} { return &DbStringFieldHelper{} }, stmt, params...)

	if err != nil {
		return "", err
	} else {
		var i *DbStringFieldHelper = data.(*DbStringFieldHelper)
		return i.Value, err
	}
}

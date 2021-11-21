package bamysqlhelper

import (
	"database/sql"
	"testing"
	"time"
)

type testTable struct {
	numberField int
	textField   string
	dateField   time.Time
	timeField   time.Time
	floatField  float32
}

func (p *testTable) MapFields(fetch *sql.Rows) error {
	err := fetch.Scan(&p.numberField, &p.textField, &p.dateField, &p.timeField)
	return err
}

func TestCRUD(t *testing.T) {
	helper := NewBAHelper(NewSQLConnection("www.db4free.net", "sbcash", "bcash", "pentasia119"))
	callback := func() interface{} {
		return &testTable{}
	}

	if _, err := helper.RunSQL("CREATE TEMPORARY TABLE TestTable(numberField int, textField varchar(80), dateField date, timeField TIME, floatField FLOAT);", nil); err != nil {
		t.Fatal(err)
	}

	if _, err := helper.Insert("insert into testTable (numberField, textField, dateField, timeField, floatField) values (1, 'Hello world', NOW(), CURTIME(), 1.99)", nil); err != nil {
		t.Fatal(err)
	}

	if _, err := helper.Update("update TestTable set floatField = ?, textField = ? where numberField = ?", 3.99, "Changed", 1); err != nil {
		t.Fatal(err)
	}

	if _, err := helper.FetchMany(callback, "select * from where numberField = ?", 3.99, "Changed", 1); err != nil {
		t.Fatal(err)
	}

	if _, err := helper.Delete("delete from testTable where numberField = ?", 1); err != nil {
		t.Fatal(err)
	}

	//	data, err := helper.FetchOne(callback, "select 1 from 1 where id = ?", 1)
}

/**
CREATE TEMPORARY TABLE TestTable(numberField int, textField varchar(80), dateField date, timeField TIME, floatField FLOAT);

insert into testTable (numberField, textField, dateField, timeField, floatField) values (1, "Hello world", NOW(), CURTIME(), 1.99)
**/

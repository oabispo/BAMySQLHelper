package bamysqlhelper

import (
	"database/sql"
	"fmt"
	"testing"
	"time"
)

type testTable struct {
	numberField int       `json:number`
	textField   string    `json:string`
	dateField   time.Time `json:date`
	timeField   time.Time `json:datetime`
	floatField  float32   `json:number`
}

func (p *testTable) MapFields(fetch *sql.Rows) error {
	var datetime, timetime string

	err := fetch.Scan(&p.numberField, &p.textField, &datetime, &timetime, &p.floatField)
	if data, err := time.Parse("02/1/2006", datetime); err == nil {
		p.dateField = data
	}

	if data, err := time.Parse("15:04:05", timetime); err == nil {
		p.timeField = data
	}

	return err
}

func convertRaw(data []interface{}) []*testTable {
	var items []*testTable = make([]*testTable, len(data))

	for _, iterator := range data {
		item := iterator.(*testTable)
		items = append(items, item)
	}
	return items
}

func printStruct(items []*testTable) {
	for _, item := range items {
		fmt.Printf("\n%+v", item)
	}
}

func TestCRUD(t *testing.T) {
	helper := NewBAMySQL(NewSQLConnection("sql10.freesqldatabase.com", "sql10452712", "sql10452712", "vBCfnfmcmj"))
	callback := func() interface{} {
		return &testTable{}
	}

	if data, err := helper.Delete("delete from TestTable where numberField = ?", 1); err == nil {
		fmt.Printf("\n%+v row(s) deleted", data)
	} else {
		t.Fatal(err)
	}

	if data, err := helper.FetchMany(callback, "select * from TestTable where numberField = ?", 1); err == nil {
		items := convertRaw(data)
		if len(items) > 0 {
			printStruct(items)
		}
	} else {
		t.Fatal(err)
	}

	if data, err := helper.Insert("insert into TestTable (numberField, textField, dateField, timeField, floatField) values (?, ?, ?, CURTIME(), ?)", 1, "Hello world", time.Now(), 1.99); err == nil {
		fmt.Printf("\nid: %+v %T", data, data)
	} else {
		t.Fatal(err)
	}

	if data, err := helper.FetchMany(callback, "select * from TestTable where numberField = ?", 1); err == nil {
		items := convertRaw(data)
		printStruct(items)
	} else {
		t.Fatal(err)
	}

	if data, err := helper.Update("update TestTable set floatField = ?, textField = ? where numberField = ?", 4.99, "texto alterado", 1); err == nil {
		fmt.Printf("\n%+v row(s) updated", data)
	} else {
		t.Fatal(err)
	}

	if data, err := helper.FetchMany(callback, "select * from TestTable where numberField = ?", 1); err == nil {
		items := convertRaw(data)
		printStruct(items)
	} else {
		t.Fatal(err)
	}
}

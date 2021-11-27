package bamysqlhelper

import (
	"database/sql"
	"fmt"
	"testing"
	"time"
)

type testTable struct {
	NumberField int       `json:number`
	TextField   string    `json:string`
	DateField   time.Time `json:date`
	TimeField   time.Time `json:datetime`
	FloatField  float32   `json:number`
}

func (p *testTable) MapFields(fetch *sql.Rows) error {
	var datetime, timetime string

	err := fetch.Scan(&p.NumberField, &p.TextField, &datetime, &timetime, &p.FloatField)
	if data, err := time.Parse("02/1/2006", datetime); err == nil {
		p.DateField = data
	}

	if data, err := time.Parse("15:04:05", timetime); err == nil {
		p.TimeField = data
	}

	return err
}

func convertRaw(data []interface{}) []*testTable {
	var items []*testTable = make([]*testTable, 0, len(data))

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
	db := NewSQLConnection("sql10.freesqldatabase.com", "sql10452712", "sql10452712", "vBCfnfmcmj")
	helper := NewBASQL(db)
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

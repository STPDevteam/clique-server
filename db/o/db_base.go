package o

import (
	"database/sql"
	oo "github.com/Anna2024/liboo"
	"github.com/jmoiron/sqlx"
)

func SelectSqler(table string, w ...[][]interface{}) string {
	return DBPre(table, w).Select()
}

func Count(table string, w ...[][]interface{}) (data int64, err error) {
	sqler := DBPre(table, w).Count()
	err = oo.SqlGet(sqler, &data)
	if err != nil {
		oo.LogW("sqler:%s", sqler)
	}
	return data, err
}

func Insert(table string, m interface{}) error {
	sqler := oo.NewSqler().Table(table).Insert(m)
	err := oo.SqlExec(sqler)
	if err != nil {
		oo.LogW("sqler:%s", sqler)
	}
	return err
}

func InsertBatch(table string, m []map[string]interface{}) error {
	sqler := oo.NewSqler().Table(table).InsertBatch(m)
	err := oo.SqlExec(sqler)
	if err != nil {
		oo.LogW("sqler:%s", sqler)
	}
	return err
}

func Update(table string, v map[string]interface{}, w ...[][]interface{}) error {
	sqler := DBPre(table, w).Update(v)
	err := oo.SqlExec(sqler)
	if err != nil {
		oo.LogW("sqler:%s", sqler)
	}
	return err
}

func InsertTx(tx *sqlx.Tx, table string, m interface{}) (sql.Result, error) {
	sqler := oo.NewSqler().Table(table).Insert(m)
	res, err := oo.SqlxTxExec(tx, sqler)
	if err != nil {
		oo.LogW("sqler:%s", sqler)
	}
	return res, err
}

func InsertBatchTx(tx *sqlx.Tx, table string, m []map[string]interface{}) (sql.Result, error) {
	sqler := oo.NewSqler().Table(table).InsertBatch(m)
	res, err := oo.SqlxTxExec(tx, sqler)
	if err != nil {
		oo.LogW("sqler:%s", sqler)
	}
	return res, err
}

func UpdateTx(tx *sqlx.Tx, table string, v map[string]interface{}, w ...[][]interface{}) (sql.Result, error) {
	sqler := DBPre(table, w).Update(v)
	res, err := oo.SqlxTxExec(tx, sqler)
	if err != nil {
		oo.LogW("sqler:%s", sqler)
	}
	return res, err
}

func Delete(table string, w ...[][]interface{}) error {
	sqler := DBPre(table, w).Delete()
	err := oo.SqlExec(sqler)
	if err != nil {
		oo.LogW("sqler:%s", sqler)
	}
	return err
}

func W(w ...interface{}) [][]interface{} {
	return [][]interface{}{w}
}

func DBPre(table string, args [][][]interface{}) *oo.Sqler {
	sqler := oo.NewSqler().Table(table)
	for x := range args {
		for y := range args[x] {
			if len(args[x][y]) == 1 {
				sqler.Where(args[x][y][0])
			} else if len(args[x][y]) == 2 {
				sqler.Where(args[x][y][0], args[x][y][1])
			} else if len(args[x][y]) == 3 {
				sqler.Where(args[x][y][0], args[x][y][1], args[x][y][2])
			}
		}
	}
	return sqler
}

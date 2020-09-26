package models

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/dubuqingfeng/bit-node-crawler/dbs"
	"github.com/dubuqingfeng/bit-node-crawler/utils"
	"os"
)

func ExportDatabase(table string) ([]string, [][]string) {
	conn := utils.Config.GlobalDatabase.Write.Name
	if !dbs.CheckDBConnExists(conn) {
		panic(errors.New("not found this conn"))
	}
	rows, _ := dbs.DBMaps[conn].Query(fmt.Sprintf("SELECT * from %s", table))
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var totalValues [][]string
	for rows.Next() {
		var s []string
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		for _, v := range values {
			s = append(s, string(v))
		}
		totalValues = append(totalValues, s)
	}

	if err = rows.Err(); err != nil {
		panic(err.Error())
	}
	return columns, totalValues
}

func WriteToCSV(file string, columns []string, totalValues [][]string) {
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	f.WriteString("\xEF\xBB\xBF")
	defer f.Close()
	w := csv.NewWriter(f)
	for a, i := range totalValues {
		if a == 0 {
			w.Write(columns)
			w.Write(i)
		} else {
			w.Write(i)
		}
	}
	w.Flush()
	fmt.Println("处理完毕：", file)
}

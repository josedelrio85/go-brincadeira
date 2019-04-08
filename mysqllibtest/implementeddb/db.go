package implementeddb

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/fatih/structs"
	_ "github.com/go-sql-driver/mysql"
)

type Where struct {
	Campos       []string
	CamposWhere  map[string]string
	CamposInsert []string
}

func OpenConnection(connString string) (*sql.DB, error) {
	db, err := sql.Open("mysql", connString)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func CreateQueryWithClass(tabla string, class interface{}) string {
	n := structs.Names(class)

	sql := "SELECT "
	for i := 0; i < len(n); i++ {
		sql += n[i] + ", "
	}
	sql = strings.TrimSuffix(sql, ", ")
	sql += " FROM " + tabla + ";"

	return sql
}

func CreateQueryWithFields(tabla string, fields []string) string {
	sql := "SELECT "
	for i := 0; i < len(fields); i++ {
		sql += fields[i] + ", "
	}
	sql = strings.TrimSuffix(sql, ", ")
	sql += " FROM " + tabla + ";"

	return sql
}

func CreateQuery(tabla string, where Where) string {
	sql := "SELECT "
	for i := 0; i < len(where.Campos); i++ {
		sql += where.Campos[i] + ", "
	}
	sql = strings.TrimSuffix(sql, ", ")
	sql += " FROM " + tabla

	tam := len(where.CamposWhere)
	if tam > 0 {
		sqlWhere := " WHERE "
		for k, v := range where.CamposWhere {
			fmt.Println(k)
			fmt.Println(v)
			sqlWhere += k + " = " + v

			if tam > 1 {
				sqlWhere += " AND "
			}
		}
		sqlWhere = strings.TrimSuffix(sqlWhere, "AND ")
		sql += sqlWhere
	}
	fmt.Println(sql)
	return sql
}

func QueryResultToMap(db *sql.DB, query string) (map[int]map[string]string, error) {

	rows, err := db.Query(query)
	cols, err := rows.Columns()
	fr := map[int]map[string]string{}
	id := 0

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		tmpstruct := map[string]string{}

		for i, colName := range cols {
			var v interface{}
			val := columns[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			tmpstruct[colName] = fmt.Sprintf("%s", v)
		}

		fr[id] = tmpstruct
		id++
	}
	// fmt.Println(fr)
	return fr, nil
}

func InsertQuery(tabla string, where Where) string {

	sql := "INSERT INTO " + tabla + " ( "
	for _, k := range where.Campos {
		sql += k + ","
	}
	sql = strings.TrimSuffix(sql, ", ")
	sql += " ) VALUES "

	tam := len(where.CamposInsert)
	if tam > 0 {
		sqlWhere := " ( "
		for _, v := range where.CamposInsert {

			sqlWhere += "'" + v + "'"

			if tam > 1 {
				sqlWhere += ", "
			}
		}
		sqlWhere = strings.TrimSuffix(sqlWhere, ", ")
		sql += sqlWhere + " ) "
	}
	return sql
}

func UpdateQuery(tabla string, where Where) string {
	sql := "UPDATE " + tabla + " SET "

	tam := len(where.CamposInsert)
	if tam > 0 {
		sqlUpdate := ""
		for k, v := range where.CamposInsert {
			sqlUpdate += where.Campos[k] + " = '" + v + "'"

			if tam > 1 {
				sqlUpdate += ", "
			}
		}
		sqlUpdate = strings.TrimSuffix(sqlUpdate, ", ")
		sql += sqlUpdate
	}

	if len(where.CamposWhere) > 0 {
		sqlWhere := " WHERE "
		for k, v := range where.CamposWhere {

			sqlWhere += k + "  = '" + v + "'"

			if tam > 1 {
				sqlWhere += " AND "
			}
		}
		sqlWhere = strings.TrimSuffix(sqlWhere, " AND ")
		sql += sqlWhere
	}

	fmt.Println(sql)
	return sql
}

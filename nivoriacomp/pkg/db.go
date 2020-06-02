package nivoriacomp

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // go mysql driver
)

// Wsmsql is a struct to manage Mysql environment configuration.
type Wsmsql struct {
	Connstring string
	db         *sql.DB
}

// Storer is an interface that declares OpenConnection, BatchInsert and SelectForRequest methods.
type Storer interface {
	Open() error
	BatchInsert([]Xmlstruct) error
	SelectForRequest() ([]Inputdata, error)
}

// Open opens a Mysql connection using a connstring parameter.
// Returns a db instance and nil if success or nil and Error instance if fails.
func (w *Wsmsql) Open() error {

	db, err := sql.Open("mysql", w.Connstring)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	w.db = db
	return nil
}

// BatchInsert iterates over an array of Xmlstruct struct, generates an sql insert
// statement and tries to store it in Mysql environment.
func (w *Wsmsql) BatchInsert(rows []Xmlstruct) error {
	splits := 10
	chunksize := (len(rows) + splits - 1) / splits

	querystr := "insert into webservice.evo_origen_idcliente (clientid, origen, suborigen, stepid, createddate) values %s "
	for i := 0; i < len(rows); i += chunksize {
		end := i + chunksize
		if end > len(rows) {
			end = len(rows)
		}

		subrows := rows[i:end]
		final := make([]string, len(subrows))
		finalargs := []interface{}{}

		for z, row := range subrows {
			r := row
			t, timerr := time.Parse("2006-01-02 15:04:05", r.CreDate)
			if timerr != nil {
				return timerr
			}

			tam := reflect.TypeOf(Xmlstruct{}).NumField()
			valuestrings := make([]string, 0, tam)
			for i := 0; i < tam; i++ {
				valuestrings = append(valuestrings, "?")
			}
			// concatenate ? with , and delete the last , => encapsule all the string into ()
			final[z] = "(" + strings.TrimSuffix(strings.Join(valuestrings, ","), ",") + ")"

			finalargs = append(
				finalargs,
				r.Idsf,
				r.Org,
				r.Suborg,
				r.Stepid.Name,
				t.Format("2006-01-02"),
			)
		}
		// creates the sentence to prepare => insert into ... (?,...,?), ... ,(?,...,?)
		stmtstr := fmt.Sprintf(querystr, strings.Join(final, ","))

		stmt, _ := w.db.Prepare(stmtstr)
		if _, stmterr := stmt.Exec(finalargs...); stmterr != nil {
			return stmterr
		}
		defer stmt.Close()
		finalargs = nil
	}
	return nil
}

// SelectForRequest queries in evo_events_sf_v2_pro table looking for the records that matches
// createddate field is equal to yesterday.
// Returns an array of Inputdata struct.
func (w *Wsmsql) SelectForRequest() ([]Inputdata, error) {
	yesterday := time.Now().Add(time.Duration(-24) * time.Hour)
	sqlselect := fmt.Sprintf(`select CLIENTID, CREATEDDATE FROM webservice.evo_events_sf_v2_pro 
	where date(CREATEDDATE) = ? group by CLIENTID;`)

	stmt, _ := w.db.Prepare(sqlselect)
	rows, stmterr := stmt.Query(yesterday.Format("2006-01-02"))
	// rows, stmterr := stmt.Query("2020-05-18")
	if stmterr != nil {
		return nil, stmterr
	}
	importer := Importer{}
	inputdata := Inputdata{}

	for rows.Next() {
		if err := rows.Scan(&inputdata.Clientid, &inputdata.Createddate); err != nil {
			return nil, err
		}
		importer.Data = append(importer.Data, inputdata)
	}
	return importer.Data, nil
}

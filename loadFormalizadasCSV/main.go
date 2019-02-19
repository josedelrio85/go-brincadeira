package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"go_components/implementeddb"
	"go_components/readparams"
	"log"
	"os"
	"strings"
	"time"
)

// Env is a struct which contains a sql.DB property
type Env struct {
	db *sql.DB
}

func main() {
	f, err := os.OpenFile("../loadFormalizadasCSV_log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	previoformalizadas("formalizadas.csv")
}

func previoformalizadas(filename string) {

	filecsv, ferr := os.Open(filename)
	if ferr != nil {
		log.Println(ferr)
		return
	}
	defer filecsv.Close()

	rows, csverr := csv.NewReader(filecsv).ReadAll()
	if csverr != nil {
		log.Println(csverr)
		return
	}

	// produccion!!!!!!!!3
	connString := readparams.GetConnString(1)
	db, conerr := implementeddb.OpenConnection(connString)
	if conerr != nil {
		log.Println(conerr)
		log.Println(db)
		return
	}

	env := &Env{db: db}

	//WEBSERVICE PRODUCCIÓN
	vuelcaFormalizadas(env.db, rows)
	defer db.Close()

	// report_panel WEBSERVICE!!!!!!!!4
	connString = readparams.GetConnString(1)
	db, err := implementeddb.OpenConnection(connString)
	if err != nil {
		log.Println(err)
		log.Println(db)
		return
	}

	env = &Env{db: db}
	vuelcaFormalizadas(env.db, rows)

	defer db.Close()
}

func vuelcaFormalizadas(db *sql.DB, rows [][]string) {

	altrows := rows[1:]

	// sqlTruncate := "delete from webservice.evo_formalizadas_sf_v2 where date(FECHA_FORMALIZACION) >= '2018-12-01';"
	sqlTruncate := "delete from webservice.evo_formalizadas_sf_v2 where date(FECHA_FORMALIZACION) >= '2019-01-01';"
	_, err := db.Query(sqlTruncate)
	if err != nil {
		log.Println(err)
		return
	}

	sql := "INSERT INTO webservice.evo_formalizadas_sf_v2 (NOMBRE,CLIENTID,NUMERO_PROCESO_CONTRATACION,PRODUCTO,NUMERO_EXPEDIENTE,ID_PERSONA_IRIS,ORIGEN_PROMOCION,FECHA_FORMALIZACION) VALUES %s"
	final := make([]string, len(altrows))
	finalArgs := []interface{}{}

	for z, row := range altrows {
		r := strings.Split(row[0], ";")
		t, timerr := time.Parse("02/01/2006 15:04", strings.Trim(r[7], " "))
		if timerr != nil {
			log.Println(timerr)
			log.Println(t)
			return
		}

		tam := len(r)
		valueStrings := make([]string, 0, tam)
		for i := 0; i < tam; i++ {
			valueStrings = append(valueStrings, "?")
		}
		final[z] = "(" + strings.TrimSuffix(strings.Join(valueStrings, ","), ",") + ")"
		finalArgs = append(finalArgs, r[0], r[1], r[2], r[3], r[4], r[5], r[6], t)
	}

	stmtStr := fmt.Sprintf(sql, strings.Join(final, ","))
	stmt, _ := db.Prepare(stmtStr)
	_, stmterr := stmt.Exec(finalArgs...)
	if stmterr != nil {
		log.Println(stmterr)
		fmt.Println(stmterr)
		return
	}
}

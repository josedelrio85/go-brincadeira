package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/bysidecar/go_components/implementeddb"
	"github.com/bysidecar/go_components/readparams"
)

// Env is a struct which contains a sql.DB property
type Env struct {
	db *sql.DB
}

func main() {

	var basepath = flag.String("basepath", "C:\\Users\\Jose\\go\\src\\github.com\\bysidecar\\go_components\\loadFormalizadasCSV", "path to read the posted file")
	var fileconfig = flag.String("fileconfig", "C:\\Users\\Jose\\go\\src\\github.com\\bysidecar\\go_components\\readparams", "path where to read config file")
	flag.Parse()

	f, err := os.OpenFile("../loadFormalizadasCSV_log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	previoformalizadas("formalizadas.csv", *basepath, *fileconfig)
}

func previoformalizadas(file string, basepath string, fileconfig string) {

	filepath := path.Join(basepath, file)

	filecsv, ferr := os.Open(filepath)
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
	connString := readparams.GetConnString(3, fileconfig)
	db, conerr := implementeddb.OpenConnection(connString)
	if conerr != nil {
		log.Println(conerr)
		log.Println(db)
		return
	}

	env := &Env{db: db}

	//WEBSERVICE PRODUCCIÓN
	vuelcaFormalizadas(env.db, rows)
	num := cuentaVolcadas(env.db)
	fmt.Println(num)
	defer db.Close()

	// report_panel WEBSERVICE!!!!!!!!4
	connString = readparams.GetConnString(4, fileconfig)
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

type excelRow struct {
	nombre                    string
	clientid                  string
	numeroprocesocontratacion string
	producto                  string
	numeroexpediente          string
	idpersonairis             string
	origenpromocion           string
	fechaformalizacion        time.Time
}

func cuentaVolcadas(db *sql.DB) (count int) {
	sql := "select count(*) as count from webservice.evo_formalizadas_sf_v2 where date(FECHA_FORMALIZACION) >= '2019-07-01';"

	rows, err := db.Query(sql)
	if err != nil {
		log.Println(err)
		return 0
	}

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			log.Println(err)
			return 0
		}
	}
	return count
}

func vuelcaFormalizadas(db *sql.DB, rows [][]string) {

	altrows := rows[1:]

	sqlTruncate := "delete from webservice.evo_formalizadas_sf_v2 where date(FECHA_FORMALIZACION) >= '2019-07-01';"
	_, err := db.Query(sqlTruncate)
	if err != nil {
		log.Println(err)
		return
	}

	sql := "INSERT INTO webservice.evo_formalizadas_sf_v2 (NOMBRE,CLIENTID,NUMERO_PROCESO_CONTRATACION,PRODUCTO,NUMERO_EXPEDIENTE,ID_PERSONA_IRIS,ORIGEN_PROMOCION,FECHA_FORMALIZACION) VALUES %s "

	splits := 1000
	chunkSize := (len(altrows) + splits - 1) / splits

	for i := 0; i < len(altrows); i += chunkSize {
		end := i + chunkSize
		if end > len(altrows) {
			end = len(altrows)
		}

		subrows := altrows[i:end]
		final := make([]string, len(subrows))
		finalArgs := []interface{}{}

		for z, row := range subrows {
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

			data := excelRow{
				nombre:                    r[0],
				clientid:                  r[1],
				numeroprocesocontratacion: r[2],
				producto:                  r[3],
				numeroexpediente:          r[4],
				idpersonairis:             r[5],
				origenpromocion:           r[6],
				fechaformalizacion:        t,
			}

			finalArgs = append(
				finalArgs,
				data.nombre,
				data.clientid,
				data.numeroprocesocontratacion,
				data.producto,
				data.numeroexpediente,
				data.idpersonairis,
				data.origenpromocion,
				data.fechaformalizacion.Format("2006-01-02"),
			)
		}

		stmtStr := fmt.Sprintf(sql, strings.Join(final, ","))
		stmt, err := db.Prepare(stmtStr)
		if err != nil {
			log.Println(err)
			return
		}
		_, stmterr := stmt.Exec(finalArgs...)
		if stmterr != nil {
			log.Println(stmterr)
			return
		}
		defer stmt.Close()
		// reset finalArgs
		finalArgs = nil
	}
}

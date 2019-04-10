package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
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
	connString := readparams.GetConnString(3)
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
	connString = readparams.GetConnString(4)
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

func vuelcaFormalizadas(db *sql.DB, rows [][]string) {

	altrows := rows[1:]

	sqlTruncate := "delete from webservice.evo_formalizadas_sf_v2 where date(FECHA_FORMALIZACION) >= '2019-01-01';"
	_, err := db.Query(sqlTruncate)
	if err != nil {
		log.Println(err)
		return
	}
	sql := "INSERT INTO webservice.evo_formalizadas_sf_v2 (NOMBRE,CLIENTID,NUMERO_PROCESO_CONTRATACION,PRODUCTO,NUMERO_EXPEDIENTE,ID_PERSONA_IRIS,ORIGEN_PROMOCION,FECHA_FORMALIZACION) VALUES "
	sqlFinal := ""

	splits := 10
	chunkSize := (len(altrows) + splits - 1) / splits

	for i := 0; i < len(altrows); i += chunkSize {
		end := i + chunkSize

		if end > len(altrows) {
			end = len(altrows)
		}
		subrows := altrows[i:end]
		sqlWhere := ""

		for _, row := range subrows {
			r := strings.Split(row[0], ";")

			t, timerr := time.Parse("02/01/2006 15:04", strings.Trim(r[7], " "))
			if timerr != nil {
				log.Println(timerr)
				log.Println(t)
				return
			}

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

			a := " ('" + data.nombre + "','" + data.clientid + "', '" + data.numeroprocesocontratacion + "', '" + data.producto + "', '" + data.numeroexpediente + "', '" + data.idpersonairis + "', '" + data.origenpromocion + "', '" + data.fechaformalizacion.Format("2006-01-02") + "'), "
			sqlWhere += a
		}

		if sqlWhere != "" {
			sqlWhere = strings.TrimSuffix(sqlWhere, ", ")
			sqlFinal = sql + sqlWhere + " ; "

			_, err = db.Query(sqlFinal)
			if err != nil {
				log.Println(err)
				log.Println(sqlFinal)
				return
			}
			sqlWhere = ""
			sqlFinal = ""
		}
	}
}

// When the array to proccess has more than 8191 rows, the proccess fails, generating panic: runtime error: invalid memory address or nil pointer dereference
// reason ?????????????
func vuelcaFormalizadasWithPreparedStatement(db *sql.DB, rows [][]string) {

	altrows := rows[1:]

	sqlTruncate := "delete from webservice.evo_formalizadas_sf_v2 where date(FECHA_FORMALIZACION) >= '2019-01-01';"
	_, err := db.Query(sqlTruncate)
	if err != nil {
		log.Println(err)
		return
	}

	sql := "INSERT INTO webservice.evo_formalizadas_sf_v2 (NOMBRE,CLIENTID,NUMERO_PROCESO_CONTRATACION,PRODUCTO,NUMERO_EXPEDIENTE,ID_PERSONA_IRIS,ORIGEN_PROMOCION,FECHA_FORMALIZACION) VALUES %s "
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

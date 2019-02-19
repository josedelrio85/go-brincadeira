package main

import (
	"database/sql"
	"encoding/csv"
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
	f, err := os.OpenFile("../loadFirmasCSV_log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	previofirmadas("firmas.csv")
}

func previofirmadas(filename string) {

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
	vuelcaFirmadas(env.db, rows)
	defer db.Close()

	// report_panel WEBSERVICE!!!!!!!!4
	// connString = readparams.GetConnString(1)
	// db, err := implementeddb.OpenConnection(connString)
	// if err != nil {
	// 	log.Println(err)
	// 	log.Println(db)
	// 	return
	// }

	// env = &Env{db: db}
	// vuelcaFirmadas(env.db, rows)

	// defer db.Close()
}

// ExcelRow is a struct with the structure of the excel file to analyze.
type ExcelRow struct {
	producto         string
	idclienteevo     string
	fechacreacion    time.Time
	estado           string
	ultimopunto      string
	gestioncaptacion string
	numerocaso       string
	estadoalt        string
	motivodes        string
	numerologalty    string
	idpersonairis    string
	numeroexp        string
	clasecliente     string
	fechafirma       time.Time
	tipoidentific    string
}

func vuelcaFirmadas(db *sql.DB, rows [][]string) {

	altrows := rows[1:]

	sqlTruncate := "delete from webservice.evo_firmados_sf_v2 where date(Fecha_de_firma) >= '2018-12-01';"
	_, err := db.Query(sqlTruncate)
	if err != nil {
		log.Println(err)
		return
	}

	sql := "INSERT INTO webservice.evo_firmados_sf_v2 (Producto,ID_Cliente_EVO,Fecha_de_creacion, Estado_cliente, Ultimo_punto_de_abandono, Gestion_Captacion, Numero_del_proceso_de_contratacion, Estado, Motivo_desestimacion, Numero_de_Logalty, ID_Persona_Iris, Numero_Expediente, Clase_de_Cliente, Fecha_de_firma, Tipo_de_identificacion) VALUES "
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

			t, timerr := time.Parse("02/01/2006", r[2])
			t2, timerr := time.Parse("02/01/2006", r[13])

			if timerr != nil {
				log.Println(timerr)
				log.Println(t)
				return
			}

			data := ExcelRow{
				producto:         r[0],
				idclienteevo:     r[1],
				fechacreacion:    t,
				estado:           r[3],
				ultimopunto:      r[4],
				gestioncaptacion: r[5],
				numerocaso:       r[6],
				estadoalt:        r[7],
				motivodes:        r[8],
				numerologalty:    r[9],
				idpersonairis:    r[10],
				numeroexp:        r[11],
				clasecliente:     r[12],
				fechafirma:       t2,
				tipoidentific:    r[14],
			}

			a := " ('" + data.producto + "','" + data.idclienteevo + "', '" + data.fechacreacion.Format("2006-01-02") + "', '" + data.estado + "', '" + data.ultimopunto + "', '" + data.gestioncaptacion + "', '" + data.numerocaso + "', '" + data.estadoalt + "', '" + data.motivodes + "', '" + data.numerologalty + "', '" + data.idpersonairis + "', '" + data.numeroexp + "', '" + data.clasecliente + "', '" + data.fechafirma.Format("2006-01-02") + "', '" + data.tipoidentific + "'), "
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

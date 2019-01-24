package main

import (
	"database/sql"
	"dev/mysqllibtest/implementeddb"
	"dev/readparams"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

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

type Env struct {
	db *sql.DB
}

var env *Env

func main() {
	f, err := os.OpenFile("../loadFirmasExcel_log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	previofirmadas("firmas.xlsx")
}

func previofirmadas(filename string) {
	//PARAM!
	xlsx, err := excelize.OpenFile(filename)
	if err != nil {
		log.Println(err)
		return
	}

	// produccion!!!!!!!!
	connString := readparams.GetConnString(3)
	db, err := implementeddb.OpenConnection(connString)
	if err != nil {
		log.Println(err)
		return
	}

	env = &Env{db: db}

	rows := [][]string{}
	for _, name := range xlsx.GetSheetMap() {
		rows = xlsx.GetRows(name)
	}

	vuelcaFirmadas(env.db, rows)
	defer db.Close()

	// report_panel WEBSERVICE!!!!!!!!
	connString = readparams.GetConnString(4)
	db, err = implementeddb.OpenConnection(connString)
	if err != nil {
		log.Println(err)
		return
	}

	env = &Env{db: db}
	vuelcaFirmadas(env.db, rows)

	defer db.Close()
}

func vuelcaFirmadas(db *sql.DB, rows [][]string) {

	altrows := rows[1:]
	tam := (len(rows) / 10)

	sqlTruncate := "delete from webservice.evo_firmados_sf_v2 where date(Fecha_de_firma) >= '2018-12-01';"
	_, err := db.Query(sqlTruncate)
	if err != nil {
		log.Println(err)
		return
	}

	sql := "INSERT INTO webservice.evo_firmados_sf_v2 (Producto,ID_Cliente_EVO,Fecha_de_creacion, Estado_cliente, Ultimo_punto_de_abandono, Gestion_Captacion, Numero_del_proceso_de_contratacion, Estado, Motivo_desestimacion, Numero_de_Logalty, ID_Persona_Iris, Numero_Expediente, Clase_de_Cliente, Fecha_de_firma, Tipo_de_identificacion) VALUES "
	var sqlWhere string
	var sqlAlt string
	fmt.Print(sqlAlt)
	aux := 0

	for i, row := range altrows {
		t, err := time.Parse("01-02-06", row[2])
		t2, err := time.Parse("01-02-06", row[13])

		if err != nil {
			log.Println(err)
			log.Println(t)
			return
		}

		data := ExcelRow{
			producto:         row[0],
			idclienteevo:     row[1],
			fechacreacion:    t,
			estado:           row[3],
			ultimopunto:      row[4],
			gestioncaptacion: row[5],
			numerocaso:       row[6],
			estadoalt:        row[7],
			motivodes:        row[8],
			numerologalty:    row[9],
			idpersonairis:    row[10],
			numeroexp:        row[11],
			clasecliente:     row[12],
			fechafirma:       t2,
			tipoidentific:    row[14],
		}

		a := " ('" + data.producto + "','" + data.idclienteevo + "', '" + data.fechacreacion.Format("2006-01-02") + "', '" + data.estado + "', '" + data.ultimopunto + "', '" + data.gestioncaptacion + "', '" + data.numerocaso + "', '" + data.estadoalt + "', '" + data.motivodes + "', '" + data.numerologalty + "', '" + data.idpersonairis + "', '" + data.numeroexp + "', '" + data.clasecliente + "', '" + data.fechafirma.Format("2006-01-02") + "', '" + data.tipoidentific + "'), "
		sqlWhere += a
		if aux == i && aux > 0 {
			sqlWhere = strings.TrimSuffix(sqlWhere, ", ")
			sqlAlt = sql + sqlWhere + " ; "

			log.Println(sqlAlt)

			_, err = db.Query(sqlAlt)
			if err != nil {
				log.Println(err)
				log.Println(sqlAlt)
				return
			}
			sqlWhere = ""
			sqlAlt = ""
			aux += tam
		}
		if aux == 0 {
			aux++
		}
	}

	if sqlWhere != "" {
		sqlWhere = strings.TrimSuffix(sqlWhere, ", ")
		sql += sqlWhere + " ; "
		log.Println(sqlWhere)

		_, err = db.Query(sql)
		if err != nil {
			log.Println(err)
			log.Println(sql)
			return
		}
	}

}

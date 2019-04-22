package models

import (
	"database/sql"

	"github.com/bysidecar/go_components/mysqllibtest/functions"
	"github.com/bysidecar/go_components/mysqllibtest/implementeddb"
)

type Sources struct {
	Souid          string `json:"sou_id"`
	Soudescription string `json:"sou_description"`
	Souactive      string `json:"sou_active"`
}

func indexAlt(db *sql.DB) ([]*Sources, error) {
	query := implementeddb.CreateQueryWithClass("report_panel.sources", Sources{Souid: ""})

	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	elements := make([]*Sources, 0)
	for rows.Next() {
		el := new(Sources)
		err := rows.Scan(&el.Souid, &el.Soudescription)
		if err != nil {
			return nil, err
		}
		elements = append(elements, el)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return elements, nil
}

func Index(db *sql.DB) ([]Sources, error) {

	campos := []string{
		"sou_id",
		"sou_description",
		"sou_active",
	}

	camposWhere := map[string]string{"sou_active": "1"}
	a := implementeddb.Where{Campos: campos, CamposWhere: camposWhere}
	query := implementeddb.CreateQuery("report_panel.sources", a)

	result, err := implementeddb.QueryResultToMap(db, query)
	if err != nil {
		return nil, err
	}

	var m []Sources
	for _, r := range result {
		res := Sources{
			r["sou_id"],
			r["sou_description"],
			r["sou_active"],
		}
		m = append(m, res)
	}
	return m, nil
}

func Get(db *sql.DB, id string) ([]Sources, error) {

	campos := []string{
		"sou_id",
		"sou_description",
		"sou_active",
	}

	camposWhere := map[string]string{"sou_active": "1", "sou_id": id}
	a := implementeddb.Where{Campos: campos, CamposWhere: camposWhere}
	query := implementeddb.CreateQuery("report_panel.sources", a)

	result, err := implementeddb.QueryResultToMap(db, query)
	if err != nil {
		return nil, err
	}

	var m []Sources
	for _, r := range result {
		res := Sources{
			r["sou_id"],
			r["sou_description"],
			r["sou_active"],
		}
		m = append(m, res)
	}
	return m, nil
}

func Insert(db *sql.DB, source Sources) (bool, error) {

	campos, values := functions.GetJsonKeysReceived(&source)

	a := implementeddb.Where{Campos: campos, CamposWhere: nil, CamposInsert: values}
	query := implementeddb.InsertQuery("report_panel.sources", a)

	rows, err := db.Query(query)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	defer db.Close()

	return true, nil
}

func Update(db *sql.DB, source Sources) (bool, error) {

	campos, values := functions.GetJsonKeysReceived(&source)
	camposWhere := map[string]string{"sou_id": source.Souid}

	a := implementeddb.Where{Campos: campos, CamposWhere: camposWhere, CamposInsert: values}
	query := implementeddb.UpdateQuery("report_panel.sources", a)

	rows, err := db.Query(query)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	defer db.Close()

	return true, nil
}

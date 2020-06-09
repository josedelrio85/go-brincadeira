package nivoriacomp

import (
	"encoding/csv"
	"os"
	"strings"
)

// Importer is the structure to handle the needed resources to obtain the input data.
type Importer struct {
	Path   string
	Date   string
	Data   []Inputdata
	Storer Storer
}

// Inputdata represents the data that will be imported to db.
// It can come from an CSV file or a select query.
type Inputdata struct {
	Clientid    string
	Createddate string
	Stepid      string
}

// Importfromcsv obtains the resources used as input from a CSV file.
func (a *Importer) Importfromcsv() error {
	filecsv, err := os.Open(a.Path)
	if err != nil {
		return err
	}
	defer filecsv.Close()

	data, err := csv.NewReader(filecsv).ReadAll()
	if err != nil {
		return err
	}

	data = data[1:]
	for _, r := range data {
		row := strings.Split(r[0], ";")
		input := Inputdata{
			Clientid:    row[0],
			Createddate: row[1],
			Stepid:      row[2],
		}
		a.Data = append(a.Data, input)
	}
	return nil
}

// Importfromdb obtains the resources used as input from DB.
func (a *Importer) Importfromdb() error {
	data, err := a.Storer.SelectForRequest(a.Date)
	if err != nil {
		return err
	}
	a.Data = data
	return nil
}

package main

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	readconfig()
}

// Note: struct fields must be public in order for unmarshal to correctly populate the data.
type Fileconfig struct {
	Source struct {
		Tipo      string `yaml:"type"`
		Path      string `yaml:"path"`
		Extension string `yaml:"extension"`
	} `yaml:"source"`
	Destination struct {
		Tipo       string `yaml:"type"`
		Tabla      string `yaml:"table"`
		Estructura []struct {
			Field1 string `yaml:"field1"`
			Field2 string `yaml:"field2"`
			Field3 string `yaml:"field3"`
			Field4 string `yaml:"field4"`
		} `yaml:"structure"`
	} `yaml:"destination"`
}

func readconfig() {
	// 	TO DO: We should pass file path by parameter

	// file, err := os.Open("./config/config_xlsx.yaml")
	// file, err := os.Open("./config/config_csv.yaml")
	// file, err := os.Open("./config/config_json.yaml")
	file, err := os.Open("./config/config_xml.yaml")

	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	fc := Fileconfig{}
	data, err := ioutil.ReadAll(file)

	err = yaml.Unmarshal([]byte(data), &fc)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
	switchextension(fc)
}

func switchextension(fc Fileconfig) {
	var rows map[string]interface{}
	var err error
	path := fc.Source.Path

	switch fc.Source.Extension {
	case "xlsx":
		rows, err = processxlsx(path)
	case "csv":
		rows, err = processcsv(path)
	case "json":
		rows, err = processjson(path)
	case "xml":
		rows, err = processxml(path)
	default:
		panic("switchextension error")
	}

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(rows)

	iterateoverrows(rows)
}

func processcsv(filename string) (res map[string]interface{}, err error) {

	file, ferr := os.Open(filename)
	if ferr != nil {
		log.Println(ferr)
		return nil, ferr
	}
	defer file.Close()

	rows, csverr := csv.NewReader(file).ReadAll()
	if csverr != nil {
		log.Println(csverr)
		return nil, csverr
	}

	res = make(map[string]interface{})
	for k, z := range rows {
		res[strconv.Itoa(k)] = z
	}

	// fmt.Println(res)
	// fmt.Println("---------------")

	// for h, j := range res {
	// 	fmt.Println("h %s \n", h)
	// 	fmt.Println("j %s \n", j)
	// }
	return res, nil
}

func processxlsx(filename string) (res map[string]interface{}, err error) {

	res = make(map[string]interface{})
	file, ferr := excelize.OpenFile(filename)
	if ferr != nil {
		log.Println(err)
		return nil, ferr
	}

	for _, name := range file.GetSheetMap() {
		rows := file.GetRows(name)
		for k, z := range rows {
			res[strconv.Itoa(k)] = z
		}
	}
	return res, nil
}

func processjson(filename string) (res map[string]interface{}, err error) {
	file, ferr := os.Open(filename)
	if ferr != nil {
		log.Println(ferr)
		return nil, ferr
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// var result map[string]interface{}
	res = make(map[string]interface{})
	json.Unmarshal([]byte(byteValue), &res)

	return res, nil
}

func processxml(filename string) (res map[string]interface{}, err error) {
	file, ferr := os.Open(filename)
	if ferr != nil {
		log.Println(ferr)
		return nil, ferr
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res = make(map[string]interface{})
	xml.Unmarshal([]byte(byteValue), &res)

	return res, nil
}

func iterateoverrows(rows map[string]interface{}) {
	// nombres de las columnas
	return
	// cols := rows[:1]
	// fmt.Print(cols)

	// columns := make([]interface{}, len(cols))
	// columnPointers := make([]interface{}, len(cols))
	// for z := range columns {
	// 	columnPointers[z] = &columns[z]
	// }

	// fmt.Println(columnPointers)
	// for _ = range rows {
	// 	columns := make([]interface{}, len(cols))
	// 	columnPointers := make([]interface{}, len(cols))
	// 	for z := range columns {
	// 		columnPointers[z] = &columns[z]
	// 	}

	// 	fmt.Println(columnPointers)
	// }
}

func StrToMap(in string) map[string]interface{} {
	res := make(map[string]interface{})
	array := strings.Split(in, " ")
	temp := make([]string, 2)
	for _, val := range array {
		temp = strings.Split(string(val), ":")
		res[temp[0]] = temp[1]
	}
	return res
}

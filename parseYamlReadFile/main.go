package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sort"
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
		Tipo           string   `yaml:"type"`
		Tabla          string   `yaml:"table"`
		Estructura     []string `yaml:"structure"`
		EstructuraJSON []struct {
			Field string `yaml:"field"`
			Order int    `yaml:"order"`
		} `yaml:"structurejson"`
	} `yaml:"destination"`
}

type RowList struct {
	Keys []int
	Rows map[int][]string
}

func (r RowList) orderedList() {
	for _, k := range r.Keys {
		fmt.Println("key: ", k, " value: ", r.Rows[k])
		// for i, zz := range r.Rows[k] {
		// 	fmt.Println("i: ", i)
		// 	fmt.Println(zz)
		// }
	}
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
	// var rows map[string]interface{}
	var rows map[string][]string
	var err error
	path := fc.Source.Path

	switch fc.Source.Extension {
	case "xlsx":
		// rows, err = processxlsx(path)
	case "csv":
		rows, err = processcsv(path)
	case "json":
		a, ferr := processjson(path)
		if ferr != nil {
			log.Println(ferr)
			return
		}
		rows = hazcosasconinterfaz(a)
	case "xml":
		// rows, err = processxml(path)
	default:
		panic("switchextension error")
	}

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(rows)
}

func processcsv(filename string) (res map[string][]string, err error) {

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

	res = make(map[string][]string)
	for k, z := range rows {
		r := strings.Split(z[0], ";")
		res[strconv.Itoa(k)] = r
	}
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

type OrderedMap struct {
	Order []string
	Map   []interface{}
}

func processjson(filename string) (res []interface{}, err error) {
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

	om := OrderedMap{}
	json.Unmarshal([]byte(byteValue), &om.Map)

	index := make(map[string]int)
	for key := range om.Map {
		om.Order = append(om.Order, string(key))
		esc, _ := json.Marshal(key) //Escape the key
		index[string(key)] = bytes.Index([]byte(byteValue), esc)
		fmt.Println(index[string(key)])
	}
	sort.Slice(om.Order, func(i, j int) bool { return index[om.Order[i]] < index[om.Order[j]] })

	// fmt.Println(om.Map)
	return om.Map, nil
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

func hazcosasconinterfaz(res []interface{}) map[string][]string {

	b := make(map[string][]string)
	value := ""
	for z, r := range res {
		s := reflect.ValueOf(r)

		m, ok := r.(map[string]interface{})
		if !ok {
			fmt.Errorf("want type map[string]interface{};  got %T", s)
		}
		var a []string

		for _, v := range m {
			// fmt.Println(k, "=>", v)
			switch v.(type) {
			case float64:
				value = strconv.FormatFloat(v.(float64), 'f', 6, 64)
			case int:
				value = strconv.FormatInt(v.(int64), 'i')
			case bool, string:
				value = v.(string)
			default:
				// fmt.Println("type unknown") // here v has type interface{}
			}
			a = append(a, value)
		}
		b[string(z)] = a
	}
	return b
}

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
	file, err := os.Open("./config/config_json_case2.yaml")
	// file, err := os.Open("./config/config_xml.yaml")

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

	rowlist := RowList{}
	var rerr error
	var a []byte

	switch fc.Source.Extension {
	case "xlsx":
		rowlist, rerr = readxlsx(fc.Source.Path)
	case "csv":
		rowlist, rerr = readcsv(fc.Source.Path)
	case "json":
		a, rerr = readjson(fc.Source.Path)
		rowlist, rerr = processjson(a, &fc)
	case "xml":
		readxml(fc.Source.Path)
	default:
		panic("switchextension error")
	}

	if rerr != nil {
		log.Println(rerr)
		panic(rerr)
	}

	rowlist.orderedList()
}

func readcsv(filename string) (result RowList, err error) {

	file, ferr := os.Open(filename)
	if ferr != nil {
		log.Println(ferr)
		return RowList{}, ferr
	}
	defer file.Close()

	rows, csverr := csv.NewReader(file).ReadAll()
	if csverr != nil {
		log.Println(csverr)
		return RowList{}, csverr
	}

	res := make(map[int][]string)
	var keys []int
	for k, z := range rows {
		keys = append(keys, k)
		r := strings.Split(z[0], ";")
		res[k] = r
	}
	sort.Ints(keys)

	result.Keys = keys
	result.Rows = res

	return result, nil
}

func readxlsx(filename string) (result RowList, err error) {

	file, ferr := excelize.OpenFile(filename)
	if ferr != nil {
		log.Println(err)
		return RowList{}, ferr
	}

	res := make(map[int][]string)
	var keys []int
	for _, name := range file.GetSheetMap() {
		rows := file.GetRows(name)
		for k, z := range rows {
			keys = append(keys, k)
			res[k] = z
		}
		sort.Ints(keys)
	}

	result.Keys = keys
	result.Rows = res
	return result, nil
}

func readxml(filename string) (res map[string]interface{}, err error) {
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

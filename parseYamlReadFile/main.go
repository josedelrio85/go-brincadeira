package main

import (
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

func readjson(filename string) (data []byte, err error) {
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
	return byteValue, nil
}

func processjson(data []byte, fc *Fileconfig) (result RowList, err error) {

	var v []interface{}
	json.Unmarshal(data, &v)

	result, err = generateorderedestructureforjson(v, fc)
	return result, err
}

func generateorderedestructureforjson(res []interface{}, fc *Fileconfig) (result RowList, err error) {

	rl := make(map[int][]string)
	var keys []int

	for k, z := range res {
		s := reflect.ValueOf(z)
		keys = append(keys, k)

		m, ok := z.(map[string]interface{})
		if !ok {
			fmt.Errorf("want type map[string]interface{};  got %T", s)
		}

		var a []string
		b := make(map[int]string, len(m))
		value := ""

		var keysalt []int
		for i, v := range m {

			switch v.(type) {
			case float64:
				value = fmt.Sprintf("%.0f", v.(float64))
			case int:
				value = strconv.FormatInt(v.(int64), 'i')
			case bool, string:
				value = v.(string)
			default:
				fmt.Println("type unknown")
			}

			for _, estrucjson := range fc.Destination.EstructuraJSON {
				if estrucjson.Field == i {
					keysalt = append(keysalt, estrucjson.Order)
					// fmt.Println("order: ", estrucjson.Order, " field: ", estrucjson.Field, " value: ", value)
					b[estrucjson.Order] = value
				}
			}
			sort.Ints(keysalt)
		}

		for _, k := range keysalt {
			a = append(a, b[k])
		}
		rl[k] = a
	}

	result.Keys = keys
	result.Rows = rl

	return result, nil
}

func (r RowList) insertStatement(fc *Fileconfig) (sql string) {

	sql = "INSERT INTO " + fc.Destination.Tabla + " ( "

	if fc.Source.Extension == "json" {
		for _, k := range fc.Destination.EstructuraJSON {
			sql += k.Field + ","
		}
	} else {
		for _, k := range fc.Destination.Estructura {
			sql += k + ","
		}
	}

	sql = strings.TrimSuffix(sql, ",")
	sql += " ) VALUES "

	for _, k := range r.Keys {
		sqlWhere := " ( "
		tam := len(r.Rows[k])
		if tam > 0 {
			for _, value := range r.Rows[k] {
				sqlWhere += "'" + value + "'"
				if tam > 1 {
					sqlWhere += ", "
				}
			}
			sqlWhere = strings.TrimSuffix(sqlWhere, ", ")
			sql += sqlWhere + " ), "
		}
	}
	sql = strings.TrimSuffix(sql, ", ")
	fmt.Println(sql)
	return sql
}

func (r RowList) insertPreparedStatement(fc *Fileconfig) (stmtStr string, finalArgs []interface{}) {

	sql := "INSERT INTO " + fc.Destination.Tabla + " ( "

	if fc.Source.Extension == "json" {
		for _, k := range fc.Destination.EstructuraJSON {
			sql += k.Field + ","
		}
	} else {
		for _, k := range fc.Destination.Estructura {
			sql += k + ","
		}
	}
	sql = strings.TrimSuffix(sql, ",")
	sql += " ) VALUES %s"
	//////// hasta aquí => inser into (c1, c2, c3, .... ) values (

	final := make([]string, len(r.Keys))
	finalArgs = []interface{}{}

	for _, k := range r.Keys {
		tam := len(r.Rows[k])
		if tam > 0 {
			// create a mini-statement like (?,?,....,?) with the size of the elements of the row
			valueStrings := make([]string, 0, tam)
			for i := 0; i < tam; i++ {
				valueStrings = append(valueStrings, "?")
			}
			final[k] = "(" + strings.TrimSuffix(strings.Join(valueStrings, ","), ",") + ")"
			//array of interfaces with the values that will be executed by stmt.Exec
			finalArgs = append(finalArgs, r.Rows[k])
		}
	}
	stmtStr = fmt.Sprintf(sql, strings.Join(final, ","))

	fmt.Println(stmtStr)
	fmt.Println(finalArgs)

	return stmtStr, finalArgs
}

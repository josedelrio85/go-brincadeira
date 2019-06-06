package main

import (
	"encoding/base64"
	"encoding/csv"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	var filename = flag.String("filename", "file.csv", "file name to read")
	flag.Parse()

	f, err := os.OpenFile("./encode64_log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer f.Close()

	data, err := procFile(*filename)
	if err != nil {
		log.SetOutput(f)
		return
	}

	writeFile(data)
}

func procFile(inputfile string) ([][]string, error) {
	basepath := "./"

	files, err := ioutil.ReadDir(basepath)
	if err != nil {
		log.Printf("Error reading directory %v", err)
		return nil, err
	}

	for _, f := range files {
		if strings.Contains(f.Name(), ".csv") && strings.Contains(f.Name(), "input") {
			inputfile = f.Name()
		}
	}

	filepath := path.Join(basepath, inputfile)

	filecsv, err := os.Open(filepath)
	if err != nil {
		log.Printf("Error reading input file %v", err)
		return nil, err
	}
	defer filecsv.Close()

	rows, csverr := csv.NewReader(filecsv).ReadAll()
	if csverr != nil {
		log.Println(csverr)
		return nil, err
	}

	var array [][]string
	rows = rows[1:]
	for _, row := range rows {
		input := strings.Join(row, "")
		str := base64.StdEncoding.EncodeToString([]byte(input))
		array = append(array, []string{input, str})
	}

	return array, nil
}

func writeFile(data [][]string) {
	output, err := os.Create("output.csv")
	if err != nil {
		log.Printf("Error creating output file %v", err)
		return
	}

	writer := csv.NewWriter(output)
	writer.Comma = ';'
	defer writer.Flush()

	if err := writer.Write([]string{"Phone", "base64"}); err != nil {
		log.Printf("Error writing headers in output file %v", err)
		return
	}
	for _, row := range data {
		if err := writer.Write(row); err != nil {
			log.Printf("Error writing row %v", err)
			return
		}
	}
}

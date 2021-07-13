package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	nivoriacomp "github.com/josedelrio85/go_components/nivoriacomp/pkg"
	"github.com/josedelrio85/go_components/readparams"
)

func main() {
	start := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("Process started at %s", start)

	var basepath = flag.String("basepath", "C:\\Users\\Jose\\go\\src\\github.com\\josedelrio85\\go_components\\nivoriacomp\\creados_EVO.csv", "path where to read csv file")
	// var fileconfig = flag.String("fileconfig", "C:\\Users\\Jose\\go\\src\\github.com\\josedelrio85\\go_components\\readparams", "path where to read config file")
	var fileconfig = flag.String("fileconfig", "/home/jose/Code/go/src/github.com/josedelrio85/go_components/readparams", "path where to read config file")
	var typeload = flag.String("typeload", "1", "type of data load. 1 => from db; 2=> from csv file")
	var dateinput = flag.String("date", "", "input date")
	flag.Parse()

	f, err := os.OpenFile("./nivoriacomp_log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// prod 3 dev 2
	connstr := readparams.GetConnString(3, *fileconfig)
	wsmsql := &nivoriacomp.Wsmsql{
		Connstring: connstr,
	}
	importer := nivoriacomp.Importer{
		Path:   *basepath,
		Date:   *dateinput,
		Storer: wsmsql,
	}

	// log.Printf("date %s ", *dateinput)
	if err := wsmsql.Open(); err != nil {
		message := fmt.Sprintf("error opening mysql connection. err %s", err)
		nivoriacomp.ResponseError(message, err)
	}

	if *typeload == "2" {
		if err := importer.Importfromcsv(); err != nil {
			message := fmt.Sprintf("error importing data. err %s", err)
			nivoriacomp.ResponseError(message, err)
		}
	} else {
		if err := importer.Importfromdb(); err != nil {
			message := fmt.Sprintf("error importing data. err %s", err)
			nivoriacomp.ResponseError(message, err)
		}
	}

	jsondata := nivoriacomp.Jsondata{
		Network: "josedelrio85_evo",
		Token:   "NVR1ab68b1a04e6bbc5029c2f0e6f5b3d64",
	}

	xmlent := nivoriacomp.XMLEntity{
		Jdata: jsondata,
	}

	if err := xmlent.Request(importer.Data); err != nil {
		message := fmt.Sprintf("error requesting data. err %s", err)
		nivoriacomp.ResponseError(message, err)
	}

	if err := wsmsql.BatchInsert(xmlent.Xdata); err != nil {
		message := fmt.Sprintf("error inserting data. err %s", err)
		nivoriacomp.ResponseError(message, err)
	}

	final := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("Process ended at %s", final)
}

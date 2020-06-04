package loaddatareport

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// Importer imports data
func Importer(newfilename string) error {

	log.Printf("Starting import process at %s", time.Now().Format("2006-01-02 15-04-05"))

	port := getSetting("DB_PORT")
	portInt, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing to string Database's port %s, Err: %s", port, err)
	}

	database := struct {
		Host   string
		User   string
		Pass   string
		Port   int64
		Dbname string
		File   string
	}{
		Host: getSetting("DB_HOST"),
		User: getSetting("DB_USER"),
		Pass: getSetting("DB_PASS"),
		// Dbname: "testing",
		Dbname: "testing-leontel",
		Port:   portInt,
		// Host:   "leads-pre.c848y92oajny.eu-west-1.rds.amazonaws.com",
		// User:   "leads",
		// Pass:   "LW3PBzuqy3zfBrqBbbFM",
		// Dbname: "leads",
		// File:   "select * from leads order by id desc limit 10;",
		File: newfilename,
	}

	host := fmt.Sprintf("-h%s", database.Host)
	user := fmt.Sprintf("-u%s", database.User)
	pass := fmt.Sprintf("-p%s", database.Pass)
	portt := fmt.Sprintf("-P %d", database.Port)
	db := fmt.Sprintf("-D%s", database.Dbname)
	// file := fmt.Sprintf("-e %s", database.File)
	file := fmt.Sprintf("source ./%s", database.File)

	// cmd := exec.Command("/usr/bin/mysql", "-h127.0.0.1", "-P 3306", "-uroot", "-proot_bsc", "-f", "-Dwebservice", "-e show tables;")
	// cmd := exec.Command("/usr/bin/mysql", host, portt, user, pass, db, "-f", file)
	// cmd := exec.Command("/usr/bin/mysql", "-h127.0.0.1", "-P 3306", "-uroot", "-proot_bsc", "-f", "-Dwebservice", "-e ", "source filename.sql")
	cmd := exec.Command("/usr/bin/mysql", host, portt, user, pass, db, "-f", "-e", file)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(stdout)
	// err = ioutil.WriteFile("./out.sql", bytes, 0644)
	if err != nil {
		return err
	}
	log.Println(string(bytes))

	byteserr, err := ioutil.ReadAll(stderr)
	// err = ioutil.WriteFile("./err.sql", byteserr, 0644)
	if err != nil {
		return err
	}
	log.Println(string(byteserr))

	log.Printf("Import process finished! At %s", time.Now().Format("2006-01-02 15-04-05"))
	return nil
}

func getSetting(setting string) string {
	value, ok := os.LookupEnv(setting)
	if !ok {
		log.Fatalf("Init error, %s ENV var not found", setting)
	}

	return value
}

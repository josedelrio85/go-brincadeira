package dumpfootel

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Properties is a struct
type Properties struct {
	Tables   []string
	Dbname   string
	Filename string
}

// Dumper imports data
func Dumper(prop Properties) error {

	log.Printf("Starting dump process at %s", time.Now().Format("2006-01-02 15-04-05"))

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
		Tables []string
	}{
		Host:   getSetting("DB_HOST"),
		User:   getSetting("DB_USER"),
		Pass:   getSetting("DB_PASS"),
		Dbname: prop.Dbname,
		Port:   portInt,
		Tables: prop.Tables,
	}

	host := fmt.Sprintf("-h%s", database.Host)
	user := fmt.Sprintf("-u%s", database.User)
	pass := fmt.Sprintf("-p%s", database.Pass)
	portt := fmt.Sprintf("-P %d", database.Port)
	// db := fmt.Sprintf("--databases %s", database.Dbname)
	// list := fmt.Sprintf("--tables ./%s", strings.Join(database.Tables[:], ","))
	// destination := fmt.Sprintf("> %s_%s.sql", prop.Filename, time.Now().Format("2006-01-02"))
	// destination2 := fmt.Sprintf("./%s_%s.sql", prop.Filename, time.Now().Format("2006-01-02"))

	args := []string{}
	args = append(args, host)
	args = append(args, user)
	args = append(args, pass)
	args = append(args, portt)
	args = append(args, "-f")
	args = append(args, "--databases")
	args = append(args, database.Dbname)
	args = append(args, "--tables")
	for _, t := range database.Tables {
		args = append(args, t)
	}
	// args = append(args, "cat_categories")
	// args = append(args, "cli_clients")

	// args = append(args, "> out.sql")

	log.Println(strings.Join(args, " "))
	// return nil

	// cmd := exec.Command("/usr/bin/mysqldump", "-h127.0.0.1", "-P 3306", "-uroot", "-proot_bsc", "-f", "crmti",
	// 	"cat_categories", "cli_clients", "dni_dnis").
	cmd := exec.Command("/usr/bin/mysqldump", args...)

	// log.Println(string(cmd))

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
	err = ioutil.WriteFile("./out.txt", bytes, 0644)
	if err != nil {
		return err
	}
	// log.Println(string(bytes))

	byteserr, err := ioutil.ReadAll(stderr)
	if err != nil {
		return err
	}
	log.Println(string(byteserr))

	log.Printf("Dump process finished! At %s", time.Now().Format("2006-01-02 15-04-05"))
	return nil
}

func getSetting(setting string) string {
	value, ok := os.LookupEnv(setting)
	if !ok {
		log.Fatalf("Init error, %s ENV var not found", setting)
	}

	return value
}

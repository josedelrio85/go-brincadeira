package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"time"
)

// State is a struct that represents the structure of a status call
type State struct {
	Status string `json:"call_status"`
	Phone  string `json:"call_telephone"`
	Weight int
}

func main() {

	var printall = flag.Bool("printall", false, "Set to true if you want to print not desired phones")
	var basepath = flag.String("basepath", "C:\\Users\\Jose\\go\\src\\github.com\\bysidecar\\go_components\\statecalllog", "path to read the posted file")
	var seconds = flag.Int("seconds", 10, "Set the amount of time you want to log")

	flag.Parse()

	logname := fmt.Sprintf("statecall_log_%s", time.Now().Format("2006-01-02_15-04-05"))
	f, err := os.OpenFile(logname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	for i := 1; i <= *seconds; i++ {
		time.Sleep(2 * time.Second)
		if err := read(*basepath, *printall); err != nil {
			log.Printf("Error in read function, Err: %v", err)
			return
		}
		log.Println("--------------------------------------")
		log.Print("\r\n")
	}
}

func read(basepath string, printall bool) error {
	file := path.Join(basepath, "dialerCalls.json")
	jsonFile, err := os.Open(file)
	if err != nil {
		log.Printf("Error reading json file. Err: %v", err)
		return err
	}

	defer jsonFile.Close()

	bytevalue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Printf("Error getting byte values from json file, Err: %v", err)
		return err
	}

	states := []State{}
	json.Unmarshal(bytevalue, &states)
	final := setStatus(states)

	sort.Slice(final[:], func(i, j int) bool {
		return final[i].Weight < final[j].Weight
	})

	print(final, printall)
	return nil
}

func setStatus(states []State) []State {
	final := []State{}

	for _, state := range states {
		switch a := state.Status; a {
		case "CONNECTING":
			state.Weight = 0
		case "ORIGINATE":
			state.Weight = 1
		case "ONQUEUE":
			state.Weight = 2
		case "ONIVR":
			state.Weight = 3
		case "ONAGENT":
			state.Weight = 4
		case "FINISHING":
			state.Weight = 5
		}
		final = append(final, state)
	}
	return final
}

func print(states []State, printall bool) {
	observed := map[string]bool{
		"665932355": true,
		"638769068": true,
		"685511584": true,
	}

	for _, state := range states {
		if observed[state.Phone] {
			log.Printf("--------------------> %v", state)
		} else {
			if printall {
				log.Println(state)
			}
		}
	}
}

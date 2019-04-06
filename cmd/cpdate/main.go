package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal("error getting all the files from the current directory, Err: %v", err)
	}

	for _, file := range files {
		path := fmt.Sprintf("./%d/%d/%d/",
			file.ModTime().Year(),
			file.ModTime().Month(),
			file.ModTime().Day(),
		)

		fmt.Printf("Copying file %v to %v...\n", file.Name(), path+file.Name())

		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Fatal("error create the destination directory for '%s', Err: %v", file.Name(), err)
		}

		input, err := ioutil.ReadFile(file.Name())
		if err != nil {
			log.Fatal("error reading file to copy '%s', Err: %v", file.Name(), err)
			return
		}

		err = ioutil.WriteFile(path+file.Name(), input, 0644)
		if err != nil {
			log.Fatal("error writing copied file '%s', Err: %v", file.Name(), err)
			return
		}
	}
}

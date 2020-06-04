package loaddatareport

import (
	"compress/gzip"
	"io"
	"log"
	"os"
	"strings"
)

// Unzip a file
func Unzip(filename string) error {
	gzipfile, err := os.Open(filename)

	if err != nil {
		return err
	}

	reader, err := gzip.NewReader(gzipfile)
	if err != nil {
		return err
	}
	defer reader.Close()

	newfilename := strings.TrimSuffix(filename, ".gz")

	writer, err := os.Create(newfilename)

	if err != nil {
		return err
	}
	defer writer.Close()

	if _, err = io.Copy(writer, reader); err != nil {
		return err
	}
	log.Printf("%s unzipped succesful", newfilename)
	return nil
}

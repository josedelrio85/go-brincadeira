package loaddatareport

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Handler is an struct
type Handler struct {
}

type payload struct {
	Name string `json:"name"`
}

// HandleFunction is a function used to manage all received requests.
func (ch *Handler) HandleFunction() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Process launched at %s", time.Now().Format("2006-01-02 15-04-05"))

		path := strings.Split(r.URL.Path, "/")
		bucket := html.EscapeString(path[2])

		payload := payload{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			log.Fatal(err)
		}
		log.Println(bucket)
		log.Println(payload.Name)

		// return

		filename := payload.Name
		// get file from s3
		if err := GetFromS3(filename, bucket); err != nil {
			message := fmt.Sprintf("Error retrieving file, Err: %v", err)
			responseError(w, message, err)
			return
		}

		//unzip it
		if err := Unzip(filename); err != nil {
			message := fmt.Sprintf("Error unzipping file, Err: %v", err)
			responseError(w, message, err)
			return
		}

		newfilename := strings.TrimSuffix(filename, ".gz")
		//import to db
		if err := Importer(newfilename); err != nil {
			message := fmt.Sprintf("Error importing data, Err: %v", err)
			responseError(w, message, err)
			return
		}

		//delete .sql and .sql.gz
		if err := os.Remove(filename); err != nil {
			message := fmt.Sprintf("Error removing file, Err: %v", err)
			responseError(w, message, err)
		}
		log.Printf("%s removed succesfully", filename)

		if err := os.Remove(newfilename); err != nil {
			message := fmt.Sprintf("Error removing file, Err: %v", err)
			responseError(w, message, err)
		}
		log.Printf("%s removed succesfully", newfilename)

		log.Printf("Process ended at %s", time.Now().Format("2006-01-02 15-04-05"))
		responseOk(w, "OK")
	})
}

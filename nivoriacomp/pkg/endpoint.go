package nivoriacomp

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

// Jsondata represents a Json object passed as parameter to the endpoint
type Jsondata struct {
	Network string `json:"network"`
	Token   string `json:"token"`
	Ide     string `json:"id"`
}

// Xmlstruct represents the response of the endpoint
type Xmlstruct struct {
	Idsf   string `xml:"id_Salesforce"`
	Suborg string `xml:"suborigen"`
	Org    string `xml:"origen"`
	Stepid struct {
		ID   int64  `xml:"id"`
		Name string `xml:"name"`
	} `xml:"stepid"`
	CreDate string
}

// XMLEntity is a structure to encapsulate Jsondata and Xmlstruct structs
type XMLEntity struct {
	Jdata Jsondata
	Xdata []Xmlstruct
}

// Request method makes a set of request to the endpoint.
// Marshal input data to a Json string before the request is made
// and then unmarshal the XML response to an array of structs.
func (xe *XMLEntity) Request(data []Inputdata) error {

	for _, dat := range data {
		xe.Jdata.Ide = dat.Clientid

		json, err := json.Marshal(xe.Jdata)
		if err != nil {
			return err
		}

		url := "http://www.nivolab.com/dev/api/evo/getGoal.php"

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(json))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var xmlResult Xmlstruct
		if err := xml.Unmarshal(data, &xmlResult); err != nil {
			return err
		}

		if (Xmlstruct{}) == xmlResult {
			xmlResult.Idsf = xe.Jdata.Ide
		}
		// if xmlResult.Stepid.Name != dat.Stepid {
		// 	xmlResult.Stepid.Name = dat.Stepid
		// }
		xmlResult.CreDate = dat.Createddate
		xe.Xdata = append(xe.Xdata, xmlResult)
	}
	return nil
}

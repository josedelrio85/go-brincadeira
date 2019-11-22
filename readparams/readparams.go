package readparams

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type dbParams struct {
	ReportPanelDev []struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		Db       string `json:"db"`
	} `json:"report_panel_dev"`
	WebserviceDev []struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		Db       string `json:"db"`
	} `json:"webservice_dev"`
	Webservice []struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		Db       string `json:"db"`
	} `json:"webservice"`
	ReportPanel []struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		Db       string `json:"db"`
	} `json:"report_panel"`
	Crmti []struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		Db       string `json:"db"`
	} `json:"crmti"`
}

func getDbParams(fileconfig string) dbParams {
	var dbparams dbParams
	fileconfig = path.Join(fileconfig, "config.development.json")

	file, err := os.Open(fileconfig)
	if err != nil {
		fmt.Println(err)
		return dbparams
	}

	defer file.Close()

	body, err := ioutil.ReadAll(file)
	err = json.Unmarshal(body, &dbparams)

	if err != nil {
		return dbparams
	}
	return dbparams
}

// GetConnString returns connstring for development enviroment (report_panel && webservice)
func GetConnString(conNumber int, fileconfig string) string {

	params := getDbParams(fileconfig)
	switch conNumber {
	case 1:
		return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", params.ReportPanelDev[0].User, params.ReportPanelDev[0].Password, params.ReportPanelDev[0].Host, params.ReportPanelDev[0].Port, params.ReportPanelDev[0].Db)
	case 2:
		return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", params.WebserviceDev[0].User, params.WebserviceDev[0].Password, params.WebserviceDev[0].Host, params.WebserviceDev[0].Port, params.WebserviceDev[0].Db)
	case 3:
		return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", params.Webservice[0].User, params.Webservice[0].Password, params.Webservice[0].Host, params.Webservice[0].Port, params.Webservice[0].Db)
	case 4:
		return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", params.ReportPanel[0].User, params.ReportPanel[0].Password, params.ReportPanel[0].Host, params.ReportPanel[0].Port, params.ReportPanel[0].Db)
	case 5:
		return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", params.Crmti[0].User, params.Crmti[0].Password, params.Crmti[0].Host, params.Crmti[0].Port, params.Crmti[0].Db)
	}
	return ""
}

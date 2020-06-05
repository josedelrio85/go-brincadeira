package main

import (
	"log"
	"net/http"
	"runtime"
	"time"

	dumper "github.com/bysidecar/go_components/dumpleontel/pkg"
	"github.com/bysidecar/voalarm"
)

func main() {
	log.Printf("Starting process at %s", time.Now().Format("2006-01-02 15-04-05"))

	propleontel := dumper.Properties{
		Tables: []string{
			"cat_categories",
			"cli_clients",
			"dni_dnis",
			"ord_lines",
			"ord_orders",
			"pro_products",
			"que_queues",
			"que_queues_description",
			"rel_gro_usr",
			"rel_pro_cat",
			"rel_pro_groups",
			"rel_prof_gro",
			"rel_prof_sub",
			"rel_que_sub",
			"rel_rep_usr",
			"rel_sal_pro",
			"rel_sou_sub",
			"sou_sources",
			"sub_subcategories",
			"typ_types",
			"user_log",
			"act_activity",
			"his_history",
			"lea_leads",
			"usr_users",
		},
		Dbname:   "crmti",
		Filename: "leontel",
	}
	if err := dumper.Dumper(propleontel); err != nil {
		sendAlarm("Error dumping leontel", http.StatusInternalServerError, err)
		log.Fatal(err)
	}

	// sixmonths := time.Now().Add(time.Duration(-4320) * time.Hour) // -6 months
	// propastcdr := dumper.Properties{
	// 	Tables: []string{
	// 		"ast_cdr",
	// 		fmt.Sprintf("--where='date(calldate) >= %s'", sixmonths.Format("2006-01-02")),
	// 	},
	// 	Dbname:   "asterisk",
	// 	Filename: "ast_cdr",
	// }
	// if err := dumper.Dumper(propastcdr); err != nil {
	// 	sendAlarm("Error dumping ast_cdr", http.StatusInternalServerError, err)
	// 	log.Fatal(err)
	// }

	// propqueue := dumper.Properties{
	// 	Tables: []string{
	// 		"tel_queue_activity",
	// 		fmt.Sprintf("--where='date(tel_queue_act_ts) >= %s'", sixmonths.Format("2006-01-02")),
	// 	},
	// 	Dbname:   "asterisk",
	// 	Filename: "tel_queue_activity",
	// }
	// if err := dumper.Dumper(propqueue); err != nil {
	// 	sendAlarm("Error dumping tel_queue_activity", http.StatusInternalServerError, err)
	// 	log.Fatal(err)
	// }

	log.Printf("Process ended  at %s", time.Now().Format("2006-01-02 15-04-05"))
}

// sendAlarm to VictorOps plattform and format the error for more info
func sendAlarm(message string, status int, err error) {
	fancyHandleError(err)

	mstype := voalarm.Acknowledgement
	switch status {
	case http.StatusInternalServerError:
		mstype = voalarm.Warning
	case http.StatusUnprocessableEntity:
		mstype = voalarm.Info
	}

	alarm := voalarm.NewClient("")
	_, err = alarm.SendAlarm(message, mstype, err)
	if err != nil {
		fancyHandleError(err)
	}
}

// fancyHandleError logs the error and indicates the line and function
func fancyHandleError(err error) (b bool) {
	if err != nil {
		// using 1 => it will actually log where the error happened, 0 = this function.
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
		b = true
	}
	return
}

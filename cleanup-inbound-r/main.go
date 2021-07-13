package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/josedelrio85/go_components/implementeddb"
	"github.com/josedelrio85/go_components/readparams"
	"github.com/josedelrio85/voalarm"
)

//Database is a struct that represents MySQL db instance
type Database struct {
	db *sql.DB
}

func main() {
	var fileconfig = flag.String("fileconfig", "C:\\Users\\Jose\\go\\src\\github.com\\josedelrio85\\go_components\\readparams", "path where to read config file")
	flag.Parse()

	file, err := os.OpenFile("./cleanup-inbound-r_log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		sendAlarm(err)
		return
	}
	defer file.Close()

	log.SetOutput(file)

	connString := readparams.GetConnString(5, *fileconfig)
	db, err := implementeddb.OpenConnection(connString)
	if err != nil {
		sendAlarm(err)
		return
	}

	database := Database{
		db: db,
	}

	today := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	sql1 := `DROP TABLE if EXISTS crmti.inboundRClose;`
	if _, err := database.ExecQuery(sql1); err != nil {
		sendAlarm(err)
		return
	}

	sql2 := fmt.Sprintf(`
		CREATE TABLE crmti.inboundRClose 
		SELECT lea_id FROM lea_leads  
		WHERE lea_closed = 0  
		AND TELEFONO in ('881920590','946858791','946858792','946858793','946858794','946858795')
		AND lea_source = 81
		AND lea_type = 9 
		AND date(lea_ts) <= date('%s')
		AND observaciones2 is null;`, today)
	if _, err := database.ExecQuery(sql2); err != nil {
		sendAlarm(err)
		return
	}

	var count2 int
	database.db.QueryRow(`select count(*) from crmti.inboundRClose;`).Scan(&count2)

	sql3 := fmt.Sprintf(`
		UPDATE crmti.lea_leads SET 
		observaciones2='CERRADO POR LIMPIEZA INBOUND R RCABLE %s' 
		where lea_id in (select lea_id from crmti.inboundRClose);`,
		today)

	count3, err := database.ExecQuery(sql3)

	if err != nil {
		sendAlarm(err)
		return
	}

	if count2 != int(count3) {
		msg := fmt.Sprintf("An error ocurred. Number of updated lines [%d] does not match the main query [%d].", count3, count2)
		err := errors.New(msg)
		sendAlarm(err)
		return
	}

	sql4 := fmt.Sprintf(`
		INSERT INTO crmti.his_history (his_lead,his_user,his_sub,his_dest_user,his_comment,his_newlead)  
		select lea_id,9001,817,null,concat(CURRENT_TIMESTAMP(),'CERRADO POR LIMPIEZA INBOUND R RCABLE %s'), 0 
		from crmti.inboundRClose;`, today)

	if _, err := database.ExecQuery(sql4); err != nil {
		sendAlarm(err)
		return
	}

	if _, err := database.ExecQuery(sql1); err != nil {
		sendAlarm(err)
		return
	}

	log.Printf(`Cleanup process succeeds! [%s] Registries updated: %d`, today, count2)
}

// ExecQuery prepares a query without parameters and execute this prepared statement
// Returns number of elements affected or error
func (d Database) ExecQuery(query string) (int64, error) {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	// rows, err := stmt.Query()
	result, err := stmt.Exec()
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	n, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return n, nil
}

func sendAlarm(err error) {
	log.Println(err)
	alarm := voalarm.NewClient("")
	alarm.SendAlarm("cleanup-inbound-r", voalarm.Acknowledgement, err)
}

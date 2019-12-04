# go_components

Little and funcional Go! components that fullfill procedure needs of loading third party data onto databases.

## loadFormalizadasCSV

Used to load information from EVO's CRM (Salesforce) into the SIRA database.

> Requires a CSV file named "formalizadas.csv" downloaded form EVO's CRM on the same folder than the executable
> Records that match this condition `FECHA_FORMALIZACION >= '2018-12-01'` are deleted before the inserts.

After executing the command check the log file generated on the parent folder, it should be empty.

Use the following statement to check that the import procedure has succeeded, you should see that there are rows from the day before on the table:

```sql
SELECT * FROM webservice.evo_formalizadas_sf_v2 where date(FECHA_FORMALIZACION) >= '2019-01-01' order by FECHA_FORMALIZACION desc limit 10;
```

* To compile as linux distribution execute:

```bash
  set GOOS=linux
  go build
```

* Upload the compiled object to Sira host and set correct permission and owner.

```bash
/var/www/vhosts/dashboard.bysidecar.es/custom/plugins/loadFormalizadasCSV
chown apache:apache loadFormalizadasCSV
chmod 755 loadFormalizadasCSV
```

## loadFirmasCSV

> Requires a CSV file named "firmas.csv" downloaded form EVO's CRM on the same folder than the executable

* Records that match this condition `FECHA_DE_FIRMA>= '2018-12-01'` are deleted before the inserts.

After executing the command check the log file generated on the parent folder, it should be empty.

Use the following statement to check that the import procedure has succeeded, you should see that there are rows from the day before on the table:

```sql
SELECT * FROM webservice.evo_firmados_sf_v2 where date(Fecha_de_firma) >= '2018-12-01' order by Fecha_de_firma desc limit 10;
```

* To compile as linux distribution execute:

```bash
  set GOOS=linux
  go build
```

* Upload the compiled object to Sira host and set correct permission and owner.

```bash
/var/www/vhosts/dashboard.bysidecar.es/custom/plugins/loadFirmasCSV
chown apache:apache loadFirmasCSV
chmod 755 loadFirmasCSV
```

## voalarm

This functionality can send an alarm to VictorOps plattform when an execption occurs.

```go
  // Create an instance of Client struct with no params. SendAlarm will set correct values.
 a := fmt.Errorf("Error. Artificial error: %v", errors.New("emit macho dwarf: elf header corrupted"))
 alarm := voalarm.Client{}
 resp, err := alarm.SendAlarm(voalarm.Acknowledgement, a)
 if err != nil {
  log.Fatalf("Error creating alarm. Err: %s", err)
 }
 log.Printf("Response: %+v\n", resp)

 // Create an instance of Client struct with apikey param setted to "". SendAlarm will set correct values too.
 a := fmt.Errorf("Error. Artificial error: %v", errors.New("emit macho dwarf: elf header corrupted"))
 alarm := voalarm.NewClient("")
 resp, err := alarm.SendAlarm(voalarm.Acknowledgement, a)
 if err != nil {
  log.Fatalf("Error creating alarm. Err: %s", err)
 }
 log.Printf("Response: %+v\n", resp)
```

## encode64

This functionality encodes to base64 string an input csv file with one column of data (mainly phone numbers but will decode any string).

Returns an output csv file with two columns, the first is the input data and the second is the encoded data.

This functionality will look, by default, for an input csv file which name starts by 'input'.

You can pass an entry param --filename indicating the file name to encode.

## State Call Log

This functionality logs the states of the current calls registered in asterisk ast_current_calls table.

It does not interact with the asterisk system directly. Instead it reads a json file generated throw an external service that is reading this table continously.

### How to launch

* Get ssh access to 192.168.50.20 server

* Execute the main file with this flags:

  * Seconds flag: ammount of time to log the results. Default 10 seconds.

  * Printall flag: print all the entries in the json file. Default false (only will write reseverd number phones)

  * basepath: json file path
  
  ```bash
  ./main --seconds=15 --printall=true --basepath=/var/www/privateBySidecar/realTimeData
  ```

## Nivoria

* This taks is launched every day at 02:20.

* This query is launched against our DB to get the data about CLIENTID and CREATEDDATE fields.

```sql
select CLIENTID, CREATEDDATE FROM webservice.evo_events_sf_v2_pro    where date(CREATEDDATE) = [yesterday]
```

* The data retrieved is stored into an array and this data is provided as param to the next endpoint (POST / JSON):

<http://www.nivolab.com/dev/api/evo/getGoal.php>

* The results obtained are stored in another array and a batch insert is made into `evo_origen_idcliente` table.

* To compile as linux distribution execute:

```bash
  set GOOS=linux
  go build
```

* Upload the compiled object to Webserice host and set correct permission and owner.

```bash
/etc/srv/bysidecar/bin/github.com/bysidecar/nivoriacomp/
chown apache:apache nivoriacomp
chmod 755 nivoriacomp
```

* Instruction to automatic execute the script

```bash
20 2 * * * cd /etc/srv/bysidecar/bin/github.com/bysidecar/nivoriacomp/ && ./nivoriacomp -fileconfig=/var/www/privateBySidecar/ > /var/log/nivoriacomp.log 2>&1
```

## Cleanup Evo leads

* This taks is launched every day at 03:20.

* First a list of ID's is retrieved from lea_leads table for Evo campaign to get leads that were auto assigned to consultants.

* Then an update is made to close this leads and mark it setting a text in observations2 field.

* Finally an insert in his_history table is made to create the row that indicates that the lead was closed.

* To compile as linux distribution execute:

```bash
  set GOOS=linux
  go build
```

* Upload the compiled object to Webserice host and set correct permission and owner.

```bash
/etc/srv/bysidecar/bin/github.com/bysidecar/cleanup-evo-leads/
chown apache:apache cleanup-evo-leads
chmod 755 cleanup-evo-leads
```

* Instruction to automatic execute the script

```bash
20 3 * * * cd /etc/srv/bysidecar/bin/github.com/bysidecar/cleanup-evo-leads/ && ./cleanup-evo-leads -fileconfig=/******/privateBySidecar/ > /var/log/cleanup-evo-leads.log 2>&1
```

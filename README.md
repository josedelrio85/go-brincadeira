# go_components
Little and funcional Go! components that fullfill procedure needs of loading third party data onto databases.

# Included procedures

## loadFormalizadasCSV

Used to load information from EVO's CRM (Salesforce) into the SIRA database.

> Requires a CSV file named "formalizadas.csv" downloaded form EVO's CRM on the same folder than the executable 
  
> Records that match this condition `FECHA_FORMALIZACION >= '2018-12-01'` are deleted before the inserts.

After executing the command check the log file generated on the parent folder, it should be empty.

Use the following statement to check that the import procedure has succeeded, you should see that there are rows from the day before on the table:

```
SELECT * FROM webservice.evo_formalizadas_sf_v2 where date(FECHA_FORMALIZACION) >= '2019-01-01' order by FECHA_FORMALIZACION desc limit 10;
```
  
## loadFirmasCSV 

> Requires a CSV file named "firmas.csv" downloaded form EVO's CRM on the same folder than the executable 

* Records that match this condition `FECHA_DE_FIRMA>= '2018-12-01'` are deleted before the inserts.

After executing the command check the log file generated on the parent folder, it should be empty.

Use the following statement to check that the import procedure has succeeded, you should see that there are rows from the day before on the table:

```
SELECT * FROM webservice.evo_firmados_sf_v2 where date(Fecha_de_firma) >= '2018-12-01' order by Fecha_de_firma desc limit 10;
```


## voalarm

This functionality can send an alarm to VictorOps plattform when an execption occurs. 

```
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

  - Get ssh access to 192.168.50.20 server

  - Execute the main file with this flags:
  
    - Seconds flag: ammount of time to log the results. Default 10 seconds.

    - Printall flag: print all the entries in the json file. Default false (only will write reseverd number phones)

    - basepath: json file path
 
 ```
./main --seconds=15 --printall=true --basepath=/var/www/privateBySidecar/realTimeData
```

## Nivoria

  - This taks is launched every day at 02:20.

  - This query is launched against our DB to get the data about CLIENTID and CREATEDDATE fields.

```
select CLIENTID, CREATEDDATE FROM webservice.evo_events_sf_v2_pro    where date(CREATEDDATE) = [yesterday]
```

  - The data retrieved is stored into an array and this data is provided as param to the next endpoint (POST / JSON):

```
http://www.nivolab.com/dev/api/evo/getGoal.php
```

   - The results obtained are stored in another array and a batch insert is made into `evo_origen_idcliente` table.
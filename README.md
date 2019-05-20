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
# go_components
Little and funcional Go! components


## loadFormalizadasCSV

  * Reads a CSV file called "formalizadas.csv" with a defined structure
  * Iterates over the rows
  * For each row, populates a custom Struct and creates a prepared statement for MySQL using (?,?,?,?)
  * When the read process is finished, the statement are prepared and executed.
  * The destination of the insertions are 2 different schemas, production environment both. 
  * The records that match this condition => FECHA_FORMALIZACION >= '2018-12-01' are deleted before the inserts.
  
  
## loadFirmasExcel

  * Reads a CSV file called "firmas.xlsx" with a defined structure
  * Iterates over the rows
  * For each row, populates a custom Struct and creates a insert statement.
  * Each insert statement are concatenated making SQL instructions of len(totalRows)/10, to avoid timeout problems.
  * When the read process is finished, the statement are prepared and executed.
  * The destination of the insertions are 2 different schemas, production environment both. 
  * The records that match this condition => Fecha_de_firma >= '2018-12-01' are deleted before the inserts.

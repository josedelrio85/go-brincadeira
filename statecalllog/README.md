# State Call Log

This functionality logs the states of the current calls registered in asterisk ast_current_calls table. 

It does not interact with the asterisk system directly. Instead it reads a json file generated throw an external service that is reading this table continously.


## How to launch

  - Get ssh access to 192.168.50.20 server

  - Execute the main file with this flags:
  
    - Seconds flag: ammount of time to log the results. Default 10 seconds.

    - Printall flag: print all the entries in the json file. Default false (only will write reseverd number phones)

    - basepath: json file path
 
 ```
./main --seconds=15 --printall=true --basepath=/var/www/privateBySidecar/realTimeData
```
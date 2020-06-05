#!/bin/sh

#----------------------------------------------------------
# a simple mysql database backup script.
# version 2, updated March 26, 2011.
# copyright 2011 alvin alexander, http://alvinalexander.com
#----------------------------------------------------------
# This work is licensed under a Creative Commons 
# Attribution-ShareAlike 3.0 Unported License;
# see http://creativecommons.org/licenses/by-sa/3.0/ 
# for more information.
#----------------------------------------------------------

# (1) set up all the mysqldump variables
FILE=ac_`date +"%Y-%m-%d"`.sql
FILE2=tq_`date +"%Y-%m-%d"`.sql
FILE3=asterisk_`date +"%Y-%m-%d"`.sql
DBSERVER=127.0.0.1
DATABASE=asterisk
USER=root
PASS=root_bsc

# (2) in case you run this more than once a day, remove the previous version of the file
unalias rm     2> /dev/null
rm ${FILE}     2> /dev/null
rm ${FILE}.gz  2> /dev/null

# (3) do the mysql database backup (dump)
mysqldump --host=${DBSERVER} --user=${USER} --password=${PASS} -P 3306 -f --databases ${DATABASE} --tables ast_cdr --where="date(calldate) >= '2020-01-01'" > ${FILE}
mysqldump --host=${DBSERVER} --user=${USER} --password=${PASS} -P 3306 -f --databases ${DATABASE} --tables tel_queue_activity --where="date(tel_queue_act_ts) >= '2020-01-01'" > ${FILE2}

# join files
cat ${FILE} ${FILE2} > ${FILE3}
# (4) gzip the mysql database dump file
gzip $FILE3

# (5) show the user the result
echo "${FILE3}.gz was created:"
ls -l ${FILE3}.gz

# (6) push to s3 bucket
echo "Uploading file ${FILE3}.gz"
# /usr/bin/aws s3 cp - s3://data.bysidecar.me/backups/leontel/${FILE}.gz
aws s3 cp ${FILE3}.gz s3://data.bysidecar.me/backups/leontel/${FILE3}.gz
echo "Upload finished"


# (7) delete created files
rm ${FILE}     2> /dev/null
rm ${FILE}.gz  2> /dev/null

rm ${FILE2}     2> /dev/null
rm ${FILE2}.gz  2> /dev/null

rm ${FILE3}     2> /dev/null
rm ${FILE3}.gz  2> /dev/null
echo "Files removed"



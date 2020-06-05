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
FILE=`date +"%Y_%m_%d"`.sql
DBSERVER=127.0.0.1
DATABASE=crmti
USER=root
PASS=root_bsc

# (2) in case you run this more than once a day, remove the previous version of the file
unalias rm     2> /dev/null
rm ${FILE}     2> /dev/null
rm ${FILE}.gz  2> /dev/null

# (3) do the mysql database backup (dump)
mysqldump --host=${DBSERVER} --user=${USER} --password=${PASS} -P 3306 -f --databases ${DATABASE} --tables cat_categories cli_clients dni_dnis ord_lines ord_orders ope_operation pro_products que_queues que_queues_description rel_gro_usr rel_pro_cat rel_pro_groups rel_prof_gro rel_prof_sub rel_que_sub rel_rep_usr rel_sal_pro rel_sou_sub sou_sources sub_subcategories typ_types user_log act_activity his_history lea_leads usr_users > ${FILE}

# (4) gzip the mysql database dump file
gzip $FILE

# (5) show the user the result
echo "${FILE}.gz was created:"
ls -l ${FILE}.gz

# (6) push to s3 bucket
echo "Uploading file ${FILE}.gz"
/usr/bin/aws s3 cp - s3://data.bysidecar.me/backups/leontel/${FILE}.gz
echo "Upload finished"


# (7) delete created files
rm ${FILE}     2> /dev/null
rm ${FILE}.gz  2> /dev/null


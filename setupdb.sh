#!/bin/bash

# TODO(ttacon): make sure the user provides username/password?
# TODO(ttacon): maybe use sql file to run after user logs in to mysql?

mysql -u $1 -p$2 -e "CREATE DATABASE pouch;"
if [ $? -ne 0 ]; then
    echo "Database pouch already exists...";
fi

mysql -u $1 -p$2 -e "GRANT ALL ON pouch.* TO pouch@localhost IDENTIFIED BY 'pouch';";
mysql -u $1 -p$2 -e "use pouch; CREATE TABLE Food (ID int primary key auto_increment not null, Name varchar(255) not null, NullableField varchar(32)) engine=InnoDB;";


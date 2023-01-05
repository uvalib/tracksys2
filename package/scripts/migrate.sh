#!/usr/bin/env bash
#
# run any necessary migrations
#

# run the migrations
#bin/migrate -path db -verbose -database "mysql://$DBUSER:$DBPASS@tcp($DBHOST:$DBPORT)/$DBNAME?x-migrations-table=ts2_schema_migrations" up
echo "WOULD DO: bin/migrate -path db -verbose -database 'mysql://$DBUSER:$DBPASS@tcp($DBHOST:$DBPORT)/$DBNAME?x-migrations-table=ts2_schema_migrations' up"


# return the status
exit $?

#
# end of file
#

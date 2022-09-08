#!/usr/bin/env bash
#
# run application
#

# run the server
umask 0002
cd bin; ./tracksys2                \
   -jwtkey $TRACKSYS_JWT_KEY       \
   -reports $REPORTING_SERVICE_URL \
   -projects $PROJECTS_SERVICE_URL \
   -virgo $VIRGO_URL               \
   -iiifman $IIIFMAN_SERVICE_URL   \
   -iiif $IIIF_SERVICE_URL         \
   -curio $CURIO_SERVICE_URL       \
   -jobs $JOBS_URL                 \
   -apollo $APOLLO_URL             \
   -pdf $PDF_SERVICE_URL           \
   -tsapi $TSAPI                   \
   -dbhost $DBHOST                 \
   -dbport $DBPORT                 \
   -dbname $DBNAME                 \
   -dbuser $DBUSER                 \
   -dbpass $DBPASS

# return the status
exit $?

#
# end of file
#

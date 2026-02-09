#!/usr/bin/env bash
#
# run application
#

# run the server
umask 0002
cd bin; ./tracksys2                \
   -jwtkey $TRACKSYS_JWT_KEY       \
   -aptrust $AP_TRUST_URL          \
   -projects $PROJECTS_SERVICE_URL \
   -virgo $VIRGO_URL               \
   -iiifman $IIIFMAN_SERVICE_URL   \
   -iiif $IIIF_SERVICE_URL         \
   -ils $ILS_CONNECTOR_URL         \
   -curio $CURIO_SERVICE_URL       \
   -jobs $JOBS_URL                 \
   -apollo $APOLLO_URL             \
   -pdf $PDF_SERVICE_URL           \
   -solr $SOLR_URL                 \
   -index $INDEX_URL               \
   -xmlhook $XML_INDEX_HOOK        \
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

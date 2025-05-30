#
# build the target application
#
FROM public.ecr.aws/docker/library/golang:1.24.0-alpine3.21 AS builder

# update the packages
RUN apk update && apk upgrade && apk add --no-cache make npm

# needed for the build
#ARG FONT_AWESOME_TOKEN

WORKDIR /build
COPY go.mod go.sum Makefile ./
COPY backend ./backend
COPY frontend ./frontend
COPY data ./data
RUN make linux-full

#
# build the target container
#
FROM public.ecr.aws/docker/library/alpine:3.21

# update the packages
RUN apk update && apk upgrade && apk add bash tzdata ca-certificates curl libxml2-utils && rm -rf /var/cache/apk/*

# Create the run user and group
RUN addgroup webservice && adduser webservice -G webservice -D

# set the timezone appropriatly
ENV TZ=UTC
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Specify home
ENV APP_HOME=/tracksys2
WORKDIR $APP_HOME

# Create necessary directories
RUN mkdir -p $APP_HOME/bin $APP_HOME/scripts $APP_HOME/db
RUN chown -R webservice $APP_HOME && chgrp -R webservice $APP_HOME

# add the migrate tool
ENV MIGRATE_VERSION=v4.18.1
RUN cd $APP_HOME/bin && curl -L https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz | tar xvz && rm LICENSE README.md

# port and run command
EXPOSE 8080
CMD ["scripts/entry.sh"]

# Move in necessary assets
COPY package/data/container_bash_profile /home/webservice/.profile
COPY package/scripts/entry.sh $APP_HOME/scripts/entry.sh
COPY package/scripts/migrate.sh $APP_HOME/scripts/migrate.sh
COPY backend/db/migrations/*.sql $APP_HOME/db/
COPY --from=builder /build/bin/tracksys2.linux $APP_HOME/bin/tracksys2
COPY --from=builder /build/bin/public $APP_HOME/bin/public
COPY --from=builder /build/bin/data $APP_HOME/bin/data

# Ensure permissions are correct
RUN chown webservice:webservice /home/webservice/.profile $APP_HOME/scripts/entry.sh $APP_HOME/bin/tracksys2 && chmod 755 /home/webservice/.profile $APP_HOME/scripts/entry.sh $APP_HOME/bin/tracksys2
RUN chown -R webservice:webservice $APP_HOME/bin/public $APP_HOME/bin/data $APP_HOME/db && chmod -R 755 $APP_HOME/bin/public $APP_HOME/bin/data $APP_HOME/db

# Add the build tag
ARG BUILD_TAG
RUN test -n "$BUILD_TAG" && touch $APP_HOME/buildtag.build-$BUILD_TAG || touch $APP_HOME/buildtag.build-0

# Specify the user
USER webservice

#
# end of file
#

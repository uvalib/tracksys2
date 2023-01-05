# TrackSys 2

This is the new version of tracking software for the Digital Production Group.
It tracks current, in-process and historic digitization orders.

### Requirements
* Go 1.19.0+
* Node 17.9+
* exiftool
* Imagemagick
* Kakadu

### Database Notes

The TrackSys2 backend uses a MySQL DB to track everything. The schema is managed by
https://github.com/golang-migrate/migrate and the scripts are in ./backend/db/migrations.

Install the migrate binary on your host system. For OSX, the easiest method is brew. Execute:

`brew install golang-migrate`.

Define your MySQL connection params in an environment variable, like this:

`export TS2DB='mysql://user:password@tcp(host:port)/dbname?x-migrations-table=ts2_schema_migrations'`

Run migrations like this:

`migrate -database ${TS2DB} -path backend/db/migrations up`

Example migrate commads to create a migration and run one:

* `migrate create -ext sql -dir backend/db/migrations -seq update_user_auth`
* `migrate -database ${TS2DB} -path backend/db/migrations/ up`